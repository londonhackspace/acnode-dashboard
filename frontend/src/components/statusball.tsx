import React from "react";

import styles from "./statusball.module.css";

export interface StyleMapType {
    good : string
    bad  : string
    meh : string
    unknown : string
}

export const StyleMap : StyleMapType = {
    good: styles.good,
    bad: styles.bad,
    meh: styles.meh,
    unknown: styles.unknown,
}

interface StatusBallProps {
    state : keyof(StyleMapType)
    text? : string
}

export default class StatusBall extends React.Component<StatusBallProps, null> {

    render() {
        if(this.props.text)
        {
            return <div className={StyleMap[this.props.state]}>{this.props.text}</div>;
        }
        else
        {
            return <div className={StyleMap[this.props.state] + " " + styles.notext}></div>;
        }
    }

}