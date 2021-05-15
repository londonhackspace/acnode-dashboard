import React from "react";
import DataSource from "../NodeDataSource";
import ExtendedNodeRecord from "../extendednoderecord";

import VersionTable from "../components/versiontable"

interface VersionProps {
    dataSource: DataSource;
}

interface VersionState {
    nodes : ExtendedNodeRecord[]
}

export default class Versions extends React.Component<VersionProps, VersionState> {
    private unsubscriber : ()=>void = null;

    constructor(props : VersionProps) {
        super(props);
        this.state = { nodes: [] };
    }

    render() {
        return <VersionTable nodes={this.state.nodes}/>;
    }

    updateNodes() {
        this.setState({
            nodes: this.props.dataSource.nodes.map((nodeName : string) => {
                return this.props.dataSource.getNode(nodeName);
            })
        });
    }

    componentDidMount() {
        this.unsubscriber = this.props.dataSource.onDataChange.subscribe(() => {
          this.updateNodes();
        });
        this.updateNodes();
    }

    componentWillUnmount() {
        this.unsubscriber();
    }
};