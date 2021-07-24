import { SimpleEventDispatcher } from "strongly-typed-events"
import {restGetRequest} from "../utils";

export interface PrinterStatus {
    mqttConnected: boolean;
    octoprintConnected: boolean;
    firmwareVersion: string;
    zHeight: number;
    piUndervoltage: boolean;
    piOverheat: boolean;
    hotendTemperature: number;
    bedTemperature: number;
}

export interface NodeRecord {
    id: number;
    name: string;
    mqttName: string;
    nodeType: string;
    inService: boolean;
    LastSeen: number;
    LastSeenAPI: number;
    LastSeenMQTT: number;
    LastStarted: number;
    MemFree: number;
    MemUsed: number;
    Status: string;
    Version: string;
    CameraId : number | undefined;
    IsTransient: boolean;
    SettingsVersion: number | undefined;
    EEPROMSettingsVersion: number | undefined;
    ResetCause: string | undefined;
    PrinterStatus : PrinterStatus | null;
}

export interface NodeProps {
    CameraId : number | undefined;
    IsTransient: boolean | undefined;
}

export interface User {
    username: string;
    name: string;
    admin: boolean;
}

export interface AccessControlEntry {
    timestamp : number;
    user_id : string;
    user_name : string;
    user_card : string;
    success : boolean;
}

export interface AccessControl {
    count : number;
    page : number;
    pageCount : number;
    entries : AccessControlEntry[];
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
        return restGetRequest(url);
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
            console.log("Auth error");
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

    getAccessControlNodes() : Promise<string[]> {
        return this.makeRequest("accesslogs").then((res : string) => {
            return JSON.parse(res);
        }, this.handleErrorCode.bind(this));
    }

    getAccessControlForNode(node : string, page : number) : Promise<AccessControl> {
        return this.makeRequest("accesslogs/" + node + "?page="+page).then((res : string) => {
            return JSON.parse(res);
        }, this.handleErrorCode.bind(this));
    }

    login(username : string, password: string) : Promise<boolean> {
        return this.makePostRequest("auth/login", JSON.stringify({
            username: username,
            password: password,
        })).then((ret : number) =>{
            this._loginRequired = false;
            this.loginRequiredSignal.dispatchAsync(false);
            return ret == 204;
        }, () => false);
    }

    logout() {
        this.makeRequest("auth/logout").then(() => {
            this._loginRequired = true;
            this.loginRequiredSignal.dispatchAsync(true);
        });
    }

    getUser() : Promise<User> {
        return this.makeRequest("auth/currentuser").then((res : string) => {
            return JSON.parse(res);
        }, this.handleErrorCode.bind(this));
    }

    public setProps(node : string, props : NodeProps) {
        JSON.stringify(props)
        this.makePostRequest("nodes/setProps/" + node, JSON.stringify(props))
    }

};
