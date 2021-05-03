import React, {ReactElement} from "react";
import dateformat from "dateformat";

import styles from "./nodetable.module.css";

import NodeDataSource from "../NodeDataSource";
import {NodeRecord} from "../apiclient/dashapi";
import ExtendedNodeRecord, {NodeHealth} from "../extendednoderecord";
import {isUndefined} from "webpack-merge/dist/utils";
import StatusBall, {StyleMap as StatusBallStyleMap} from "./statusball"

interface NodeTableProps {
    dataSource : NodeDataSource
}

interface NodeTableState {
    nodes : ExtendedNodeRecord[]
    activeRow : string
}

interface NodeLastSeenProps {
    lastseen : number
    source? : string
}

interface NodeLastSeenState {
    totalTime : number
}

const nodeHealthMapping = new Map<NodeHealth, keyof(typeof StatusBallStyleMap)>([
    [NodeHealth.GOOD, "good"],
    [NodeHealth.MEH, "meh"],
    [NodeHealth.BAD, "bad"],
    [NodeHealth.UNKNOWN, "unknown"],
]);

class NodeLastSeen extends React.Component<NodeLastSeenProps, NodeLastSeenState> {
    timer : number;

    constructor(props : NodeLastSeenProps) {
        super(props);
        this.state = {totalTime : Math.floor(Date.now()/1000) - this.props.lastseen};
    }

    // converts a number in seconds to a human readable duration
    private sinceTime(t : number) {
        let out = ""
        let hours = Math.floor(t/3600)
        let minutes = Math.floor((t%3600)/60)
        let seconds = t%60

        if(hours > 0) {
            if(out.length > 0) out += " ";
            out += hours + " hour";
            if(hours > 1) out += "s";
        }
        if(minutes > 0) {
            if(out.length > 0) out += " ";
            out += minutes + " minute";
            if(minutes > 1) out += "s";
        }
        if(seconds > 0) {
            if(out.length > 0) out += " ";
            out += seconds + " second";
            if(seconds > 1) out += "s";
        }
        return out
    }

    render() {
        let extratext = this.props.source ? " ("+this.props.source +")" : "";
        if(this.props.lastseen == -1) {
            return "Never" + extratext;
        } else if(this.state.totalTime < 7200) {
            return this.sinceTime(this.state.totalTime) + " ago" + extratext;
        } else {
            let timestamp = new Date(this.props.lastseen*1000);
            const today = new Date();
            if(timestamp.getDate()==today.getDate() &&
                timestamp.getMonth()==today.getMonth() &&
                timestamp.getFullYear() == today.getFullYear()) {
                return "Today " + dateformat(timestamp, "HH:MM:ss o") + extratext;
            } else {
                return dateformat(timestamp, "ddd dd mmm yyyy HH:MM:ss o") + extratext;
            }

            // if we have an update timer, cancel it since we won't be making any more changes
            if(this.timer) {
                window.clearInterval(this.timer);
                this.timer = null;
            }
        }
    }

    componentDidMount() {
        // no point setting up an update timer if it won't ever do anything
        if(this.props.lastseen != -1 && this.state.totalTime < 7200) {
            this.timer = window.setInterval(() => {
                this.setState({totalTime : Math.floor(Date.now()/1000) - this.props.lastseen});
            }, 1000);
        }
    }

    componentWillUnmount() {
        if(this.timer) {
            window.clearInterval(this.timer);
            this.timer = null;
        }
    }

    componentDidUpdate(prevProps: Readonly<NodeLastSeenProps>) {
        if(prevProps.lastseen != this.props.lastseen) {
            // reset counter
            this.setState({totalTime : Math.floor(Date.now()/1000) - this.props.lastseen});
        }
    }
}

export default class NodeTable extends React.Component<NodeTableProps, NodeTableState> {

    private unsubscribers : (() => void)[] = []

    constructor(props : NodeTableProps) {
        super(props);
        this.state = { nodes: [], activeRow: "" }
    }

    render() {
        let nodes = this.state.nodes.sort((a: NodeRecord, b: NodeRecord) => {
            // put ones without ID at the end
            if(isUndefined(a.id) && !isUndefined(b.id)) {
                return -1;
            }
            if(!isUndefined(a.id) && isUndefined(b.id)) {
                return 1;
            }

            if(a.id == b.id) return 0;

            return a.id > b.id ? 1 : -1;
        });
        let rows = nodes.map((node) => {
            let rowStyle = styles.InactiveRow;
            if(node.mqttName == this.state.activeRow) {
                rowStyle = styles.ActiveRow;
            }
            let setActive = () => {
                this.props.dataSource.setActiveRow(node.mqttName);
            }

            let lastSeenNodes : ReactElement[] = []

            if(node.LastSeenMQTT >= 0 || node.LastSeenAPI >= 0) {
                lastSeenNodes.push(<NodeLastSeen source="MQTT" lastseen={node.LastSeenMQTT}/>);
                lastSeenNodes.push(<br/>);
                lastSeenNodes.push(<NodeLastSeen  source="API" lastseen={node.LastSeenAPI}/>);
            } else {
                lastSeenNodes.push(<NodeLastSeen lastseen={node.LastSeen}/>);
            }

            return <tr key={node.mqttName} className={rowStyle}>
                <td>{node.id || "-"}</td>
                <td><a href="#" onClick={setActive}>{node.name}</a></td>
                <td>{node.SettingsVersion || ""}</td>
                <td>{ lastSeenNodes }</td>
                <td><NodeLastSeen lastseen={node.LastStarted}/></td>
                <td><StatusBall state={nodeHealthMapping.get(node.objectHealth)}/></td>
            </tr>;
        })
        return <table className={styles.NodeTable}>
            <thead><tr><th>Id</th><th>Name</th><th>Settings Version</th><th>Last Seen</th><th>Last Started</th><th>Health</th></tr></thead>
            <tbody>{rows}</tbody>
        </table>
    }

    private refresh() {
        let nodes = this.props.dataSource.nodes.map((n) => this.props.dataSource.getNode(n));
        this.setState({nodes: nodes, activeRow: this.props.dataSource.activeRow });
    }

    componentDidMount() {
        this.unsubscribers.push(this.props.dataSource.onDataChange.subscribe(
            this.refresh.bind(this)
        ));
        this.unsubscribers.push(this.props.dataSource.onActiveRowChange.subscribe((row) => {
            this.setState({nodes: this.state.nodes, activeRow: row});
        }));
        this.refresh()
    }

    componentWillUnmount() {
        for(let us of this.unsubscribers) {
            us();
        }
        this.unsubscribers = []
    }
}