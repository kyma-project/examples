import React, {Component, Fragment} from 'react';
import PageTemplate from '../components/PageTemplate'
import {StoreContext} from '../contexts/StoreContext';
import Lambda from '../components/Lambda'

class Home extends Component {

    render() {
        return (
            <PageTemplate>
                <StoreContext.Consumer>
                    {({store}) => {

                        return (
                            <Fragment>
                                {store.getState().isUserLoggedIn && (
                                    <Lambda withUser={store.getState().user} />
                                )}
                                {!store.getState().isUserLoggedIn && (
                                    <div>Home</div>
                                )}
                            </Fragment>
                        );
                    }}
                </StoreContext.Consumer>
            </PageTemplate>
        );
    }
}

export default Home