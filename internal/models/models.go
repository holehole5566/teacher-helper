package models

// Student represents a student with duty and lunch participation flags.
type Student struct {
	SeatNumber   int    `json:"seat_number"`
	Name         string `json:"name"`
	DutyEnabled  bool   `json:"duty_enabled"`
	LunchEnabled bool   `json:"lunch_enabled"`
}

// Settings holds all configurable parameters.
type Settings struct {
	SemesterStartDate string   `json:"semester_start_date"`
	DutyGroupSize     int      `json:"duty_group_size"`
	LunchGroupSize    int      `json:"lunch_group_size"`
	DutyStartNumber   int      `json:"duty_start_number"`
	LunchStartNumber  int      `json:"lunch_start_number"`
	AutoStart         bool     `json:"auto_start"`
	MealBuckets       []string `json:"meal_buckets"`
	CountdownTimes    []string `json:"countdown_times"`
	PeriodTimes       []string `json:"period_times"`
	CountdownMusic    string   `json:"countdown_music"`
	CountdownVolume   float64  `json:"countdown_volume"`
}

// Config is the top-level JSON structure persisted to config.json.
type Config struct {
	Students  []Student    `json:"students"`
	Settings  Settings     `json:"settings"`
	Holidays  []string     `json:"holidays"`
	Timetable [5][8]string `json:"timetable"`
}

// LunchAssignment pairs a student with a meal bucket.
type LunchAssignment struct {
	Student Student `json:"student"`
	Bucket  string  `json:"bucket"`
}

// TodayDutyResult is returned by GetTodayDuty to the frontend.
type TodayDutyResult struct {
	Date             string            `json:"date"`
	DisplayDate      string            `json:"displayDate"`
	IsWorkday        bool              `json:"isWorkday"`
	DutyStudents     []Student         `json:"dutyStudents"`
	LunchAssignments []LunchAssignment `json:"lunchAssignments"`
}
