package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"teacher-wails/internal/models"
	"teacher-wails/internal/utils"
	"testing"
	"time"
)

func TestDailyMigrationPreservesHomeworkAndDoesNotReimport(t *testing.T) {
	dir := t.TempDir()
	legacy := models.Config{MissingHomework: []models.MissingHomework{{Subject: "數學", Students: []int{3}}}}
	raw, _ := json.Marshal(legacy)
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	dh, err := NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	records, err := dh.GetMissingHomework()
	if err != nil || len(records) != 1 || records[0].Subject != "數學" {
		t.Fatalf("migration: %v %v", records, err)
	}
	backup, err := os.ReadFile(path + ".before-daily-split.bak")
	if err != nil || string(backup) != string(raw) {
		t.Fatalf("backup: %v", err)
	}
	settings, _ := dh.GetSettings()
	if err := dh.SaveSettings(settings); err != nil {
		t.Fatal(err)
	}
	saved, _ := os.ReadFile(path)
	var fields map[string]json.RawMessage
	json.Unmarshal(saved, &fields)
	if _, exists := fields["missing_homework"]; exists {
		t.Fatal("homework leaked back into config")
	}
	if err := dh.SaveMissingHomework([]models.MissingHomework{}); err != nil {
		t.Fatal(err)
	}
	dh, err = NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	records, err = dh.GetMissingHomework()
	if err != nil || len(records) != 0 {
		t.Fatal("deleted homework reimported", err)
	}
}

func TestDailyRetentionUsesSevenCalendarDays(t *testing.T) {
	dh, err := NewDataHandlerWithPath(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	err = dh.UpdateDailyData(func(data *DailyData) error {
		for i := 0; i < 10; i++ {
			date := utils.FormatDate(now.AddDate(0, 0, -i))
			data.Records[date] = DailyRecord{Date: date, MorningReview: "原始內容"}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		date := utils.FormatDate(now.AddDate(0, 0, -i))
		record, err := dh.GetDailyRecord(date)
		if i < 7 && (err != nil || record.MorningReview != "原始內容") {
			t.Fatalf("retained %s: %v", date, err)
		}
		if i >= 7 && err == nil {
			t.Fatalf("expired %s retained", date)
		}
	}
}

func TestMalformedDailyFileIsNotOverwritten(t *testing.T) {
	dir := t.TempDir()
	dh, err := NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "data.json")
	raw := []byte(`{"broken":`)
	os.WriteFile(path, raw, 0600)
	if err := dh.SaveMissingHomework(nil); err == nil {
		t.Fatal("expected invalid JSON failure")
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(raw) {
		t.Fatal("corrupt source overwritten")
	}
}

// An existing daily file from an earlier version has no root homework field.
// Configuration migration must retain the legacy working list until it is copied.
func TestMigrationWithExistingSnapshots(t *testing.T) {
	dir := t.TempDir()
	legacy := `{"students":[],"settings":{},"missing_homework":[{"subject":"國語","students":[1]}]}`
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	date := utils.FormatDate(time.Now())
	raw, _ := json.Marshal(DailyData{Version: 1, Records: map[string]DailyRecord{date: {Date: date, MorningReview: "既有快照"}}})
	if err := os.WriteFile(filepath.Join(dir, "data.json"), raw, 0600); err != nil {
		t.Fatal(err)
	}
	dh, err := NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	homework, err := dh.GetMissingHomework()
	if err != nil || len(homework) != 1 {
		t.Fatal("lost legacy homework", err)
	}
	record, err := dh.GetDailyRecord(date)
	if err != nil || record.MorningReview != "既有快照" {
		t.Fatal("lost existing snapshot", err)
	}
}
