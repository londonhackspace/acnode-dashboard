
export interface NodeRecord {
    id: number;
    name: string;
    mqttName: string;
    nodeType: string;
    LastSeen: number;
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

    constructor(baseurl : string) {
        this.baseurl = baseurl;
    }

    private makeRequest(req : string) : Promise<string> {
        let url : string = this.baseurl+req;

        let p = new Promise<string>((resolve, reject) => {
            let httpReq = new XMLHttpRequest();
            httpReq.open("GET", url, true);
            httpReq.withCredentials = true;
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

    getNodes() : Promise<string[]> {
        return this.makeRequest("nodes").then((res) => {
           return JSON.parse(res);
        });
    }

    getNode(name : string) : Promise<NodeRecord> {
        return this.makeRequest("nodes/"+name).then((res : string) => {
            return JSON.parse(res);
        });
    }
};
