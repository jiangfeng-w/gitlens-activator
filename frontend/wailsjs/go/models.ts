export namespace main {
	
	export class actionResult {
	    dirPath: string;
	    ok: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new actionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirPath = source["dirPath"];
	        this.ok = source["ok"];
	        this.message = source["message"];
	    }
	}
	export class gitlensExtension {
	    dirName: string;
	    dirPath: string;
	    version: string;
	    universal: boolean;
	    hasBackup: boolean;
	    activated: boolean;
	
	    static createFrom(source: any = {}) {
	        return new gitlensExtension(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirName = source["dirName"];
	        this.dirPath = source["dirPath"];
	        this.version = source["version"];
	        this.universal = source["universal"];
	        this.hasBackup = source["hasBackup"];
	        this.activated = source["activated"];
	    }
	}
	export class editorCandidate {
	    key: string;
	    name: string;
	    extensionsDir: string;
	    installed: boolean;
	    custom: boolean;
	    extensions?: gitlensExtension[];
	
	    static createFrom(source: any = {}) {
	        return new editorCandidate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.name = source["name"];
	        this.extensionsDir = source["extensionsDir"];
	        this.installed = source["installed"];
	        this.custom = source["custom"];
	        this.extensions = this.convertValues(source["extensions"], gitlensExtension);
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
	export class detectResult {
	    presets: editorCandidate[];
	    customs: editorCandidate[];
	
	    static createFrom(source: any = {}) {
	        return new detectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.presets = this.convertValues(source["presets"], editorCandidate);
	        this.customs = this.convertValues(source["customs"], editorCandidate);
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

