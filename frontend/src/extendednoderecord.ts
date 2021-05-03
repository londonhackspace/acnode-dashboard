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
    LastSeenMQTT: number;
    LastSeenAPI: number;
    LastStarted: number;
    MemFree: number;
    MemUsed: number;
    Status: string;
    Version: string;
    VersionDate : Date;
    SettingsVersion: number | undefined;
    EEPROMSettingsVersion: number | undefined;
    ResetCause: string | undefined;

    health : NodeHealth
    _healthHints : string[] = [];
    healthCalculated : number = 0;

    constructor(noderec : NodeRecord) {
        Object.assign(this, noderec)
        this.health = this.calulateObjectHealth()

        // null for now - can be filled in later
        this.VersionDate = null;
        // some fudging to make lastseen an aggregate
        if(this.LastSeenAPI > this.LastSeen) {
            this.LastSeen = this.LastSeenAPI;
        }
        if(this.LastSeenMQTT > this.LastSeen) {
            this.LastSeen = this.LastSeenMQTT;
        }
    }

    get MemTotal() : number {
        return this.MemFree+this.MemUsed;
    }

    get objectHealth() : NodeHealth {
        return this.health;
    }

    get healthHints() : string[] {
        return this._healthHints;
    }

    get healthStatusAge() : number {
        return Date.now()-this.healthCalculated;
    }

    refreshObjectHealth() : boolean {
        let newHealth = this.calulateObjectHealth();
        if(newHealth != this.health) {
            this.healthCalculated = Date.now();
            this.health = newHealth;
            return true;
        }
        return false;
    }

    private calulateObjectHealth() : NodeHealth {
        // basic idea: Start at good, degrade if needed
        let health = NodeHealth.GOOD;
        this._healthHints = [];

        let lastSeen = (Date.now()/1000) - this.LastSeen;
        let lastSeenMQTT = (Date.now()/1000) - this.LastSeenMQTT;
        let lastSeenAPI = (Date.now()/1000) - this.LastSeenAPI;

        // use the newer LastSeen values?
        if(this.LastSeenMQTT > -1 || this.LastSeenAPI > -1) {
            if(this.LastSeenAPI == -1 && this.LastSeenMQTT == -1) {
                return NodeHealth.UNKNOWN;
            }

            // if we're seeing neither MQTT or ACServer log entries,
            // it's probably dead
            if((lastSeenAPI > 600 || this.LastSeenAPI == -1)
                && (lastSeenMQTT > 600 || this.LastSeenMQTT == -1)) {
                this._healthHints.push("Has not been seen online in any form in over 10 minutes");
                return NodeHealth.BAD
            }

            if(this.LastSeenMQTT == -1 || lastSeenMQTT > 120) {
                this._healthHints.push("Has not sent a message via MQTT in over 2 minutes");
                health = NodeHealth.MEH;
            }

            if(this.LastSeenAPI == -1 || lastSeenAPI > 600) {
                this._healthHints.push("Has not contacted ACServer in over 10 minutes");
                health = NodeHealth.MEH;
            }

        } else {
            if(this.LastSeen == -1) {
                return NodeHealth.UNKNOWN;
            }

            if(lastSeen > 600) {
                this._healthHints.push("Has not been seen online in over 10 minutes");
                health = NodeHealth.BAD;
            } else if(lastSeen > 60) {
                this._healthHints.push("Has not been seen online in over a minute");
                return NodeHealth.MEH;
            }
        }

        // lower the health if the node watchdog'd recently
        if(this.LastStarted > 0 && this.LastStarted < 600) {
            if(this.ResetCause == "Watchdog") {
                this._healthHints.push("Watchdog reset detected in last 10 minutes");
                health = NodeHealth.MEH;
            }
        }

        // low on memory?
        if(this.MemUsed > 0 && this.MemFree < (this.MemTotal/10)) {
            this._healthHints.push("Very low on memory (<10% left)");
            health = NodeHealth.BAD;
        }

        return health;
    }
}