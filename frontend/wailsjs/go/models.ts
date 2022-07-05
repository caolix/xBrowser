export namespace db {
	
	export class LoginInfo {
	    endpoint: string;
	    ak: string;
	    sk: string;
	    remark: string;
	    prepath: string;
	    // Go type: time.Time
	    loginTime: any;
	
	    static createFrom(source: any = {}) {
	        return new LoginInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpoint = source["endpoint"];
	        this.ak = source["ak"];
	        this.sk = source["sk"];
	        this.remark = source["remark"];
	        this.prepath = source["prepath"];
	        this.loginTime = this.convertValues(source["loginTime"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

export namespace app {
	
	export class ListBucketResult {
	    buckets: string[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ListBucketResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.buckets = source["buckets"];
	        this.err = source["err"];
	    }
	}

}

