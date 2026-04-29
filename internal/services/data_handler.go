package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"teacher-wails/internal/models"
)

// DataHandler manages JSON persistence of the config file.
type DataHandler struct {
	mu         sync.Mutex
	configFile string
}

// NewDataHandler creates a DataHandler with the config file path resolved
// relative to the executable directory.
func NewDataHandler() (*DataHandler, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(filepath.Dir(exe), "data")
	dh := &DataHandler{
		configFile: filepath.Join(dataDir, "config.json"),
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	if err := dh.initConfig(); err != nil {
		return nil, err
	}
	return dh, nil
}

// NewDataHandlerWithPath creates a DataHandler with a custom data directory
// (useful for development where exe path may not be appropriate).
func NewDataHandlerWithPath(dataDir string) (*DataHandler, error) {
	dh := &DataHandler{
		configFile: filepath.Join(dataDir, "config.json"),
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	if err := dh.initConfig(); err != nil {
		return nil, err
	}
	return dh, nil
}

func (dh *DataHandler) initConfig() error {
	if _, err := os.Stat(dh.configFile); os.IsNotExist(err) {
		defaultConfig := models.Config{
			Students: []models.Student{},
			Settings: models.Settings{
				SemesterStartDate: "",
				DutyGroupSize:     2,
				LunchGroupSize:    5,
				DutyStartNumber:   1,
				LunchStartNumber:  1,
				AutoStart:         false,
				MealBuckets:       []string{"飯桶", "菜桶1", "菜桶2", "湯桶", "餐具"},
				CountdownTimes:    []string{},
			},
			Holidays: []string{},
		}
		return dh.SaveConfig(defaultConfig)
	}
	// Migrate legacy data
	config, err := dh.LoadConfig()
	if err != nil {
		return err
	}
	if dh.migrateConfig(&config) {
		return dh.SaveConfig(config)
	}
	return nil
}

// migrateConfig handles legacy field migration. Returns true if changes were made.
func (dh *DataHandler) migrateConfig(config *models.Config) bool {
	// The legacy "enabled" → "duty_enabled"/"lunch_enabled" migration is handled
	// via a raw JSON parse since the Go struct only has the new fields.
	raw, err := os.ReadFile(dh.configFile)
	if err != nil {
		return false
	}
	var rawConfig struct {
		Students []map[string]interface{} `json:"students"`
		Settings map[string]interface{}   `json:"settings"`
	}
	if err := json.Unmarshal(raw, &rawConfig); err != nil {
		return false
	}

	updated := false
	for i, s := range rawConfig.Students {
		if _, hasOld := s["enabled"]; hasOld {
			if _, hasNew := s["duty_enabled"]; !hasNew {
				enabled, _ := s["enabled"].(bool)
				config.Students[i].DutyEnabled = enabled
				config.Students[i].LunchEnabled = enabled
				updated = true
			}
		}
	}

	if rawConfig.Settings != nil {
		if _, ok := rawConfig.Settings["duty_start_number"]; !ok {
			config.Settings.DutyStartNumber = 1
			config.Settings.LunchStartNumber = 1
			updated = true
		}
	}

	// Migrate 7-period timetable → 8-period (insert 午休 at index 4)
	var rawTT struct {
		Timetable [][]string `json:"timetable"`
	}
	if err := json.Unmarshal(raw, &rawTT); err == nil && len(rawTT.Timetable) == 5 {
		if len(rawTT.Timetable[0]) == 7 {
			for d := 0; d < 5; d++ {
				old := rawTT.Timetable[d]
				config.Timetable[d] = [8]string{old[0], old[1], old[2], old[3], "", old[4], old[5], old[6]}
			}
			updated = true
		}
	}

	return updated
}

// LoadConfig reads the config file from disk.
func (dh *DataHandler) LoadConfig() (models.Config, error) {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	return dh.loadConfigUnsafe()
}

func (dh *DataHandler) loadConfigUnsafe() (models.Config, error) {
	var config models.Config
	data, err := os.ReadFile(dh.configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	// Ensure slices are non-nil for clean JSON output
	if config.Students == nil {
		config.Students = []models.Student{}
	}
	if config.Holidays == nil {
		config.Holidays = []string{}
	}
	if config.Settings.MealBuckets == nil {
		config.Settings.MealBuckets = []string{}
	}
	if config.Settings.CountdownTimes == nil {
		config.Settings.CountdownTimes = []string{}
	}
	if config.Settings.PeriodTimes == nil {
		config.Settings.PeriodTimes = []string{}
	}
	return config, err
}

// SaveConfig writes the config to disk with 2-space indent, UTF-8.
func (dh *DataHandler) SaveConfig(config models.Config) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	return dh.saveConfigUnsafe(config)
}

func (dh *DataHandler) saveConfigUnsafe(config models.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dh.configFile, data, 0644)
}

// GetStudents returns the student list from config.
func (dh *DataHandler) GetStudents() ([]models.Student, error) {
	config, err := dh.LoadConfig()
	if err != nil {
		return nil, err
	}
	return config.Students, nil
}

// SaveStudents saves the student list to config.
func (dh *DataHandler) SaveStudents(students []models.Student) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	config, err := dh.loadConfigUnsafe()
	if err != nil {
		return err
	}
	config.Students = students
	return dh.saveConfigUnsafe(config)
}

// GetSettings returns settings from config.
func (dh *DataHandler) GetSettings() (models.Settings, error) {
	config, err := dh.LoadConfig()
	if err != nil {
		return models.Settings{}, err
	}
	return config.Settings, nil
}

// SaveSettings saves settings to config.
func (dh *DataHandler) SaveSettings(settings models.Settings) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	config, err := dh.loadConfigUnsafe()
	if err != nil {
		return err
	}
	config.Settings = settings
	return dh.saveConfigUnsafe(config)
}

// GetHolidays returns the holiday list.
func (dh *DataHandler) GetHolidays() ([]string, error) {
	config, err := dh.LoadConfig()
	if err != nil {
		return nil, err
	}
	return config.Holidays, nil
}

// SaveHolidays saves the holiday list.
func (dh *DataHandler) SaveHolidays(holidays []string) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	config, err := dh.loadConfigUnsafe()
	if err != nil {
		return err
	}
	config.Holidays = holidays
	return dh.saveConfigUnsafe(config)
}

// SaveTimetable saves the timetable to config.
func (dh *DataHandler) SaveTimetable(timetable [5][8]string) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	config, err := dh.loadConfigUnsafe()
	if err != nil {
		return err
	}
	config.Timetable = timetable
	return dh.saveConfigUnsafe(config)
}
