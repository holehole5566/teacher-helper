package services

import (
	"teacher-wails/internal/utils"
	"testing"
)

func TestSchoolDayReview(t *testing.T) {
	data := DailyData{Records: map[string]DailyRecord{
		"2026-11-05": {NextDayReview: "較舊的考試"},
		"2026-11-06": {NextDayReview: "週五預告"},
		"2026-11-08": {NextDayReview: "週日不應採用"},
		"2026-11-09": {NextDayReview: "週一預告"},
	}}
	for _, tc := range []struct {
		name, date string
		holidays   []string
		want       string
	}{
		{"Monday skips weekend", "2026-11-09", nil, "週五預告"},
		{"Tuesday uses Monday", "2026-11-10", nil, "週一預告"},
		{"Tuesday skips holiday Monday", "2026-11-10", []string{"2026-11-09"}, "週五預告"},
		{"Weekend hides reminder", "2026-11-08", nil, ""},
		{"Holiday hides reminder", "2026-11-09", []string{"2026-11-09"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			date, _ := utils.ParseDate(tc.date)
			if got := data.SchoolDayReview(date, tc.holidays); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
	data.Records["2026-11-09"] = DailyRecord{}
	date, _ := utils.ParseDate("2026-11-10")
	if got := data.SchoolDayReview(date, nil); got != "" {
		t.Fatalf("blank school day reused stale reminder: %q", got)
	}
	delete(data.Records, "2026-11-09")
	if got := data.SchoolDayReview(date, nil); got != "" {
		t.Fatalf("missing school day reused stale reminder: %q", got)
	}
}

func TestReviewSurvivesLongHolidayAndClearsOnNextSchoolDay(t *testing.T) {
	holidays := []string{"2026-11-09", "2026-11-10", "2026-11-11", "2026-11-12", "2026-11-13"}
	data := DailyData{Records: map[string]DailyRecord{"2026-11-06": {NextDayReview: "連假後複習考"}}}
	friday, _ := utils.ParseDate("2026-11-06")
	data.preserveReview(friday, holidays)
	delete(data.Records, "2026-11-06") // seven-calendar-day retention expires
	monday, _ := utils.ParseDate("2026-11-16")
	data.preserveReview(monday, holidays)
	if got := data.SchoolDayReview(monday, holidays); got != "連假後複習考" {
		t.Fatalf("lost long-holiday reminder: %q", got)
	}
	data.Records["2026-11-16"] = DailyRecord{} // no new exam announced
	data.preserveReview(monday, holidays)
	delete(data.Records, "2026-11-16")
	tuesday, _ := utils.ParseDate("2026-11-17")
	if got := data.SchoolDayReview(tuesday, holidays); got != "" {
		t.Fatalf("stale reminder carried twice: %q", got)
	}
}
