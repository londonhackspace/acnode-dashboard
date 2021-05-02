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
    api : Api;
}

interface MainWidgetState {
    activeRow : string
}

export default class MainWidget extends React.Component<MainWidgetProps, MainWidgetState> {
    private ds : DataSource;

    constructor(props : MainWidgetProps) {
        super(props);
        this.state = {activeRow: ""};
        this.createAndSetupDataSource()
    }

    private createAndSetupDataSource() {
        this.ds = new DataSource(this.props.api);
        // force a redraw when things change
        this.ds.onDataChange.subscribe(() => {
            this.forceUpdate();
        })
        this.ds.onActiveRowChange.subscribe((current : string) => {
            this.setState({activeRow: current});
        });
    }

    componentDidUpdate(prevProps : Readonly<MainFrameProps>) {
        if(prevProps.api !== this.props.api) {
            this.ds.stop();
            this.createAndSetupDataSource();
            this.ds.start();
        }
    }

    render() {
        let activeNode : ExtendedNodeRecord = null;

        if(this.state.activeRow && this.state.activeRow != "") {
            activeNode = this.ds.getNode(this.state.activeRow);
        }

        let mainNode : ReactElement;
        if(this.ds.nodes.length > 0)
        {
            mainNode = <NodeTable dataSource={this.ds}></NodeTable>;
        } else {
            mainNode = <Spinner></Spinner>;
        }


        return <div className={styles.mainwidget}>
            <div className={styles.mainRow}>Nodes:<br/>{mainNode}</div>
            <div className={styles.mainRow}><NodeDetailPanel node={activeNode}></NodeDetailPanel></div>
        </div>
    }

    componentDidMount() {
        this.ds.start();
    }

    componentWillUnmount() {
        this.ds.stop();
    }

}
