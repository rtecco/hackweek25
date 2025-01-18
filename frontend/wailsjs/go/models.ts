export namespace main {
	
	export class Profile {
	    ID: number;
	    Name: string;
	    ProfileSVG: string;
	    Summary: string;
	    SpecialtyTagsRaw: string;
	    SpecialtyTags: string[];
	    Vouches: number;
	    PortfolioDescriptor: string;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.ProfileSVG = source["ProfileSVG"];
	        this.Summary = source["Summary"];
	        this.SpecialtyTagsRaw = source["SpecialtyTagsRaw"];
	        this.SpecialtyTags = source["SpecialtyTags"];
	        this.Vouches = source["Vouches"];
	        this.PortfolioDescriptor = source["PortfolioDescriptor"];
	    }
	}
	export class BuyerChatResponse {
	    followup: string;
	    profiles: Profile[];
	
	    static createFrom(source: any = {}) {
	        return new BuyerChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.followup = source["followup"];
	        this.profiles = this.convertValues(source["profiles"], Profile);
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

