import React, { Component } from 'react';
import './App.css';
import CarList from './components/Carlist.js';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';

const styles = {
  root: {
    flexGrow: 1,
  },
  grow: {
    flexGrow: 1,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
};

class App extends Component {
  
  render() {
    return (
      <div className="App">
        <AppBar position="static">
          <Toolbar>
          <Typography variant="h6" color="inherit">
            Car List
          </Typography>
          </Toolbar>
        </AppBar>
        <CarList />
      </div>
    );
  }
}

App.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(App);