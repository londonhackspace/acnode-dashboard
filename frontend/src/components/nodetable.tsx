import React from "react";
import dateformat from "dateformat";

import styles from "./nodetable.module.css";

import NodeDataSource from "../NodeDataSource";
import {NodeRecord} from "../apiclient/dashapi";
import {isUndefined} from "webpack-merge/dist/utils";

interface NodeTableProps {
    dataSource : NodeDataSource
}

interface NodeTableState {
    nodes : NodeRecord[]
    activeRow : string
}

interface NodeLastSeenProps {
    lastseen : number
}

interface NodeLastSeenState {
    sinceLoad : number
}

class NodeLastSeen extends React.Component<NodeLastSeenProps, NodeLastSeenState> {
    timer : number;

    constructor(props : NodeLastSeenProps) {
        super(props);
        this.state = {sinceLoad: 0};
    }

    render() {
        let totalTime = this.props.lastseen + this.state.sinceLoad;
        if(this.props.lastseen == -1) {
            return "Never";
        } else if(totalTime < 600) {
            return totalTime + " seconds ago";
        } else {
            let timestamp = new Date(Date.now()-(totalTime*1000));
            return dateformat(timestamp, "ddd dd mmm yyyy hh:mm:ss o")
        }
    }

    componentDidMount() {
        this.timer = window.setInterval(() => {
            this.setState({sinceLoad: this.state.sinceLoad + 1});
        }, 1000);
    }

    componentWillUnmount() {
        window.clearInterval(this.timer);
    }

    componentDidUpdate(prevProps: Readonly<NodeLastSeenProps>) {
        if(prevProps.lastseen != this.props.lastseen) {
            // reset counter
            this.setState({sinceLoad: 0});
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
            return <tr key={node.mqttName} className={rowStyle}>
                <td>{node.id || "-"}</td>
                <td><a href="#" onClick={setActive}>{node.name}</a></td>
                <td>{node.SettingsVersion || ""}</td>
                <td><NodeLastSeen lastseen={node.LastSeen}/></td>
            </tr>;
        })
        return <table className={styles.NodeTable}>
            <thead><tr><th>Id</th><th>Name</th><th>Settings Version</th><th>Last Seen</th></tr></thead>
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