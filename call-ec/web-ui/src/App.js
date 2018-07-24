import React, {Component} from 'react';
import {BrowserRouter as Router, Route, Switch} from 'react-router-dom';
import Loadable from 'react-loadable';
import './App.css';

const Home = Loadable({
    loader: () => import('./routes/Home'),
    loading: () => <div>Loading...</div>,
});

const AuthCallback = Loadable({
    loader: () => import('./routes/AuthCallback'),
    loading: () => <div>Loading...</div>,
});

class App extends Component {

    render() {

        return (
            <Router>
                <Switch>
                    <Route exact path="/" component={Home}/>
                    <Route path="/auth-callback" component={AuthCallback}/>
                </Switch>
            </Router>
        );
    }
}

export default App;
