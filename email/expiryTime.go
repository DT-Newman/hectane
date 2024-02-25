package email

import "time"

type ExpiryTime struct {
	time.Time // Embed time.Time to allow calling of normal time.Time methods
}

func (c *ExpiryTime) UnmarshalJSON(b []byte) error {
	if len(b) < 3 {
		// Empty string: ""
		// If no expiry time is provided then set the expiry time to 24 hours
		c.Time = time.Now().Add(time.Hour * 24)
		return nil
	}

	t, err := time.Parse("2006-01-02T15:04:05.000-0700", string(b[1:len(b)-1])) // b[1:len(b)-1] removes the first and last character, as they are quotes
	if err != nil {
		return err
	}

	c.Time = t

	return nil
}
