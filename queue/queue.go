package queue

import (
	"github.com/hectane/hectane/util"

	"fmt"
	"log"
	"sync"
	"time"
)

// Queue status information.
type QueueStatus struct {
	Uptime int                    `json:"uptime"`
	Hosts  map[string]*HostStatus `json:"hosts"`
}

// Mail queue managing the sending of messages to hosts.
type Queue struct {
	sync.Mutex
	Storage    *Storage
	config     *Config
	hosts      map[string]*Host
	newMessage *util.NonBlockingChan
	startTime  time.Time
	stop       chan bool
}

// Log the specified message.
func (q *Queue) log(msg string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[Queue] %s", msg), v...)
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) deliverMessage(m *Message) {
	q.Lock()
	if _, ok := q.hosts[m.Host]; !ok {
		q.hosts[m.Host] = NewHost(m.Host, q.Storage, q.config)
	}
	q.hosts[m.Host].Deliver(m)
	q.Unlock()
}

// Check for inactive host queues and shut them down.
func (q *Queue) checkForInactiveQueues() {
	q.Lock()
	for n, h := range q.hosts {
		if h.Idle() > time.Minute {
			h.Stop()
			delete(q.hosts, n)
		}
	}
	q.Unlock()
}

// Receive new messages and deliver them to the specified host queue.
func (q *Queue) run() {
	defer close(q.stop)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
loop:
	for {
		select {
		case i := <-q.newMessage.Recv:
			q.deliverMessage(i.(*Message))
		case <-ticker.C:
			q.checkForInactiveQueues()
		case <-q.stop:
			break loop
		}
	}
	q.log("shutting down host queues")
	q.Lock()
	for h := range q.hosts {
		q.hosts[h].Stop()
	}
	q.Unlock()
	q.log("mail queue shutdown")
}

// Create a new message queue. Any undelivered messages on disk will be added
// to the appropriate queue.
func NewQueue(c *Config) (*Queue, error) {
	q := &Queue{
		Storage:    NewStorage(c.Directory),
		config:     c,
		hosts:      make(map[string]*Host),
		newMessage: util.NewNonBlockingChan(),
		startTime:  time.Now(),
		stop:       make(chan bool),
	}
	if messages, err := q.Storage.LoadMessages(); err == nil {
		q.log("loaded %d message(s) from %s", len(messages), c.Directory)
		for _, m := range messages {
			q.deliverMessage(m)
		}
	} else {
		return nil, err
	}
	go q.run()
	return q, nil
}

// Provide the status of each host queue.
func (q *Queue) Status() *QueueStatus {
	s := &QueueStatus{
		Uptime: int(time.Now().Sub(q.startTime) / time.Second),
		Hosts:  make(map[string]*HostStatus),
	}
	q.Lock()
	for n, h := range q.hosts {
		s.Hosts[n] = h.Status()
	}
	q.Unlock()
	return s
}

// Deliver the specified message to the appropriate host queue.
func (q *Queue) Deliver(m *Message) {
	q.newMessage.Send <- m
}

// Stop all active host queues.
func (q *Queue) Stop() {
	q.stop <- true
	<-q.stop
}
