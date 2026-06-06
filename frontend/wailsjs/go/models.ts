export namespace mail {
	
	export class Account {
	    email: string;
	    provider: string;
	    status: string;
	    syncedAt: string;
	    scopes: string[];
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.provider = source["provider"];
	        this.status = source["status"];
	        this.syncedAt = source["syncedAt"];
	        this.scopes = source["scopes"];
	    }
	}
	export class Activity {
	    time: string;
	    text: string;
	    cat?: string;
	
	    static createFrom(source: any = {}) {
	        return new Activity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.text = source["text"];
	        this.cat = source["cat"];
	    }
	}
	export class Fact {
	    key: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Fact(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}
	export class AgentAnalysis {
	    summary: string;
	    facts: Fact[];
	    tasks: string[];
	
	    static createFrom(source: any = {}) {
	        return new AgentAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.summary = source["summary"];
	        this.facts = this.convertValues(source["facts"], Fact);
	        this.tasks = source["tasks"];
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
	export class Category {
	    key: string;
	    label: string;
	    short: string;
	    color: string;
	    desc: string;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.short = source["short"];
	        this.color = source["color"];
	        this.desc = source["desc"];
	    }
	}
	export class DraftLine {
	    text: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new DraftLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.kind = source["kind"];
	    }
	}
	export class Draft {
	    tone: string;
	    lines: DraftLine[];
	
	    static createFrom(source: any = {}) {
	        return new Draft(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tone = source["tone"];
	        this.lines = this.convertValues(source["lines"], DraftLine);
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
	
	
	export class Mail {
	    id: string;
	    cat: string;
	    from: string;
	    fromAddr: string;
	    avatar: string;
	    subject: string;
	    snippet: string;
	    time: string;
	    day: string;
	    unread: boolean;
	    deadline?: string;
	    confidence: number;
	    threadCount: number;
	    body: string[];
	    agent: AgentAnalysis;
	    draft?: Draft;
	
	    static createFrom(source: any = {}) {
	        return new Mail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.cat = source["cat"];
	        this.from = source["from"];
	        this.fromAddr = source["fromAddr"];
	        this.avatar = source["avatar"];
	        this.subject = source["subject"];
	        this.snippet = source["snippet"];
	        this.time = source["time"];
	        this.day = source["day"];
	        this.unread = source["unread"];
	        this.deadline = source["deadline"];
	        this.confidence = source["confidence"];
	        this.threadCount = source["threadCount"];
	        this.body = source["body"];
	        this.agent = this.convertValues(source["agent"], AgentAnalysis);
	        this.draft = this.convertValues(source["draft"], Draft);
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
	export class Rule {
	    on: boolean;
	    text: string;
	    scope: string;
	    locked?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Rule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.on = source["on"];
	        this.text = source["text"];
	        this.scope = source["scope"];
	        this.locked = source["locked"];
	    }
	}

}

