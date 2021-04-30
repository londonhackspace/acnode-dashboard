
import * as ReactDOM from "react-dom";
import React from "react";

import {MainFrame, MainFrameProps} from "./mainframe/main";

import APIClient from "./apiclient/dashapi"
import NodeDataSource from "./NodeDataSource"

let client = new APIClient("http://localhost:8080/api/");

ReactDOM.render(
    React.createElement(MainFrame, {api: client}, null),
    document.getElementById('app')
);

// development autoreload fun
if(module.hot) {
    module.hot.accept();
}