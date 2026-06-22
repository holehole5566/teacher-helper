package models

// Student represents a student with duty and lunch participation flags.
type Student struct {
	SeatNumber   int    `json:"seat_number"`
	Name         string `json:"name"`
	DutyEnabled  bool   `json:"duty_enabled"`
	LunchEnabled bool   `json:"lunch_enabled"`
}

// MusicTrack is a single entry in the global countdown music playlist.
type MusicTrack struct {
	Path     string `json:"path"`
	InRandom bool   `json:"in_random"`
}

// CountdownTimeMusic maps a countdown_time to a playback mode and optional track index.
type CountdownTimeMusic struct {
	Time  string `json:"time"`
	Mode  string `json:"mode"`  // "random" | "index"
	Index int    `json:"index"` // 0-based index into Settings.CountdownMusics
}

// Settings holds all configurable parameters.
type Settings struct {
	SemesterStartDate     string               `json:"semester_start_date"`
	DutyGroupSize         int                  `json:"duty_group_size"`
	LunchGroupSize        int                  `json:"lunch_group_size"`
	DutyStartNumber       int                  `json:"duty_start_number"`
	LunchStartNumber      int                  `json:"lunch_start_number"`
	AutoStart             bool                 `json:"auto_start"`
	MealBuckets           []string             `json:"meal_buckets"`
	CountdownTimes        []string             `json:"countdown_times"`
	PeriodTimes           []string             `json:"period_times"`
	CountdownVolume       float64              `json:"countdown_volume"`
	CountdownMusics       []MusicTrack         `json:"countdown_musics"`
	CountdownTimeMusicMap []CountdownTimeMusic `json:"countdown_time_music_map"`
	DiscordWebhook        string               `json:"discord_webhook"`
	AudioOutputDevice     string               `json:"audio_output_device"`
}

// Config is the top-level JSON structure persisted to config.json.
type Config struct {
	Password        string            `json:"password"`
	Students        []Student         `json:"students"`
	Settings        Settings          `json:"settings"`
	Holidays        []string          `json:"holidays"`
	Timetable       [5][8]string      `json:"timetable"`
	MissingHomework []MissingHomework `json:"missing_homework"`
}

// LunchAssignment pairs a student with a meal bucket.
type LunchAssignment struct {
	Student Student `json:"student"`
	Bucket  string  `json:"bucket"`
}

// MissingHomework tracks one homework assignment's missing submissions.
type MissingHomework struct {
	Subject  string `json:"subject"`
	Students []int  `json:"students"`
	Note     string `json:"note"`
}

// TodayDutyResult is returned by GetTodayDuty to the frontend.
type TodayDutyResult struct {
	Date             string            `json:"date"`
	DisplayDate      string            `json:"displayDate"`
	IsWorkday        bool              `json:"isWorkday"`
	DutyStudents     []Student         `json:"dutyStudents"`
	LunchAssignments []LunchAssignment `json:"lunchAssignments"`
}
