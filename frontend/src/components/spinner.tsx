import React from "react";

import  style from "./spinner.module.css";

// Spinner taken from https://loading.io/css/
export default class Spinner extends React.Component<any,any> {
    render() {
        return <div className={style['lds-spinner']}>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
        </div>;
    }
}