
import ReactDOM from "react-dom";
import React from "react";

import {MainFrame} from "./mainframe/main";

import APIClient from "./apiclient/dashapi"

let client = new APIClient("/api/");

ReactDOM.render(
    React.createElement(MainFrame, {api: client}, null),
    document.getElementById('app')
);

// development autoreload fun
if(module.hot) {
    module.hot.accept();
}