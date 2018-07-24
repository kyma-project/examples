import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
// import registerServiceWorker from './registerServiceWorker';
import {createStore} from 'redux';
import {StoreContext} from './contexts/StoreContext';
import producers from './reducers';
import CssBaseline from '@material-ui/core/CssBaseline';
import {MuiThemeProvider, createMuiTheme} from '@material-ui/core/styles';
import lightBlue from '@material-ui/core/colors/lightBlue';

const theme = createMuiTheme({
    palette: {
        primary: lightBlue,
    },
});

const store = createStore(producers);

function render() {
    ReactDOM.render(
        <StoreContext.Provider value={{store: store}}>
            <MuiThemeProvider theme={theme}>
                <CssBaseline/>
                <App/>
            </MuiThemeProvider>
        </StoreContext.Provider>,
        document.getElementById('root')
    );
}

store.subscribe(render);

render();

// uncomment if you need some optimizations on production (to serve assets from local cache)
// registerServiceWorker();
