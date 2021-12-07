import React from "react";

import styles from "./mainframe.module.css";
import Api, {User} from "../apiclient/dashapi"

import Login from "../login/login"

import MainWidget from "./mainwidget"
import Versions from "./versions"
import TabBar from "../components/tabbar";
import DataSource from "../NodeDataSource";
import AccessLogs from "./accesslogs";

//injected by webpack
declare var gitHash : string;

export interface MainFrameProps {
    api : Api;
}

interface MainFrameState {
    loginRequired : boolean
    userName : string
    userIsAdmin: boolean
    activeTab : number
}

interface tabRecord {
    getWidget : ()=>React.ReactElement;
    name : string;
    requiresAdmin: boolean;
    shortName : string;
}

export class MainFrame extends React.Component<MainFrameProps,MainFrameState> {
    private unsubscriber : ()=>void = null;
    private ds : DataSource;

    private widgets : tabRecord[] = [
        {
            name: "main",
            getWidget: () => <MainWidget ds={this.ds} isAdmin={this.state.userIsAdmin}/>,
            requiresAdmin: false,
            shortName: "main",
        },
        {
            name: "Versions",
            getWidget: () => <Versions dataSource={this.ds} />,
            requiresAdmin: false,
            shortName: "versions",
        },
        {
            name: "Access Logs",
            getWidget: ()=> <AccessLogs api={this.props.api}/>,
            requiresAdmin: true,
            shortName: "accesslogs",
        },
    ];

    constructor(props : MainFrameProps) {
        super(props);
        this.state = {loginRequired: false, userName: "", userIsAdmin: false, activeTab: 0 };
        this.ds = new DataSource(this.props.api);

        window.onpopstate = this.onHistoryPop.bind(this);
    }

    private onHistoryPop(evt : PopStateEvent) {
        let newIndex = this.getTabIdFromUrl();
        this.setState((prev: MainFrameState): MainFrameState => {
            return {
                loginRequired: prev.loginRequired,
                userName: prev.userName,
                userIsAdmin: prev.userIsAdmin,
                activeTab: newIndex,
            };
        })
    }

    private updateUserDetails() {
        this.props.api.getUser().then((user : User) => {
            // we shouldn't get called if there isn't a user but
            // it seems sometimes we do?
            if(!user) return;
            this.setState((prev : MainFrameState) : MainFrameState => {
                let username = user.name;
                if(username.length == 0) {
                    username = user.username
                }
                return {
                    loginRequired: prev.loginRequired,
                    userName: username,
                    userIsAdmin: user.admin,
                    activeTab: prev.activeTab,
                };
            });
        });
    }

    onLoginRequiredChanged(loginRequired : boolean) {
        if(!loginRequired) this.updateUserDetails();
        console.log(loginRequired ? "Login required" : "Login not required");
        this.setState((prev : MainFrameState) : MainFrameState => {
            return {
                loginRequired: loginRequired,
                userName: loginRequired ? "" : prev.userName,
                userIsAdmin: loginRequired ? false : prev.userIsAdmin,
                activeTab: prev.activeTab,
            };
        });
    }

    onLogout() {
        this.props.api.logout();
    }

    tabChangeHandler(tabId : number) {
        this.setState((prev) : MainFrameState => {
            history.pushState({}, document.title, "/" + this.widgets[tabId].shortName);
            return {
              loginRequired: prev.loginRequired,
              userName: prev.userName,
              userIsAdmin: prev.userIsAdmin,
              activeTab: tabId,
            };
        });
    }

    getTabIdFromUrl() : number {
        let u = new URL(document.URL);
        let part = u.pathname.substr(1);

        let idx = 0;
        for(let w of this.widgets)
        {
            if(part.startsWith(w.shortName))
            {
                return idx;
            }
            idx++;
        }

        return 0;
    }

    render() {
        if(this.state.loginRequired) {
            return <div className={styles.container}>
                <div className={styles.header} >
                    <div className={styles.title}>ACNode Dashboard</div>
                    <div className={styles.versionInfo}>version {gitHash}</div>
                </div>

                <Login api={this.props.api}></Login>
            </div>
        }

        // Each record is a name and a factory function

        let gitHashClean = gitHash.split('-')[0]

        let gitHashLink = <a className={styles.versionLink} href={"https://github.com/londonhackspace/acnode-dashboard/commit/"+gitHashClean}>{gitHash}</a>

        let widgetsList = this.widgets.filter((w) => {
            if(w.requiresAdmin)
            {
                return this.state.userIsAdmin;
            }
            return true;
        }).map(w => w.name);

        return <div className={styles.container}>
                <div className={styles.header} >
                    <div className={styles.title}>ACNode Dashboard</div>
                    <div className={styles.versionInfo}>version {gitHashLink}</div>
                    <div className={styles.authButtons}><a href="#" onClick={this.onLogout.bind(this)}>Logout</a></div>
                </div>
                <div className={styles.pageBody}>
                    <TabBar
                        widgets={widgetsList}
                        active={this.state.activeTab}
                        onTabChange={this.tabChangeHandler.bind(this)}
                    />
                    <div className={styles.selectedPage}>
                        {this.widgets[this.state.activeTab].getWidget()}
                    </div>
                </div>
            </div>;
    }

    componentDidMount() {
        this.unsubscriber = this.props.api.onLoginRequired.subscribe(this.onLoginRequiredChanged.bind(this));
        this.updateUserDetails()
        this.setState((prev : MainFrameState) : MainFrameState => {
           return {
               loginRequired: prev.loginRequired,
               userName: prev.userName,
               userIsAdmin: prev.userIsAdmin,
               activeTab: this.getTabIdFromUrl()
           }
        });
        this.ds.start();
    }

    componentWillUnmount() {
        this.unsubscriber();
        this.unsubscriber = null;
        this.ds.stop();
    }
}