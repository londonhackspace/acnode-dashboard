import { SimpleEventDispatcher } from "strongly-typed-events"

export interface NodeRecord {
    id: number;
    name: string;
    mqttName: string;
    nodeType: string;
    LastSeen: number;
    LastSeenAPI: number;
    LastSeenMQTT: number;
    LastStarted: number;
    MemFree: number;
    MemUsed: number;
    Status: string;
    Version: string;
    SettingsVersion: number | undefined;
    EEPROMSettingsVersion: number | undefined;
    ResetCause: string | undefined;
}

export default class DashAPI {
    private readonly baseurl : string;
    private _loginRequired : boolean;

    private loginRequiredSignal : SimpleEventDispatcher<boolean>

    constructor(baseurl : string) {
        this.baseurl = baseurl;
        this._loginRequired = false;
        this.loginRequiredSignal = new SimpleEventDispatcher<boolean>();
    }

    private makeRequest(req : string) : Promise<string> {
        let url : string = this.baseurl+req;

        let p = new Promise<string>((resolve, reject) => {
            let httpReq = new XMLHttpRequest();
            httpReq.open("GET", url, true);
            httpReq.onreadystatechange = function() {
                if(this.readyState == 4) {
                    if(this.status >= 200 && this.status < 300) {
                        resolve(this.response);
                    } else {
                        reject(this.status);
                    }
                }
            }
            httpReq.send();
        });
        return p;
    }

    private makePostRequest(req : string, body : string) : Promise<number> {
        let url : string = this.baseurl+req;

        let p = new Promise<number>((resolve, reject) => {
            let httpReq = new XMLHttpRequest();
            httpReq.open("POST", url, true);
            httpReq.onreadystatechange = function() {
                if(this.readyState == 4) {
                    if(this.status >= 200 && this.status < 300) {
                        resolve(this.status);
                    } else {
                        reject(this.status);
                    }
                }
            }
            httpReq.send(body);
        });
        return p;
    }

    get loginRequired() : boolean {
        return this._loginRequired;
    }

    get onLoginRequired() {
        return this.loginRequiredSignal.asEvent();
    }

    handleErrorCode(code : number) {
        if(code == 401) {
            this._loginRequired = true;
            this.loginRequiredSignal.dispatchAsync(this._loginRequired);
        }
    }

    getNodes() : Promise<string[]> {
        return this.makeRequest("nodes").then((res) => {
           return JSON.parse(res);
        }, this.handleErrorCode.bind(this));
    }

    getNode(name : string) : Promise<NodeRecord> {
        return this.makeRequest("nodes/"+name).then((res : string) => {
            return JSON.parse(res);
        }, this.handleErrorCode.bind(this));
    }

    login(username : string, password: string) : Promise<boolean> {
        return this.makePostRequest("/auth/login", JSON.stringify({
            username: username,
            password: password,
        })).then((ret : number) =>{
            this._loginRequired = false;
            this.loginRequiredSignal.dispatchAsync(false);
            return ret == 204;
        }, () => false);
    }

    logout() {
        this.makeRequest("/auth/logout").then(() => {
            this._loginRequired = true;
            this.loginRequiredSignal.dispatchAsync(true);
        });
    }
};
