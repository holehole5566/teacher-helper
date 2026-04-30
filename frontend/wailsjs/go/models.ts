export namespace models {
	
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
	    countdown_music: string;
	    countdown_volume: number;
	
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
	        this.countdown_music = source["countdown_music"];
	        this.countdown_volume = source["countdown_volume"];
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

