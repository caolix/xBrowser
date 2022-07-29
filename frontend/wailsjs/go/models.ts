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
	
	export class ObjectHandlerResult {
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ObjectHandlerResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.err = source["err"];
	    }
	}
	export class DeleteKey {
	    key: string;
	    keyType: string;
	
	    static createFrom(source: any = {}) {
	        return new DeleteKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.keyType = source["keyType"];
	    }
	}
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
	export class Object {
	    key: string;
	    size: number;
	    human_size: string;
	    // Go type: time.Time
	    last_modified: any;
	
	    static createFrom(source: any = {}) {
	        return new Object(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.human_size = source["human_size"];
	        this.last_modified = this.convertValues(source["last_modified"], null);
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
	export class ListObjectResult {
	    contents: Object[];
	    prefixes: string[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ListObjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.contents = this.convertValues(source["contents"], Object);
	        this.prefixes = source["prefixes"];
	        this.err = source["err"];
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
	export class SelectedFile {
	    key: string;
	    size: number;
	    human_size: string;
	    // Go type: time.Time
	    last_modified: any;
	    source: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectedFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.human_size = source["human_size"];
	        this.last_modified = this.convertValues(source["last_modified"], null);
	        this.source = source["source"];
	        this.name = source["name"];
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
	export class SelectFilesResult {
	    files: SelectedFile[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectFilesResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], SelectedFile);
	        this.err = source["err"];
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

