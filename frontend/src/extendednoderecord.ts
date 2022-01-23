import {NodeRecord, PrinterStatus} from "./apiclient/dashapi";

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
    inService: boolean;
    LastSeen: number;
    LastSeenMQTT: number;
    LastSeenAPI: number;
    LastStarted: number;
    MemFree: number;
    MemUsed: number;
    Status: string;
    Version: string;
    CameraId : number | undefined;
    IsTransient: boolean;
    InUse: boolean;
    VersionDate : Date;
    VersionMessage : string;
    SettingsVersion: number | undefined;
    EEPROMSettingsVersion: number | undefined;
    ResetCause: string | undefined;
    PrinterStatus: PrinterStatus | null;

    health : NodeHealth
    _healthHints : string[] = [];
    healthCalculated : number = 0;

    printerHealth : NodeHealth
    _printerHealthHints : string[] = [];

    constructor(noderec : NodeRecord) {
        Object.assign(this, noderec)
        this.health = this.calulateObjectHealth();
        this.printerHealth = this.calculatePrinterHealth();

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

    get printerHealthHints() : string[] {
        return this._printerHealthHints;
    }

    get healthStatusAge() : number {
        return Date.now()-this.healthCalculated;
    }

    refreshObjectHealth() : boolean {
        let newHealth = this.calulateObjectHealth();
        let printerHealth = this.calculatePrinterHealth();
        if(newHealth != this.health || printerHealth != this.printerHealth) {
            this.healthCalculated = Date.now();
            this.health = newHealth;
            this.printerHealth = printerHealth;
            return true;
        }
        return false;
    }

    private calculatePrinterHealth() : NodeHealth {
        this._printerHealthHints = [];
        if(this.PrinterStatus == null) {
            return NodeHealth.UNKNOWN;
        }

        if(!this.PrinterStatus.mqttConnected) {
            this._printerHealthHints.push("Octoprint disconnected from MQTT");
            return NodeHealth.BAD;
        }

        if(this.PrinterStatus.piOverheat) {
            this._printerHealthHints.push("Pi Overheating");
            return NodeHealth.MEH;
        }

        if(this.PrinterStatus.piUndervoltage) {
            this._printerHealthHints.push("Pi Undervoltage");
            return NodeHealth.MEH;
        }

        return NodeHealth.GOOD;
    }

    private calulateObjectHealth() : NodeHealth {
        // basic idea: Start at good, degrade if needed
        let health = NodeHealth.GOOD;
        this._healthHints = [];

        if(!this.inService) {
            this._healthHints.push("Tool out of service");
        }

        // Calculate relative values for decisions
        let lastSeen = (Date.now()/1000) - this.LastSeen;
        let lastSeenMQTT = (Date.now()/1000) - this.LastSeenMQTT;
        let lastSeenAPI = (Date.now()/1000) - this.LastSeenAPI;
        let lastStarted = (Date.now()/1000) - this.LastStarted;

        var apiThreshold = 610;
        var apiThresholdText = "over 10 minutes";

        var mqttThreshold = 130;
        var mqttThresholdText = "over 2 minutes";

        // unrestricted doors don't check in nearly so often if they're not
        // running firmware new enough to periodically revalidate the cache
        // or aren't running a maintainer cache
        if(this.nodeType == "Door") {
            apiThreshold = (3600*12) + 10;
            apiThresholdText = "over 12 hours"
        }

        // use the newer LastSeen values?
        if(this.LastSeenMQTT > -1 || this.LastSeenAPI > -1) {

            // if we're seeing neither MQTT or ACServer log entries,
            // it's probably dead
            if((lastSeenAPI > apiThreshold || this.LastSeenAPI == -1)
                && (lastSeenMQTT > mqttThreshold || this.LastSeenMQTT == -1)) {
                var text  = apiThresholdText;
                if(mqttThreshold > apiThreshold) {
                    text = mqttThresholdText;
                }
                this._healthHints.push("Has not been seen online in any form in " + text);
                return this.IsTransient ? NodeHealth.UNKNOWN : NodeHealth.BAD
            }

            if(this.LastSeenMQTT == -1 || lastSeenMQTT > mqttThreshold) {
                this._healthHints.push("Has not sent a message via MQTT in " + mqttThresholdText);
                health = NodeHealth.MEH;
            }

            if((this.LastSeenAPI == -1 || lastSeenAPI > apiThreshold) && !this.InUse) {
                this._healthHints.push("Has not contacted ACServer in " + apiThresholdText);
                health = NodeHealth.MEH;
            }

        } else {
            if(this.LastSeen == -1) {
                return NodeHealth.UNKNOWN;
            }

            if(lastSeen > 600) {
                this._healthHints.push("Has not been seen online in over 10 minutes");
                health = this.IsTransient ? NodeHealth.UNKNOWN : NodeHealth.BAD;
            } else if(lastSeen > 60) {
                this._healthHints.push("Has not been seen online in over a minute");
                return NodeHealth.MEH;
            }
        }

        // lower the health if the node watchdog'd recently
        if(this.LastStarted > 0 && lastStarted < 600) {
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