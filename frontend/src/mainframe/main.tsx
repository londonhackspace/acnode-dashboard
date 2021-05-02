import React from "react";

import styles from "./mainframe.module.css";
import Api from "../apiclient/dashapi"

import Login from "../login/login"

import MainWidget from "./mainwidget"

export interface MainFrameProps {
    api : Api;
}

interface MainFrameState {
    loginRequired : boolean
}

export class MainFrame extends React.Component<MainFrameProps,MainFrameState> {
    private unsubscriber : ()=>void = null;

    constructor(props : MainFrameProps) {
        super(props);
        this.state = {loginRequired: false};
    }

    onLoginRequiredChanged(loginRequired : boolean) {
        this.setState({loginRequired: loginRequired});
    }

    onLogout() {
        this.props.api.logout();
    }

    render() {
        if(this.state.loginRequired) {
            return <div className={styles.container}>
                <div className={styles.header} >ACNode Dashboard</div>
                <Login api={this.props.api}></Login>
            </div>
        }
        return <div className={styles.container}>
                <div className={styles.header} >
                    <div className={styles.title}>ACNode Dashboard</div>
                    <div className={styles.authButtons}><a href="#" onClick={this.onLogout.bind(this)}>Logout</a></div>
                </div>
                <div className={styles.pageBody}>
                    <MainWidget api={this.props.api}></MainWidget>
                </div>
            </div>;
    }

    componentDidMount() {
        this.unsubscriber = this.props.api.onLoginRequired.subscribe(this.onLoginRequiredChanged.bind(this));
    }

    componentWillUnmount() {
        this.unsubscriber();
        this.unsubscriber = null;
    }
}