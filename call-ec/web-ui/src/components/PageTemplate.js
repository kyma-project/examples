import React, {Component, Fragment} from 'react';
import LogoutButton from './LogoutButton';
import LoginWithECButton from './LoginWithECButton';
import {StoreContext} from '../contexts/StoreContext';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import {withStyles} from '@material-ui/core/styles';

const styles = {
    root: {
        flexGrow: 1,
    },
    flex: {
        flex: 1,
    },
};

class PageTemplate extends Component {

    render() {

        const {classes} = this.props;

        return (
            <StoreContext.Consumer>
                {({store}) => {

                    return (
                        <Fragment>
                            <AppBar position="static">
                                <Toolbar>
                                    <Typography variant="title" color="inherit" className={classes.flex}>
                                        Kyma Surfer
                                    </Typography>
                                    {store.getState().isUserLoggedIn && (
                                        <Fragment>
                                            <span style={{marginRight: 10}}>{store.getState().user.id}</span>
                                            <LogoutButton color="inherit"/>
                                        </Fragment>
                                    )}
                                </Toolbar>
                            </AppBar>
                            <div className="App-content">
                                {!store.getState().isUserLoggedIn && (<LoginWithECButton/>)}
                                {store.getState().isUserLoggedIn && this.props.children}
                            </div>
                        </Fragment>
                    );
                }}
            </StoreContext.Consumer>
        );


    }
}

export default withStyles(styles)(PageTemplate);

