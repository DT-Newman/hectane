package email

import (
	"encoding/json"
	"strings"
	"time"
)

type ExpiryTime struct {
	time.Time
}

func (c ExpiryTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Time.Format("2006-01-02T15:04:05.000-0700"))
}

func (c *ExpiryTime) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || strings.Contains(string(b), "null") || string(b[1:len(b)-1]) == "0" {
		// Empty string: ""
		// If no expiry time is provided then set the expiry time to 24 hours
		c.Time = time.Now().AddDate(0, 0, 1)
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05.000-0700", string(b[1:len(b)-1])) // b[1:len(b)-1] removes the first and last character, as they are quotes
	if err != nil {
		return err
	}

	c.Time = t

	return nil
}
