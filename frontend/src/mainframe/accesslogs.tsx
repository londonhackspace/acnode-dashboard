import React, {ReactElement} from "react";
import "../apiclient/dashapi"
import DashAPI, {AccessControl, AccessControlEntry} from "../apiclient/dashapi";

interface AccessLogsState {
    nodeList : string[];
    data : AccessControl;
}

interface AccessLogsProps {
    api : DashAPI
}

export default class AccessLogs extends React.Component<AccessLogsProps, AccessLogsState>
{
    constructor(props : AccessLogsProps) {
        super(props);
        this.state = {
            nodeList : [] as string[],
            data : null,
        };
        props.api.getAccessControlNodes().then((nodes : string[]) => {
           this.setState({ nodeList : nodes, data: null});
        });
    }

    render() {
        let lst = this.state.nodeList.map((node : string ) : ReactElement => {
            let onclick = () => {
                this.props.api.getAccessControlForNode(node, 1).then((data : AccessControl) => {
                    this.setState({ nodeList: this.state.nodeList, data : data});
                });
            };
           return <a href="#" onClick={onclick}>{node}</a> ;
        });

        let logs = null;
        if(this.state.data != null) {
            let entries = this.state.data.entries.map((entry : AccessControlEntry) => {
                console.log("Entry: " + entry.user_name);
                let d = new Date(entry.timestamp*1000);
                return <tr><td>{d.toDateString() + " " + d.toTimeString()}</td>
                        <td>{entry.user_name}</td>
                    </tr>;
            });
            logs = <table>
                <thead><tr><td>When</td><td>Who</td></tr></thead>
                <tbody>{entries}</tbody>
            </table>;
        }

        let shell = <div>{lst}{logs}</div>;
        return shell;
    }
}