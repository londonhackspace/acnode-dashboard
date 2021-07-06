import React, {ReactElement} from "react";
import "../apiclient/dashapi"
import DashAPI, {AccessControl, AccessControlEntry} from "../apiclient/dashapi";
import styles from "./accesslogs.module.css"
import tablestyles from "../components/nodetable.module.css"

interface AccessLogsState {
    nodeList : string[];
    data : AccessControl;
    page : number;
    node : string;
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
            page: 1,
            node : null,
        };
        props.api.getAccessControlNodes().then((nodes : string[]) => {
           this.setState({ nodeList : nodes, data: null, page: 1, node: null});
        });
    }

    getData(node : string, page : number) {
        this.props.api.getAccessControlForNode(node, page).then((data : AccessControl) => {
            this.setState({ nodeList: this.state.nodeList, data : data, page: page, node: node});
        });
    }

    render() {
        let lst = this.state.nodeList.map((node : string ) : ReactElement => {
            let onclick = () => {
                this.getData(node, 1);
            };
           return <a href="#" onClick={onclick} className={styles.logId}>{node}</a>;
        });

        let logs = null;
        let pagenumbers = null;
        if(this.state.data != null) {
            let entries = this.state.data.entries.map((entry : AccessControlEntry) => {
                console.log("Entry: " + entry.user_name);
                let d = new Date(entry.timestamp*1000);
                return <tr><td>{d.toDateString() + " " + d.toTimeString()}</td>
                        <td>{entry.user_name}</td><td>{entry.user_card}</td>
                    </tr>;
            });
            logs = <table className={tablestyles.NodeTable}>
                <thead><tr className={tablestyles.NodeTable}><td>When</td><td>Who</td><td>Card Id</td></tr></thead>
                <tbody>{entries}</tbody>
            </table>;
            let pagenumberparts : ReactElement[] = [];
            for(let i = 1; i <= this.state.data.pageCount; i++) {
                let onClick = () => {
                  this.getData(this.state.node, i);
                };
                pagenumberparts.push(<a href="#" onClick={onClick}>{i}</a>);
                pagenumberparts.push(<span> </span>);
            }
            pagenumbers = <div>Page: {pagenumberparts}</div>
        }

        let shell = <div className={styles.accessLogs}>{lst}<br/>{logs}<br/>{pagenumbers}</div>;
        return shell;
    }
}