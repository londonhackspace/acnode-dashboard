import APIClient, {NodeProps, NodeRecord} from "./apiclient/dashapi"

import { SignalDispatcher, SimpleEventDispatcher } from "strongly-typed-events"
import ExtendedNodeRecord from "./extendednoderecord";
import GithubCommitInfo from "./githubcommitinfo"

const useWebsockets = true;

interface RowChangeSignal {
    node : string;
    edit : boolean;
}

export default class NodeDataSource {
    private api : APIClient;
    private _activeRow : string;
    private _isEdit : boolean;

    private rowChangeSig = new SimpleEventDispatcher<RowChangeSignal>();
    private dataChangeSig = new SignalDispatcher();
    private nodeData = new Map<string, ExtendedNodeRecord>();

    private ws : WebSocket;
    private wsRetryMultiplier : number;

    private recalculateTimer : number;
    private running : boolean;

    private polltimer : number

    constructor(api: APIClient) {
        this.api = api;
        this._activeRow = "";
        this._isEdit = false;
        this.running = false;
    }

    private connectWS() {
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
            this.polltimer = window.setInterval(() => {
                console.log("Stuffing junk into websocket to keep it alive (hopefully)");
                // send some data once in a while to prevent timeouts
                this.ws.send("dummy");
            },10000);
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
            window.clearInterval(this.polltimer)
            this.polltimer = null
        }
        this.ws.onclose = (evt: CloseEvent) => {
            window.clearInterval(this.polltimer)
            this.polltimer = null
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

    private connect() {
        if(!this.running) {
            return;
        }
        if(useWebsockets) {
            this.connectWS()
        } else {
            console.log("Setting up periodic refresh timer");
            this.polltimer = window.setInterval(() => {
                console.log("Refreshing data");
                this.refresh();
            }, 10000);
            this.refresh();
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
        if(this.polltimer) {
            console.log("Clearing up poll timer");
            window.clearInterval(this.polltimer);
            this.polltimer = null;
        }
        window.clearInterval(this.recalculateTimer);
        this.recalculateTimer = -1;
    }

    refresh() {
        this.api.getNodes().then((nodes : string[]) => {
            return Promise.all(nodes.map(n => this.api.getNode(n)));
        }).then((results) => {
            for(let res of results) {
                if(!res) continue;
                let extendedRec = new ExtendedNodeRecord(res);

                if(extendedRec.Version && extendedRec.Version != "") {
                    GithubCommitInfo.getCommit(extendedRec.Version)
                        .then((res) => {
                            if(res) {
                                extendedRec.VersionDate = new Date(res.commit.committer.date);
                                extendedRec.VersionMessage = res.commit.message;
                            }
                            this.dataChangeSig.dispatch();
                        });
                }

                this.nodeData.set(res.mqttName, extendedRec);
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

    setActiveRow(row : string, edit : boolean) {
        this._activeRow = row;
        this._isEdit = edit;
        this.rowChangeSig.dispatch({node : row, edit: edit});
    }

    public get activeRow() {
        return this._activeRow;
    }

    public get isEdit() {
        return this._isEdit;
    }

    public get onActiveRowChange() {
        return this.rowChangeSig.asEvent();
    }

    public get onDataChange() {
        return this.dataChangeSig.asEvent();
    }

    public setNodeProps(node : string, props : NodeProps) {
        this.api.setProps(node ,props);
    }
}