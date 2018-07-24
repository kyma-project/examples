import React, {Component} from 'react';
import {Redirect} from 'react-router';
import qs from "query-string";
import {StoreContext} from '../contexts/StoreContext';
import {logIn} from "../actions";

import AuthService from '../services/AuthService';

class AuthCallback extends Component {

    constructor(props) {
        super(props);
        this.state = {
            user: null,
            err: null
        };
        this.authService = new AuthService();
    }

    async componentDidMount() {

        const tokenParams = qs.parse(this.props.location.hash);

        try {
            const user = await this.authService.logInByIdToken(tokenParams);
            this.setState({user: user, err: null});
        }
        catch (err) {
            this.setState({user: null, err: err.toString()});
        }
    }

    render() {

        return (
            <StoreContext.Consumer>
                {({store}) => {

                    if (this.state.user) {

                        store.dispatch(logIn(this.state.user));
                        return (
                            <Redirect to="/"/>
                        )
                    }
                    if (this.state.err) {
                        return (<p style={{color: "red"}}>{this.state.err}</p>);
                    }
                }}
            </StoreContext.Consumer>
        );
    }
}

export default AuthCallback