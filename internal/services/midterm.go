package services

import (
	"fmt"
	"strings"
	"teacher-wails/internal/models"
	"teacher-wails/internal/utils"
)

func (dh *DataHandler) GetMidtermExam(date string) (models.MidtermExam, error) {
	if _, err := utils.ParseDate(date); err != nil {
		return models.MidtermExam{}, fmt.Errorf("日期格式須為 YYYY-MM-DD")
	}
	config, err := dh.LoadConfig()
	if err != nil {
		return models.MidtermExam{}, err
	}
	record, exists := config.MidtermExams[date]
	if !exists {
		record.Date = date
	}
	return record, nil
}

func (dh *DataHandler) SaveMidtermExam(record models.MidtermExam) error {
	if _, err := utils.ParseDate(record.Date); err != nil {
		return fmt.Errorf("日期格式須為 YYYY-MM-DD")
	}
	for _, count := range []*int{record.Expected, record.Actual} {
		if count != nil && *count < 0 {
			return fmt.Errorf("人數不可小於 0")
		}
	}
	if record.Expected != nil && record.Actual != nil && *record.Actual > *record.Expected {
		return fmt.Errorf("實到不可大於應到")
	}
	for i := range record.Subjects {
		record.Subjects[i] = strings.TrimSpace(record.Subjects[i])
	}
	dh.mu.Lock()
	defer dh.mu.Unlock()
	config, err := dh.loadConfigUnsafe()
	if err != nil {
		return err
	}
	if config.MidtermExams == nil {
		config.MidtermExams = map[string]models.MidtermExam{}
	}
	config.MidtermExams[record.Date] = record
	return dh.saveConfigUnsafe(config)
}
