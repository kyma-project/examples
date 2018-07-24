import React, {Component} from 'react';
import Button from '@material-ui/core/Button';
import {logOut} from '../actions';
import {StoreContext} from '../contexts/StoreContext';

class LogoutButton extends Component {

    constructor(props) {
        super(props);
        this.store = this.props.store;
    }

    render() {
        return (
            <StoreContext.Consumer>
                {({store}) => {

                    return (
                        <Button onClick={() => store.dispatch(logOut())}
                                variant="raised" color="primary"
                                title={"Logout from Kyma Surfer"}>Logout</Button>
                    )
                }}
            </StoreContext.Consumer>
        );
    }
}

export default LogoutButton