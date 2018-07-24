import React, { Component } from 'react';
import {Button} from '@material-ui/core';

class OAuth2LoginButton extends Component {

    constructor(props) {
        super(props);
        this.authorizeUrl = this.props.authorizeUrl;
        this.clientId = this.props.clientId;
        this.providerName = this.props.providerName;
        this.baseUrl = process.env.REACT_APP_BASE_URL;
    }

    render() {
        const url= `${this.authorizeUrl}?response_type=token&client_id=${this.clientId}&redirect_uri=${this.baseUrl}/auth-callback&scope=openid profile email`;
        return (
            <Button href={url}
                    variant="raised" color="primary"
                    title={`Login to Kyma Surfer with ${this.providerName}`}>Login with {this.providerName}</Button>
        );
    }
}

export default OAuth2LoginButton