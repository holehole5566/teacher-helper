export namespace main {

	export class UpdateResult {
	    has_update: boolean;
	    current_version: string;
	    latest_version: string;
	    release_notes: string;

	    static createFrom(source: any = {}) {
	        return new UpdateResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.has_update = source["has_update"];
	        this.current_version = source["current_version"];
	        this.latest_version = source["latest_version"];
	        this.release_notes = source["release_notes"];
	    }
	}

}

export namespace models {

	export class CountdownTimeMusic {
	    time: string;
	    mode: string;
	    index: number;

	    static createFrom(source: any = {}) {
	        return new CountdownTimeMusic(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.mode = source["mode"];
	        this.index = source["index"];
	    }
	}
	export class Student {
	    seat_number: number;
	    name: string;
	    duty_enabled: boolean;
	    lunch_enabled: boolean;

	    static createFrom(source: any = {}) {
	        return new Student(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.seat_number = source["seat_number"];
	        this.name = source["name"];
	        this.duty_enabled = source["duty_enabled"];
	        this.lunch_enabled = source["lunch_enabled"];
	    }
	}
	export class LunchAssignment {
	    student: Student;
	    bucket: string;

	    static createFrom(source: any = {}) {
	        return new LunchAssignment(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.student = this.convertValues(source["student"], Student);
	        this.bucket = source["bucket"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MidtermExam {
	    date: string;
	    subjects: string[];
	    expected?: number;
	    actual?: number;

	    static createFrom(source: any = {}) {
	        return new MidtermExam(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.subjects = source["subjects"];
	        this.expected = source["expected"];
	        this.actual = source["actual"];
	    }
	}
	export class MissingHomework {
	    subject: string;
	    students: number[];
	    note: string;

	    static createFrom(source: any = {}) {
	        return new MissingHomework(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.subject = source["subject"];
	        this.students = source["students"];
	        this.note = source["note"];
	    }
	}
	export class MusicTrack {
	    path: string;
	    in_random: boolean;

	    static createFrom(source: any = {}) {
	        return new MusicTrack(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.in_random = source["in_random"];
	    }
	}
	export class Settings {
	    semester_start_date: string;
	    duty_group_size: number;
	    lunch_group_size: number;
	    duty_start_number: number;
	    lunch_start_number: number;
	    auto_start: boolean;
	    meal_buckets: string[];
	    countdown_times: string[];
	    period_times: string[];
	    countdown_volume: number;
	    countdown_musics: MusicTrack[];
	    countdown_time_music_map: CountdownTimeMusic[];
	    discord_webhook: string;
	    audio_output_device: string;

	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.semester_start_date = source["semester_start_date"];
	        this.duty_group_size = source["duty_group_size"];
	        this.lunch_group_size = source["lunch_group_size"];
	        this.duty_start_number = source["duty_start_number"];
	        this.lunch_start_number = source["lunch_start_number"];
	        this.auto_start = source["auto_start"];
	        this.meal_buckets = source["meal_buckets"];
	        this.countdown_times = source["countdown_times"];
	        this.period_times = source["period_times"];
	        this.countdown_volume = source["countdown_volume"];
	        this.countdown_musics = this.convertValues(source["countdown_musics"], MusicTrack);
	        this.countdown_time_music_map = this.convertValues(source["countdown_time_music_map"], CountdownTimeMusic);
	        this.discord_webhook = source["discord_webhook"];
	        this.audio_output_device = source["audio_output_device"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

	export class TodayDutyResult {
	    date: string;
	    displayDate: string;
	    isWorkday: boolean;
	    dutyStudents: Student[];
	    lunchAssignments: LunchAssignment[];

	    static createFrom(source: any = {}) {
	        return new TodayDutyResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.displayDate = source["displayDate"];
	        this.isWorkday = source["isWorkday"];
	        this.dutyStudents = this.convertValues(source["dutyStudents"], Student);
	        this.lunchAssignments = this.convertValues(source["lunchAssignments"], LunchAssignment);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace services {

	export class DailyRecord {
	    date: string;
	    updated_at: string;
	    duty: models.TodayDutyResult;
	    students: models.Student[];
	    timetable: string[][];
	    period_times: string[];
	    missing_homework: models.MissingHomework[];
	    contact_items: string[];
	    next_day_review: string;
	    review_from_yesterday: string;
	    morning_review: string;
	    manual_duty: boolean;
	    manual_lunch: boolean;
	    layout_version: number;

	    static createFrom(source: any = {}) {
	        return new DailyRecord(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.updated_at = source["updated_at"];
	        this.duty = this.convertValues(source["duty"], models.TodayDutyResult);
	        this.students = this.convertValues(source["students"], models.Student);
	        this.timetable = source["timetable"];
	        this.period_times = source["period_times"];
	        this.missing_homework = this.convertValues(source["missing_homework"], models.MissingHomework);
	        this.contact_items = source["contact_items"];
	        this.next_day_review = source["next_day_review"];
	        this.review_from_yesterday = source["review_from_yesterday"];
	        this.morning_review = source["morning_review"];
	        this.manual_duty = source["manual_duty"];
	        this.manual_lunch = source["manual_lunch"];
	        this.layout_version = source["layout_version"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

