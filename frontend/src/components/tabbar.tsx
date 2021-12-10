import React from "react";

import styles from "./tabbar.module.css";

interface TabBarProps {
    widgets : string[]
    active : number
    onTabChange : (newtab : number) => void
}

interface TabBarState {

}

export default class TabBar extends React.Component<TabBarProps, TabBarState> {
    render() {
        let tabs : React.ReactElement[] = []
        for(let i = 0; i < this.props.widgets.length; i++) {
            tabs.push(
                <button key={this.props.widgets[i]} className={i == this.props.active ? styles.activetab : styles.inactivetab}
                    onClick={() => this.props.onTabChange(i)}>
                    {this.props.widgets[i]}
                </button>);
        }
        return <div className={styles.tabbar}>
            {tabs}
        </div>
    }
}