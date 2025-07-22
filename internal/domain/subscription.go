package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID   `json:"id"`
	ServiceName string      `json:"service_name"`
	Price       int         `json:"price"`
	UserID      uuid.UUID   `json:"user_id"`
	StartDate   YearMonth   `json:"start_date"`
	EndDate     *YearMonth  `json:"end_date,omitempty"`
}

type YearMonth struct {
	time.Time
}

func (ym *YearMonth) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	ym.Time = t
	return nil
}

func (ym YearMonth) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ym.Time.Format("01-2006") + `"`), nil
}

