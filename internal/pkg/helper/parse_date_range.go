package helper

import (
	"time"

	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
)

func (h *Helper) ParseDateRange(startDate string, endDate string) (time.Time, time.Time, *res.Err) {
	layout := "2006-01-02"

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, res.ErrBadRequest("Invalid start date format. Use YYYY-MM-DD.")
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, res.ErrBadRequest("Invalid start date format. Use YYYY-MM-DD.")
	}

	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return start, end, nil
}
