package email

import (
	"testing"
	"time"
)

func TestUnmarshallJSON(t *testing.T) {

	tests := []struct {
		testName       string
		test           []byte
		expectedResult time.Time
	}{
		{"empty", []byte(""), time.Now().AddDate(0, 0, 1)},
		{"null", []byte("null"), time.Now().AddDate(0, 0, 1)},
		{"zero", []byte("\"0\""), time.Now().AddDate(0, 0, 1)},
		{"zeroint", []byte("0"), time.Now().AddDate(0, 0, 1)},
		{"nil", []byte(nil), time.Now().AddDate(0, 0, 1)},
		{"UTCDate", []byte("\"2024-01-02T15:04:05.000+0000\""), time.Date(2024, 01, 2, 15, 4, 5, 0, time.UTC)},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			et := ExpiryTime{}
			et.UnmarshalJSON(test.test)

			if !et.Time.Truncate(time.Hour).Equal(test.expectedResult.Truncate(time.Hour)) {
				t.Errorf("The string %v was parsed to time %v and does not match the expected parsed time of %v", string(test.test), et, test.expectedResult)
			}

		})
	}

}
