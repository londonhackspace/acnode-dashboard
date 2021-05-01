import APIClient, {NodeRecord} from "./apiclient/dashapi"

import { SignalDispatcher, SimpleEventDispatcher } from "strongly-typed-events"
import ExtendedNodeRecord from "./extendednoderecord";



export default class NodeDataSource {
    private api : APIClient;
    private _activeRow : string;

    private rowChangeSig = new SimpleEventDispatcher<string>();
    private dataChangeSig = new SignalDispatcher();
    private refreshTimer : number;
    private started = false;
    private nodeData = new Map<string, ExtendedNodeRecord>();

    constructor(api: APIClient) {
        this.api = api;
        this._activeRow = "";
    }

    start() {
        this.stop();
        console.log("Setting up timer");
        this.refreshTimer = window.setInterval(this.refresh.bind(this), 10000);
        this.started = true;
        this.refresh();
    }

    stop() {
        if(this.started) {
            console.log("Stopping the timer");
            window.clearInterval(this.refreshTimer);
            this.started = false;
        }
    }

    refresh() {
        this.api.getNodes().then((nodes : string[]) => {
            return Promise.all(nodes.map(n => this.api.getNode(n)));
        }).then((results) => {
            for(let res of results) {
                this.nodeData.set(res.mqttName, new ExtendedNodeRecord(res));
            }
            this.dataChangeSig.dispatch();
        });
    }

    get nodes() : string[] {
        return Array.from(this.nodeData.keys());
    }

    getNode(name : string) : ExtendedNodeRecord | null {
        if(this.nodeData.has(name)) {
            return this.nodeData.get(name);
        }
        return null;
    }

    setActiveRow(row : string) {
        this._activeRow = row;
        this.rowChangeSig.dispatch(row);
    }

    public get activeRow() {
        return this._activeRow;
    }

    public get onActiveRowChange() {
        return this.rowChangeSig.asEvent();
    }

    public get onDataChange() {
        return this.dataChangeSig.asEvent();
    }

}