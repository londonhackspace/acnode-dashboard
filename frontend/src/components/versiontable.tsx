import React from "react";
import ExtendedNodeRecord from "../extendednoderecord";

import styles from "./nodetable.module.css";

interface VersionTableProps {
    nodes : ExtendedNodeRecord[]
}

interface VersionTableState {

}

export default class Versions extends React.Component<VersionTableProps, VersionTableState> {

    constructor(props : VersionTableProps) {
        super(props);
    }

    render() {
        return <table className={styles.NodeTable}>
            <thead>
            <tr>
                <th>Node</th><th>Settings Version</th><th>Firmware Version</th>
                <th>Version Date</th><th>Message</th>
            </tr>
            </thead>
            <tbody>
            {this.props.nodes.map((node : ExtendedNodeRecord) => {

                if(!node.Version) {
                    return null;
                }
                let versionDate = null
                if(node.VersionDate) {
                    versionDate = node.VersionDate.toDateString();
                }

                let versionCleaned = node.Version.split('-')[0]
                let versionLink = <a href={"https://github.com/londonhackspace/acnode-cl/commit/"+versionCleaned}>{node.Version}</a>;
                return <tr className={styles.InactiveRow}>
                    <td>{node.name}</td><td>{node.SettingsVersion}</td><td>{versionLink}</td>
                    <td>{versionDate || "Unknown" }</td>
                    <td>{node.VersionMessage ? node.VersionMessage.split("\n")[0] : ""}</td>
                </tr>
            })}
            </tbody>
        </table>
    }
}