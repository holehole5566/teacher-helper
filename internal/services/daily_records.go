package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"teacher-wails/internal/models"
	"teacher-wails/internal/utils"
	"time"
)

// DailyRecord preserves the values displayed on a particular calendar date.
type DailyRecord struct {
	Date                string                   `json:"date"`
	UpdatedAt           string                   `json:"updated_at"`
	Duty                models.TodayDutyResult   `json:"duty"`
	Students            []models.Student         `json:"students"`
	Timetable           [5][8]string             `json:"timetable"`
	PeriodTimes         []string                 `json:"period_times"`
	MissingHomework     []models.MissingHomework `json:"missing_homework"`
	ContactItems        []string                 `json:"contact_items"`
	NextDayReview       string                   `json:"next_day_review"`
	ReviewFromYesterday string                   `json:"review_from_yesterday"`
	MorningReview       string                   `json:"morning_review"`
	ManualDuty          bool                     `json:"manual_duty"`
	ManualLunch         bool                     `json:"manual_lunch"`
	LayoutVersion       int                      `json:"layout_version"`
}

type DailyData struct {
	ReviewCarry     *ReviewAnnouncement      `json:"review_carry,omitempty"`
	MissingHomework []models.MissingHomework `json:"missing_homework"`
	Version         int                      `json:"version"`
	Records         map[string]DailyRecord   `json:"records"`
}

// UpdateDailyData serializes daily writes separately from configuration writes.
func (dh *DataHandler) UpdateDailyData(update func(*DailyData) error) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	path := filepath.Join(filepath.Dir(dh.configFile), "data.json")
	data, err := dh.loadDailyUnsafe()
	if err != nil {
		return err
	}
	if err := update(&data); err != nil {
		return err
	}
	config, err := dh.loadConfigUnsafe()
	if err != nil {
		return err
	}
	data.preserveReview(time.Now(), config.Holidays)
	cutoff := utils.FormatDate(time.Now().AddDate(0, 0, -6))
	for date := range data.Records {
		if date < cutoff {
			delete(data.Records, date)
		}
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".daily-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	if _, err = temp.Write(raw); err != nil {
		temp.Close()
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), path)
}

func (dh *DataHandler) loadDailyUnsafe() (DailyData, error) {
	data := DailyData{Version: 1, Records: map[string]DailyRecord{}}
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(dh.configFile), "data.json"))
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return data, err
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return data, err
	}
	if data.Records == nil {
		data.Records = map[string]DailyRecord{}
	}
	return data, nil
}

func (dh *DataHandler) GetDailyRecord(date string) (DailyRecord, error) {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	data, err := dh.loadDailyUnsafe()
	if err != nil {
		return DailyRecord{}, err
	}
	record, exists := data.Records[date]
	if !exists {
		return DailyRecord{}, fmt.Errorf("此日期沒有紀錄")
	}
	return record, nil
}
