package email

import (
	"encoding/json"
	"time"

	"github.com/hectane/hectane/queue"
)

// Raw represents a raw email message ready for delivery.
type Raw struct {
	From   string     `json:"from"`
	To     []string   `json:"to"`
	Body   string     `json:"body"`
	Expiry ExpiryTime `json:"expiry"`
}

// Custom unmarshalliing, set default values when not present
func (r *Raw) UnmarshalJSON(data []byte) error {

	type Alias Raw
	tmp := struct {
		*Alias
		Expiry *ExpiryTime `json:"expiry"`
	}{
		Alias: (*Alias)(r),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	// If no expiry set then set to default
	if tmp.Expiry == nil {
		r.Expiry.Time = time.Now().AddDate(0, 0, 1) // default
	} else {
		r.Expiry = *tmp.Expiry
	}

	return nil
}

// DeliverToQueue delivers raw messages to the queue.
func (r *Raw) DeliverToQueue(q *queue.Queue) error {
	w, body, err := q.Storage.NewBody()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(r.Body)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	hostMap, err := GroupAddressesByHost(r.To)
	if err != nil {
		return err
	}
	for h, to := range hostMap {
		m := &queue.Message{
			Host:   h,
			From:   r.From,
			To:     to,
			Expiry: r.Expiry.Time,
		}
		if err := q.Storage.SaveMessage(m, body); err != nil {
			return err
		}
		q.Deliver(m)
	}
	return nil
}
