import React, {ReactElement} from "react";

import styles from "./dash.module.css";
import Api from "../apiclient/dashapi"
import DataSource from "../NodeDataSource"
import Chart from "../components/chart";

const chartColors : string[] = [
    'rgb(255, 99, 132)',
    'rgb(54, 162, 235)'
];

export interface MainFrameProps {
    api : Api;
}

interface MainFrameState {
    activeRow : string
}

export class MainFrame extends React.Component<MainFrameProps,MainFrameState> {
    private ds : DataSource;

    constructor(props : MainFrameProps) {
        super(props);
        this.createAndSetupDataSource();
        this.state = {activeRow: ""};
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

        let els : ReactElement[] = []

        for(let n of this.ds.nodes) {
            let linkHandler = () => {
                this.ds.setActiveRow(n);
            };

            let node = this.ds.getNode(n);

            //els.push(<li key={n}><a href="#" onClick={linkHandler}>{n}</a></li>);
            let lastSeen = "Never";
            if(node.LastSeen != -1) {
                lastSeen = node.LastSeen.toString();
            }
            els.push(<tr key={n} onClick={linkHandler}><td>{node.name}</td><td>{node.id}</td><td>{node.SettingsVersion}</td>
            <td>{lastSeen}</td></tr>);
        }

        let nodePanel : ReactElement = null;

        if(this.state.activeRow && this.state.activeRow != "") {
            let node = this.ds.getNode(this.state.activeRow);
            let data = new Map<string,number>();
            data.set("Used Memory", node.MemUsed);
            data.set("Free Memory", node.MemFree);
            nodePanel = <div>{node.name}<br/><Chart type="doughnut" data={data} colors={chartColors}></Chart></div>;
        }

        return <div className={styles.container}>
            <div className={styles.header} >ACNode Dashboard</div>
            <div className={styles.mainArea}>
                <table>
                    <tr><th>Name</th><th>Id</th><th>Settings Version</th><th>Last Seen</th></tr>
                    {els}
                </table>
                <div>{nodePanel}</div>
                <div className={styles.wideRow}>LogsHere?</div>
            </div>
        </div> ;
    }

    componentDidMount() {
        this.ds.start();
    }

    componentWillUnmount() {
        this.ds.stop();
    }
}