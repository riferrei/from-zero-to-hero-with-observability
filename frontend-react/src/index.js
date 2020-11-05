import apm from './rum'
import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';

apm.setInitialPageLoadName("Car List")

ReactDOM.render(<App />, document.getElementById('root'));

serviceWorker.unregister();
