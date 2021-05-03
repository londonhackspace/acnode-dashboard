import React from "react";

import styles from "./mainframe.module.css";
import Api, {User} from "../apiclient/dashapi"

import Login from "../login/login"

import MainWidget from "./mainwidget"

//injected by webpack
declare var gitHash : string;

export interface MainFrameProps {
    api : Api;
}

interface MainFrameState {
    loginRequired : boolean
    userName : string
    userIsAdmin: boolean
}

export class MainFrame extends React.Component<MainFrameProps,MainFrameState> {
    private unsubscriber : ()=>void = null;

    constructor(props : MainFrameProps) {
        super(props);
        this.state = {loginRequired: false, userName: "", userIsAdmin: false };
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
            };
        });
    }

    onLogout() {
        this.props.api.logout();
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
        return <div className={styles.container}>
                <div className={styles.header} >
                    <div className={styles.title}>ACNode Dashboard</div>
                    <div className={styles.versionInfo}>version {gitHash}</div>
                    <div className={styles.authButtons}><a href="#" onClick={this.onLogout.bind(this)}>Logout</a></div>
                </div>
                <div className={styles.pageBody}>
                    <MainWidget api={this.props.api}></MainWidget>
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