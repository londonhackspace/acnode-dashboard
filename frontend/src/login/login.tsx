import React from "react";
import Spinner from "../components/spinner";

import style from "./login.module.css"
import DashAPI from "../apiclient/dashapi";


interface LoginState {
    username : string
    password : string
    working : boolean
}

interface LoginProps {
    api : DashAPI
}

export default class Login extends React.Component<LoginProps, LoginState> {

    constructor(props : LoginProps) {
        super(props);
        this.state = { username: "", password: "", working: false};
    }

    onUsernameChange(event : React.ChangeEvent<HTMLInputElement>) {
        this.setState((state, props) => ({
            username: event.target.value,
            password: state.password,
            working: state.working,
        }));
    }

    onPasswordChange(event : React.ChangeEvent<HTMLInputElement>) {
        this.setState((state, props) => ({
            username: state.username,
            password: event.target.value,
            working: state.working,
        }));
    }

    handleSubmit(event : React.FormEvent<HTMLFormElement>) {
        event.preventDefault();
        this.props.api.login(this.state.username, this.state.password).then((result) => {
            if(!result) {
                this.setState((state, props) => ({
                    username: state.username,
                    password: state.password,
                    working: false,
                }));
            }
            // we actually don't need to handle the success case, because it
            // will cause this element to be removed due to the api object's state
        })
        this.setState((state, props) => ({
            username: state.username,
            password: state.password,
            working: true,
        }));
    }

    render() {
        let header = <h1>ACNode Dashboard Login</h1>;
        if(this.state.working) {
            return <div className={style.loginframe}>
                {header}
                <Spinner></Spinner>
            </div>;
        } else {
            return <div className={style.loginframe}>
                {header}
                <form onSubmit={this.handleSubmit.bind(this)}>
                    Username: <input name="username" value={this.state.username} onChange={this.onUsernameChange.bind(this)}/><br/>
                    Password: <input type="password" value={this.state.password} onChange={this.onPasswordChange.bind(this)}/><br/>
                    <button type="submit">Go Go Go!</button>
                </form>
            </div>;
        }
    }
}