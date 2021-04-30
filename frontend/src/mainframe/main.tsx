import React from "react";

import styles from "./mainframe.module.css";
import Api from "../apiclient/dashapi"

import MainWidget from "./mainwidget"

export interface MainFrameProps {
    api : Api;
}

interface MainFrameState {

}

export class MainFrame extends React.Component<MainFrameProps,MainFrameState> {
    constructor(props : MainFrameProps) {
        super(props);
    }

    render() {
        return <div className={styles.container}>
                <div className={styles.header} >ACNode Dashboard</div>
                <MainWidget api={this.props.api}></MainWidget>
            </div>;
    }


}