import React, {ReactElement} from "react";
import Api from "../apiclient/dashapi";
import DataSource from "../NodeDataSource";
import {MainFrameProps} from "./main";
import Chart from "../components/chart";
import styles from "./mainwidget.module.css";
import NodeTable from "../components/nodetable";

const chartColors : string[] = [
    'rgb(255, 99, 132)',
    'rgb(54, 162, 235)'
];

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
        let nodePanel : ReactElement = null;

        if(this.state.activeRow && this.state.activeRow != "") {
            let node = this.ds.getNode(this.state.activeRow);
            let data = new Map<string,number>();
            data.set("Used Memory", node.MemUsed);
            data.set("Free Memory", node.MemFree);
            nodePanel = <div>{node.name}<br/><Chart type="doughnut" data={data} colors={chartColors}></Chart></div>;
        }

        return <div className={styles.mainwidget}>
            <div className={styles.mainRow}>Nodes:<br/><NodeTable dataSource={this.ds}></NodeTable></div>
            <div className={styles.mainRow}>{nodePanel}</div>
            <div className={styles.wideRow}>LogsHere?</div>
        </div>
    }

    componentDidMount() {
        this.ds.start();
    }

    componentWillUnmount() {
        this.ds.stop();
    }

}
