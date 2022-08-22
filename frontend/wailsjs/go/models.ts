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
	    humanSize: string;
	    // Go type: time.Time
	    lastModified: any;
	
	    static createFrom(source: any = {}) {
	        return new Object(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.humanSize = source["humanSize"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
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
	    nextMarker: string;
	    isTruncated: boolean;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ListObjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.contents = this.convertValues(source["contents"], Object);
	        this.prefixes = source["prefixes"];
	        this.nextMarker = source["nextMarker"];
	        this.isTruncated = source["isTruncated"];
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
	export class SelectDownloadPathResult {
	    path: string;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectDownloadPathResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.err = source["err"];
	    }
	}
	export class SelectedUploadFile {
	    key: string;
	    size: number;
	    humanSize: string;
	    // Go type: time.Time
	    lastModified: any;
	    type: string;
	    source: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectedUploadFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.humanSize = source["humanSize"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.type = source["type"];
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
	export class SelectUploadFilesResult {
	    files: SelectedUploadFile[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectUploadFilesResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], SelectedUploadFile);
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
	export class SelectUploadFolderResult {
	    path: string;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectUploadFolderResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.err = source["err"];
	    }
	}

}

export namespace db {
	
	export class LoginInfo {
	    accountId: string;
	    endpoint: string;
	    ak: string;
	    sk: string;
	    remark: string;
	    prepath: string;
	    // Go type: time.Time
	    loginTime: any;
	    useSSL: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LoginInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accountId = source["accountId"];
	        this.endpoint = source["endpoint"];
	        this.ak = source["ak"];
	        this.sk = source["sk"];
	        this.remark = source["remark"];
	        this.prepath = source["prepath"];
	        this.loginTime = this.convertValues(source["loginTime"], null);
	        this.useSSL = source["useSSL"];
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
	export class CompletedPart {
	    etag: string;
	    partNumber: number;
	
	    static createFrom(source: any = {}) {
	        return new CompletedPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.etag = source["etag"];
	        this.partNumber = source["partNumber"];
	    }
	}
	export class UploadTask {
	    accountId: string;
	    taskId: string;
	    bucket: string;
	    key: string;
	    name: string;
	    source: string;
	    size: number;
	    humanSize: string;
	    uploadedSize: number;
	    uploadId: string;
	    isMultipart: boolean;
	    partSize: number;
	    status: number;
	    completedPart: CompletedPart[];
	    // Go type: time.Time
	    modifiedTime: any;
	
	    static createFrom(source: any = {}) {
	        return new UploadTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accountId = source["accountId"];
	        this.taskId = source["taskId"];
	        this.bucket = source["bucket"];
	        this.key = source["key"];
	        this.name = source["name"];
	        this.source = source["source"];
	        this.size = source["size"];
	        this.humanSize = source["humanSize"];
	        this.uploadedSize = source["uploadedSize"];
	        this.uploadId = source["uploadId"];
	        this.isMultipart = source["isMultipart"];
	        this.partSize = source["partSize"];
	        this.status = source["status"];
	        this.completedPart = this.convertValues(source["completedPart"], CompletedPart);
	        this.modifiedTime = this.convertValues(source["modifiedTime"], null);
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

