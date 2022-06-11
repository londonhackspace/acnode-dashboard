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
            let u = new URL(document.URL);
            let node = null;
            if(u.pathname.indexOf('/', 1) > -1) {
                let part = u.pathname.substring(u.pathname.indexOf('/', 1)+1);
                if(nodes.indexOf(part) != -1) {
                    this.getData(part, 1);
                }
            }

           this.setState({ nodeList : nodes, data: null, page: 1, node: node});
        });
    }

    getData(node : string, page : number) {
        this.props.api.getAccessControlForNode(node, page).then((data : AccessControl) => {
            this.setState({ nodeList: this.state.nodeList, data : data, page: page, node: node});
        });
    }

    render() {
        let lst = this.state.nodeList.map((node : string ) : ReactElement => {
            let onclick = (evt : React.MouseEvent<HTMLElement>) => {
                let u = new URL(document.URL);
                let base = u.pathname.substring(0, u.pathname.indexOf("/", 1))
                if( u.pathname.indexOf("/", 1) == -1) {
                    base = u.pathname;
                }
                history.pushState({}, "", base + "/" + node);
                this.getData(node, 1);
                evt.preventDefault();
            };

            let entry = <a key={node} href="#" onClick={onclick} className={styles.logId}>{node}</a>;

            if(this.state.node == node) {
                return <span className={styles.activeList}>{entry}</span>
            }
            return entry
        });

        let logs = null;
        let pagenumbers = null;
        if(this.state.data != null) {
            let entries = this.state.data.entries.map((entry : AccessControlEntry) => {
                console.log("Entry: " + entry.user_name);
                let d = new Date(entry.timestamp*1000);
                let nameObject : ReactElement;

                if(entry.user_id  && entry.user_id.length > 0) {
                    let shortenedid = entry.user_id.substring(2);
                    nameObject = <a href={"https://london.hackspace.org.uk/members/member.php?id="+shortenedid}>{entry.user_name}</a>
                } else {
                    nameObject = <span>{entry.user_name}</span>
                }
                return <tr key={entry.timestamp+entry.user_id}><td>{d.toDateString() + " " + d.toTimeString()}</td>
                        <td>{nameObject}</td><td>{entry.user_card}</td>
                        <td>{entry.success ? "Granted" : "Denied"}</td>
                    </tr>;
            });
            logs = <table className={tablestyles.NodeTable} key="usagetable">
                <thead key="usagetable-head"><tr className={tablestyles.NodeTable}><td>When</td><td>Who</td><td>Card Id</td><td>Granted</td></tr></thead>
                <tbody key="usagetable-body">{entries}</tbody>
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