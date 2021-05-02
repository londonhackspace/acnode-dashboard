import APIClient, {NodeRecord} from "./apiclient/dashapi"

import { SignalDispatcher, SimpleEventDispatcher } from "strongly-typed-events"
import ExtendedNodeRecord from "./extendednoderecord";



export default class NodeDataSource {
    private api : APIClient;
    private _activeRow : string;

    private rowChangeSig = new SimpleEventDispatcher<string>();
    private dataChangeSig = new SignalDispatcher();
    private nodeData = new Map<string, ExtendedNodeRecord>();

    private ws : WebSocket;
    private wsRetryMultiplier : number;

    private recalculateTimer : number;
    private running : boolean;

    constructor(api: APIClient) {
        this.api = api;
        this._activeRow = "";
        this.running = false;
    }

    private connect() {
        if(!this.running) {
            return;
        }
        let wsurl : string;
        if(window.location.protocol === 'https:') {
            wsurl = "wss://"+window.location.host+"/ws";
        } else {
            wsurl = "ws://"+window.location.host+"/ws";
        }
        this.ws = new WebSocket(wsurl);
        this.ws.onopen = (evt : Event) => {
            console.log("Opened websocket connection");
            // once we're connected, do a full refresh so we don't miss anything
            this.refresh();
        }
        this.ws.onmessage = (evt: MessageEvent) => {
            // on a message, we can reset the retry backoff
            this.wsRetryMultiplier = 1;
            let node : NodeRecord =  JSON.parse(evt.data);
            this.nodeData.set(node.mqttName, new ExtendedNodeRecord(node));
            this.dataChangeSig.dispatch();
            console.log("Got updated status for node " + node.mqttName);
        }
        this.ws.onerror = (evt: Event) => {
            console.log("Websocket error")
            // try to make a normal refresh anyway
            // this might well fix the error if it was an auth error!
            this.refresh()
        }
        this.ws.onclose = (evt: CloseEvent) => {
            console.log("Websocket closed: " + evt.reason + " - " + evt.wasClean ? "clean" : "unclean")
            if(!evt.wasClean) {
                setTimeout(() => {
                    this.connect()
                }, 1000*this.wsRetryMultiplier);
                this.wsRetryMultiplier++;
                if(this.wsRetryMultiplier > 10) {
                    this.wsRetryMultiplier = 10;
                }
            }
        }
    }

    start() {
        if(this.running) {
            this.stop();
        }
        console.log("Starting NodeDataSource");
        this.running = true;
        this.connect()
        this.recalculateTimer = window.setInterval(() => {
            let changed = false;
            this.nodeData.forEach((node, name) => {
                let previousHealth = node.health;
                let thisChanged = node.refreshObjectHealth();
                changed = changed || thisChanged;
                if(thisChanged) {
                   console.log("Node " + node.name + " changed from " + previousHealth + " to " + node.health);
                }
            });
            if(changed) {
                this.dataChangeSig.dispatch();
            }
        }, 5000);
    }

    stop() {
        console.log("Stopping NodeDataSource");
        this.running = false;
        if(this.ws) {
            console.log("Closing websocket connection");
            this.ws.close()
            this.ws = null;
        }
        window.clearInterval(this.recalculateTimer);
        this.recalculateTimer = -1;
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