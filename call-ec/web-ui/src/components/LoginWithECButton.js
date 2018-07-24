import React, {Component} from 'react';
import OAuth2LoginButton from './OAuth2LoginButton';

class LoginWithECButton extends Component {

    constructor(props) {
        super(props);
        this.authorizeEndpoint = process.env.REACT_APP_OAUTH2_AUTHORIZE_URL;
    }

    render() {
        return (
            <OAuth2LoginButton
                providerName="EC"
                authorizeUrl={this.authorizeEndpoint}
                clientId="kyma" />
        );
    }
}

export default LoginWithECButton