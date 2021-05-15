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

interface MainWidgetProps {
    ds : DataSource
}

interface MainWidgetState {
    activeRow : string
}

export default class MainWidget extends React.Component<MainWidgetProps, MainWidgetState> {
    private unsubscribers : {():void}[] = [];

    constructor(props : MainWidgetProps) {
        super(props);
        this.state = {activeRow: ""};
        this.createAndSetupDataSource()
    }

    private createAndSetupDataSource() {
        // force a redraw when things change
        this.unsubscribers.push(this.props.ds.onDataChange.subscribe(() => {
            this.forceUpdate();
        }));
        this.unsubscribers.push(this.props.ds.onActiveRowChange.subscribe((current : string) => {
            this.setState({activeRow: current});
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
            mainNode = <NodeTable dataSource={this.props.ds}></NodeTable>;
        } else {
            mainNode = <Spinner></Spinner>;
        }


        return <div className={styles.mainwidget}>
            <div className={styles.mainRow}>{mainNode}</div>
            <div className={styles.mainRow}><NodeDetailPanel node={activeNode}></NodeDetailPanel></div>
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
