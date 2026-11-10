package services

import (
	"teacher-wails/internal/utils"
	"time"
)

// ReviewAnnouncement preserves only the latest school day's announcement across
// the seven-calendar-day history cleanup. Empty text must also replace old text.
type ReviewAnnouncement struct {
	Date   string `json:"date"`
	Review string `json:"review"`
}

// SchoolDayReview returns the immediately preceding school day's announcement,
// never an older nonempty reminder. Weekends and configured holidays show none.
func (data *DailyData) SchoolDayReview(date time.Time, holidays []string) string {
	if !utils.IsWorkday(date, holidays) {
		return ""
	}
	previous := date.AddDate(0, 0, -1)
	for !utils.IsWorkday(previous, holidays) {
		previous = previous.AddDate(0, 0, -1)
	}
	key := utils.FormatDate(previous)
	if record, ok := data.Records[key]; ok {
		return record.NextDayReview
	}
	if data.ReviewCarry != nil && data.ReviewCarry.Date == key {
		return data.ReviewCarry.Review
	}
	return ""
}

func (data *DailyData) preserveReview(today time.Time, holidays []string) {
	latest := ""
	if data.ReviewCarry != nil {
		parsed, err := utils.ParseDate(data.ReviewCarry.Date)
		if err == nil && data.ReviewCarry.Date <= utils.FormatDate(today) && utils.IsWorkday(parsed, holidays) {
			latest = data.ReviewCarry.Date
		} else {
			data.ReviewCarry = nil
		}
	}
	for date, record := range data.Records {
		parsed, err := utils.ParseDate(date)
		if err != nil || date > utils.FormatDate(today) || !utils.IsWorkday(parsed, holidays) || date < latest {
			continue
		}
		latest = date
		data.ReviewCarry = &ReviewAnnouncement{Date: date, Review: record.NextDayReview}
	}
}
