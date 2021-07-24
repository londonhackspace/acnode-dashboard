import React, {ReactElement} from "react";
import Api, {NodeRecord} from "../apiclient/dashapi";
import DataSource from "../NodeDataSource";
import {MainFrameProps} from "./main";
import Chart from "../components/chart";
import styles from "./mainwidget.module.css";
import NodeTable from "../components/nodetable";
import NodeDetailPanel from "../components/nodedetailpanel"
import Spinner from "../components/spinner"
import ExtendedNodeRecord from "../extendednoderecord";
import NodeEditPanel from "../components/nodeeditpanel";

interface MainWidgetProps {
    ds : DataSource
    isAdmin : boolean
}

interface MainWidgetState {
    activeRow : string
    isEdit : boolean
}

export default class MainWidget extends React.Component<MainWidgetProps, MainWidgetState> {
    private unsubscribers : {():void}[] = [];

    constructor(props : MainWidgetProps) {
        super(props);
        this.state = {activeRow: "", isEdit: false};
        this.createAndSetupDataSource()
    }

    private createAndSetupDataSource() {
        // force a redraw when things change
        this.unsubscribers.push(this.props.ds.onDataChange.subscribe(() => {
            this.forceUpdate();
        }));
        this.unsubscribers.push(this.props.ds.onActiveRowChange.subscribe((current ) => {
            this.setState({activeRow: current.node, isEdit: current.edit});
        }));
    }

    render() {
        let activeNode : ExtendedNodeRecord = null;

        if(this.state.activeRow && this.state.activeRow != "") {
            activeNode = this.props.ds.getNode(this.state.activeRow);
        }

        let mainNode : ReactElement;
        if(this.props.ds.nodes.length > 0)
        {
            mainNode = <NodeTable dataSource={this.props.ds} userIsAdmin={this.props.isAdmin}></NodeTable>;
        } else {
            mainNode = <Spinner></Spinner>;
        }

        let key : string = activeNode == null ? null : activeNode.mqttName;
        let sidePanel = this.state.isEdit ?
            <NodeEditPanel key={key} node={activeNode} ds={this.props.ds}></NodeEditPanel> :
            <NodeDetailPanel key={key} node={activeNode}></NodeDetailPanel>;

        return <div className={styles.mainwidget}>
            <div className={styles.mainRow}>{mainNode}</div>
            <div className={styles.mainRow}>{sidePanel}</div>
        </div>
    }

    componentDidMount() {
        this.createAndSetupDataSource()
    }

    componentWillUnmount() {
        for(let f of this.unsubscribers) {
            f();
        }
        this.unsubscribers = []
    }

}
