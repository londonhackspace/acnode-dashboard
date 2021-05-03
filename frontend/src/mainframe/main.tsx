import React from "react";

import styles from "./mainframe.module.css";
import Api, {User} from "../apiclient/dashapi"

import Login from "../login/login"

import MainWidget from "./mainwidget"
import TabBar from "../components/tabbar";

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

export class MainFrame extends React.Component<MainFrameProps,MainFrameState> {
    private unsubscriber : ()=>void = null;

    constructor(props : MainFrameProps) {
        super(props);
        this.state = {loginRequired: false, userName: "", userIsAdmin: false, activeTab: 0 };
    }

    private updateUserDetails() {
        this.props.api.getUser().then((user : User) => {
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
            return {
              loginRequired: prev.loginRequired,
              userName: prev.userName,
              userIsAdmin: prev.userIsAdmin,
              activeTab: tabId,
            };
        });
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

        let activeObject : React.ReactElement = null;

        switch(this.state.activeTab) {
            case 0:
                activeObject = <MainWidget api={this.props.api}/>;
                break;
            default:
                activeObject = <MainWidget api={this.props.api}/>;
        }

        return <div className={styles.container}>
                <div className={styles.header} >
                    <div className={styles.title}>ACNode Dashboard</div>
                    <div className={styles.versionInfo}>version {gitHash}</div>
                    <div className={styles.authButtons}><a href="#" onClick={this.onLogout.bind(this)}>Logout</a></div>
                </div>
                <div className={styles.pageBody}>
                    <TabBar
                        widgets={["main"]}
                        active={this.state.activeTab}
                        onTabChange={this.tabChangeHandler.bind(this)}
                    />
                    {activeObject}
                </div>
            </div>;
    }

    componentDidMount() {
        this.unsubscriber = this.props.api.onLoginRequired.subscribe(this.onLoginRequiredChanged.bind(this));
        this.updateUserDetails()
    }

    componentWillUnmount() {
        this.unsubscriber();
        this.unsubscriber = null;
    }
}