package main

import (
	"teacher-wails/internal/models"
	"teacher-wails/internal/services"
	"teacher-wails/internal/utils"
	"testing"
	"time"
)

func TestDailyAssignmentsAndHistoricalContents(t *testing.T) {
	dir := t.TempDir()
	dh, err := services.NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	sm := services.NewStudentManager(dh)
	a := &App{dh: dh, sm: sm, dc: services.NewDutyCalculator(dh, sm)}
	today := utils.FormatDate(time.Now())
	yesterday := utils.FormatDate(time.Now().AddDate(0, 0, -1))
	original := models.Student{SeatNumber: 1, Name: "原始姓名"}
	dh.SaveStudents([]models.Student{original})
	dh.UpdateDailyData(func(data *services.DailyData) error {
		data.Records[yesterday] = services.DailyRecord{Date: yesterday, Students: []models.Student{original}, ContactItems: []string{"昨天事項"}, MorningReview: "昨天複習考", LayoutVersion: 1}
		return nil
	})
	changed := models.Student{SeatNumber: 2, Name: "替換同學"}
	lunch := []models.LunchAssignment{{Student: changed, Bucket: "飯"}}
	if err := a.SaveDailyAssignments(today, []models.Student{changed}, lunch, false); err != nil {
		t.Fatal(err)
	}
	fresh, err := a.GetDailyOverview(today)
	if err != nil || len(fresh.ContactItems) != 0 || fresh.MorningReview != "" {
		t.Fatal("new day inherited previous contact book", err)
	}
	restarted, err := services.NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	restartRecord, err := restarted.GetDailyRecord(today)
	if err != nil || len(restartRecord.Duty.DutyStudents) != 1 || restartRecord.Duty.DutyStudents[0].Name != "替換同學" {
		t.Fatal("restart lost replacement", err)
	}
	if err := a.SaveContactBook(today, []string{" 今天事項 ", ""}, " 複習考 ", " 明日國語複習 "); err != nil {
		t.Fatal(err)
	}
	dh.SaveStudents([]models.Student{{SeatNumber: 1, Name: "後來改名"}})
	current, err := a.GetDailyOverview(today)
	if err != nil || len(current.Duty.DutyStudents) != 1 || current.Duty.DutyStudents[0].Name != "替換同學" || len(current.Duty.LunchAssignments) != 1 {
		t.Fatalf("manual assignment lost: %+v %v", current, err)
	}
	if len(current.ContactItems) != 1 || current.ContactItems[0] != "今天事項" || current.MorningReview != "複習考" {
		t.Fatal("contact normalization failed")
	}
	past, err := a.GetDailyOverview(yesterday)
	if err != nil || past.Students[0].Name != "原始姓名" || past.ContactItems[0] != "昨天事項" {
		t.Fatal("historical record changed", err)
	}
	if err := a.SaveContactBook(yesterday, nil, "", ""); err == nil {
		t.Fatal("history was editable")
	}
	if err := a.SaveDailyAssignments(yesterday, nil, nil, true); err == nil {
		t.Fatal("historical assignments editable")
	}
	if _, err := a.GetDailyOverview(utils.FormatDate(time.Now().AddDate(0, 0, 1))); err == nil {
		t.Fatal("future accepted")
	}
	if err := a.SaveDailyAssignments(today, nil, nil, true); err != nil {
		t.Fatal(err)
	}
	current, err = a.GetDailyOverview(today)
	if err != nil || current.ManualDuty || current.ManualLunch {
		t.Fatal("reset failed", err)
	}
	// Restart preserves the saved contact book.
	dh2, err := services.NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	saved, err := dh2.GetDailyRecord(today)
	if err != nil || saved.ContactItems[0] != "今天事項" {
		t.Fatal("restart lost contact book", err)
	}
}

func TestPreviousCalendarDayReviewIsCopiedAndFrozen(t *testing.T) {
	dh, err := services.NewDataHandlerWithPath(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	sm := services.NewStudentManager(dh)
	a := &App{dh: dh, sm: sm, dc: services.NewDutyCalculator(dh, sm)}
	now := time.Now()
	today := utils.FormatDate(now)
	yesterday := utils.FormatDate(now.AddDate(0, 0, -1))
	oldest := utils.FormatDate(now.AddDate(0, 0, -6))
	expiredSource := utils.FormatDate(now.AddDate(0, 0, -7))
	if err := dh.UpdateDailyData(func(data *services.DailyData) error {
		data.Records[yesterday] = services.DailyRecord{Date: yesterday, NextDayReview: "數學第四單元"}
		data.Records[oldest] = services.DailyRecord{Date: oldest, ReviewFromYesterday: "當時的複習考"}
		data.Records[expiredSource] = services.DailyRecord{Date: expiredSource, NextDayReview: "當時的複習考"}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	current, err := a.GetDailyOverview(today)
	if err != nil || current.ReviewFromYesterday != "數學第四單元" {
		t.Fatalf("yesterday's announcement missing: %+v %v", current, err)
	}
	if current.NextDayReview != "" {
		t.Fatal("new contact book inherited yesterday's announcement")
	}
	if err := a.SaveContactBook(today, nil, "早上複習考", " 明天國語 "); err != nil {
		t.Fatal(err)
	}
	current, err = a.GetDailyOverview(today)
	if err != nil || current.NextDayReview != "明天國語" || current.ReviewFromYesterday != "數學第四單元" {
		t.Fatal("tomorrow's review overwrote today's reminder", err)
	}
	frozen, err := a.GetDailyOverview(oldest)
	if err != nil || frozen.ReviewFromYesterday != "當時的複習考" {
		t.Fatal("history lost its reminder when source expired", err)
	}
	// An empty or missing immediately preceding date must not pick an older review.
	if err := dh.UpdateDailyData(func(data *services.DailyData) error {
		record := data.Records[yesterday]
		record.NextDayReview = ""
		data.Records[yesterday] = record
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	current, err = a.GetDailyOverview(today)
	if err != nil || current.ReviewFromYesterday != "" {
		t.Fatal("empty review displayed", err)
	}
	if err := dh.UpdateDailyData(func(data *services.DailyData) error { delete(data.Records, yesterday); return nil }); err != nil {
		t.Fatal(err)
	}
	current, err = a.GetDailyOverview(today)
	if err != nil || current.ReviewFromYesterday != "" {
		t.Fatal("missing yesterday should have no reminder", err)
	}
}
