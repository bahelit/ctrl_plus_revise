export namespace main {
	
	export class AppSettings {
	    replaceHighlightedText: boolean;
	    speakAIResponse: boolean;
	    showPopUp: boolean;
	    showStartWindow: boolean;
	    currentPrompt: number;
	    currentChatModel: string;
	    currentModel: string;
	    currentFromLang: string;
	    currentToLang: string;
	    stopOllamaOnShutDown: boolean;
	    useRemoteOllama: boolean;
	    ollamaURL: string;
	    askAIShortcut: string[];
	    ctrlReviseShortcut: string[];
	    translateShortcut: string[];
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.replaceHighlightedText = source["replaceHighlightedText"];
	        this.speakAIResponse = source["speakAIResponse"];
	        this.showPopUp = source["showPopUp"];
	        this.showStartWindow = source["showStartWindow"];
	        this.currentPrompt = source["currentPrompt"];
	        this.currentChatModel = source["currentChatModel"];
	        this.currentModel = source["currentModel"];
	        this.currentFromLang = source["currentFromLang"];
	        this.currentToLang = source["currentToLang"];
	        this.stopOllamaOnShutDown = source["stopOllamaOnShutDown"];
	        this.useRemoteOllama = source["useRemoteOllama"];
	        this.ollamaURL = source["ollamaURL"];
	        this.askAIShortcut = source["askAIShortcut"];
	        this.ctrlReviseShortcut = source["ctrlReviseShortcut"];
	        this.translateShortcut = source["translateShortcut"];
	    }
	}
	export class ChatInfo {
	    id: number;
	    title: string;
	    model: number;
	    questions: string[];
	    responses: string[];
	
	    static createFrom(source: any = {}) {
	        return new ChatInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.model = source["model"];
	        this.questions = source["questions"];
	        this.responses = source["responses"];
	    }
	}
	export class ChatResponse {
	    response: string;
	    chatId: number;
	    context: number[];
	
	    static createFrom(source: any = {}) {
	        return new ChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.response = source["response"];
	        this.chatId = source["chatId"];
	        this.context = source["context"];
	    }
	}
	export class MealInfo {
	    consumers: string;
	    meal: string;
	    theme: string;
	    cookware: string[];
	    dairy: string[];
	    veggies: string[];
	    protein: string[];
	    freezer: string[];
	    herbs: string[];
	    pantry: string[];
	    allergies: string[];
	
	    static createFrom(source: any = {}) {
	        return new MealInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.consumers = source["consumers"];
	        this.meal = source["meal"];
	        this.theme = source["theme"];
	        this.cookware = source["cookware"];
	        this.dairy = source["dairy"];
	        this.veggies = source["veggies"];
	        this.protein = source["protein"];
	        this.freezer = source["freezer"];
	        this.herbs = source["herbs"];
	        this.pantry = source["pantry"];
	        this.allergies = source["allergies"];
	    }
	}
	export class ModelInfo {
	    name: string;
	    isInstalled: boolean;
	    isRecommended: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.isInstalled = source["isInstalled"];
	        this.isRecommended = source["isRecommended"];
	    }
	}
	export class PromptInfo {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new PromptInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class ShortcutsConfig {
	    askAI: string[];
	    ctrlRevise: string[];
	    translate: string[];
	
	    static createFrom(source: any = {}) {
	        return new ShortcutsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.askAI = source["askAI"];
	        this.ctrlRevise = source["ctrlRevise"];
	        this.translate = source["translate"];
	    }
	}

}

