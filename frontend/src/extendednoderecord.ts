import { NodeRecord } from "./apiclient/dashapi";


export enum NodeHealth {
    BAD,
    MEH,
    GOOD,
    UNKNOWN,
};

// Basic idea is to take a NodeRecord and add extra derived values
export default class ExtendedNodeRecord implements NodeRecord {

    // properties
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


    constructor(noderec : NodeRecord) {
        Object.assign(this, noderec)
    }

    get MemTotal() : number {
        return this.MemFree+this.MemUsed;
    }

    get objectHealth() : NodeHealth {
        // basic idea: Start at good, degrade if needed
        let health = NodeHealth.GOOD;

        if(this.LastSeen > 60) {
            health = NodeHealth.MEH;
        }

        if(this.LastSeen > 600) {
            return NodeHealth.BAD;
        }

        if(this.LastSeen == -1) {
            return NodeHealth.UNKNOWN;
        }

        // lower the health if the node watchdog'd recently
        if(this.LastStarted > 0 && this.LastStarted < 600) {
            if(this.ResetCause == "Watchdog") {
                health = NodeHealth.MEH;
            }
        }

        // low on memory?
        if(this.MemUsed > 0 && this.MemFree < (this.MemTotal/10)) {
            health = NodeHealth.BAD;
        }

        return health;
    }
}