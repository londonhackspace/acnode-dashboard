import React, {ReactElement} from "react";
import ExtendedNodeRecord from "../extendednoderecord";
import Chart from "./chart";
import styles from "./nodedetailspanel.module.css"

const chartColors : string[] = [
    'rgb(255, 99, 132)',
    'rgb(54, 162, 235)'
];

interface NodeDetailPanelProps {
    node : ExtendedNodeRecord
}

interface NodeDetailPanelState {

}

export default class NodeDetailPanel extends React.Component<NodeDetailPanelProps, NodeDetailPanelState> {

    constructor(props : NodeDetailPanelProps) {
        super(props);
        this.state = {};
    }

    render() {
        if(this.props.node) {

            let parts : ReactElement[] = []

            let node = this.props.node;

            let addNodeProps = (name : string, value : string | number) => {
                parts.push(<div className={styles.nodepropsline} key={name}>
                    <span className={styles.nodepropstitle}>{name}:</span>
                    <span className={styles.nodepropsvalue}>{value || "Unknown"}</span>
                </div>);
            }

            addNodeProps("Type", node.nodeType);
            addNodeProps("MQTT Name", node.mqttName);
            addNodeProps("Status", node.Status);
            addNodeProps("Settings Version", node.SettingsVersion);
            addNodeProps("Settings Version (EEPROM)", node.EEPROMSettingsVersion);
            addNodeProps("Version", node.Version);
            if(node.VersionDate) {
                addNodeProps("Version Date", node.VersionDate.toDateString());
            }
            addNodeProps("Reset Cause", node.ResetCause);

            if(node.healthHints.length > 0) {
                let hintCounter = 0;
                parts.push(<div className={styles.nodehealthints}>
                    <span className={styles.nodepropstitle}>Health Hints:</span>
                    <ul>{node.healthHints.map(hh => <li key={hintCounter++}>{hh}</li>)}</ul>
                </div>);
            }

            if(node.MemUsed > 0) {
                let data = new Map<string,number>();
                data.set("Used Memory", node.MemUsed);
                data.set("Free Memory", node.MemFree);
                parts.push(<Chart type="doughnut" data={data} colors={chartColors}/>);
            }

            return <div>
                    <div className={styles.nodetitle}>{this.props.node.name}</div>
                    {parts}
                </div>;
        } else {
            return <div>Please select a node from the table</div>;
        }

    }
}