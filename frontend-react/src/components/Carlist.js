import apm from '../rum'
import React, { Component } from 'react';
import { SERVER_URL } from '../constants.js'
import ReactTable from "react-table";
import 'react-table/react-table.css';
import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css'
import AddCar from './AddCar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

class Carlist extends Component {
    constructor(props) {
        super(props);
        this.state = { cars: [], open: false, message: '' };
    }

    addCar(car) {

        console.log("Before Add Car Start Transaction")
        apm.startTransaction("Add Car", "Car", { managed: true });

        fetch(SERVER_URL + 'api/cars',
        {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },

                body: JSON.stringify(car)
            })
        .then(res => this.fetchCars())
        .catch(err => console.error(err))

    }

    renderEditable = (cellInfo) => {
        return (
            <div
                style={{ backgroundColor: "#fafafa" }}
                contentEditable
                suppressContentEditableWarning
                onBlur={e => {
                    const data = [...this.state.cars];
                    data[cellInfo.index][cellInfo.column.id] = e.target.innerHTML;
                    this.setState({ cars: data });
                }}
                dangerouslySetInnerHTML={{
                    __html: this.state.cars[cellInfo.index][cellInfo.column.id]
                }}
            />);
    }

    // Delete car
    onDelClick = (car) => {
        fetch( SERVER_URL + 'api/cars/' + car.id, { method: 'DELETE' })
            .then(res => {
                this.setState({ open: true, message: 'Car deleted' });
                this.fetchCars();
            })
            .catch(err => {
                this.setState({ open: true, message: 'Error when deleting' });
                console.error(err)
            })
    }

    // Update car
    updateCar(car) {
        fetch(SERVER_URL + 'api/cars/' + car.id,
            {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(car)
            })
            .then(res =>
                this.setState({ open: true, message: 'Changes saved' })
            )
            .catch(err =>
                this.setState({ open: true, message: 'Error when saving' })
            )
    }

    componentDidMount() {
        this.fetchCars();
    }

    fetchCars = () => {
        fetch(SERVER_URL + 'api/cars')
            .then((response) => response.json())
            .then((responseData) => {
                this.setState({
                    cars: responseData
                    //cars: responseData._embedded.cars,
                });
            })
            .catch(err => console.error(err));
            console.log("Before End Transaction")
            //apm.getCurrentTransaction().end()
    }

    generateError() {
        const err = new Error("This is an error generated on purpose for testing!");
        throw err;
    }

    confirmDelete = (row) => {
        confirmAlert({
            message: 'Are you sure to delete?',
            buttons: [
                {
                    label: 'Yes',
                    onClick: () => this.onDelClick(row)
                },
                {
                    label: 'No',
                }]
        })
    }

    render() {
        const columns = [
          {
              Header: 'ID',
              accessor: 'id',
              show: false
          },{
            Header: 'Brand',
            accessor: 'brand',
            Cell: this.renderEditable
        }, {
            Header: 'Model',
            accessor: 'model',
            Cell: this.renderEditable
        }, {
            Header: 'Color',
            accessor: 'color',
            Cell: this.renderEditable
        }, {
            Header: 'Year',
            accessor: 'year',
            Cell: this.renderEditable
        }, {
            Header: 'List $',
            accessor: 'price',
            Cell: this.renderEditable
        }, {
            Header: 'Market $',
            accessor: 'marketEstimate',
        },
        {
            id: 'savebutton',
            sortable: false,
            filterable: false,
            width: 100,
            accessor: 'id',
            Cell: ({ value, row }) => (<Button size="small" variant="text"
                color="primary"
                onClick={() => { this.updateCar(row) }}>Save</Button>)
        }, {
            id: 'delbutton',
            sortable: false,
            filterable: false,
            width: 100,
            accessor: 'id',
            Cell: ({ value, row  }) => (<Button size="small" variant="text" color="secondary"
                onClick={() => { this.confirmDelete(row) }}>Delete</Button>)
        }]

        // Carlist.js render() method's return statement
        return (
            <div className="App">
                <Grid container>
                    <Grid item>
                        <AddCar addCar={this.addCar} fetchCars={this.fetchCars} />
                    </Grid>
                    <Grid>
                        <Button variant="contained" color="secondary"
                            style={{ 'margin': '10px' }}
                            onClick={() => { this.generateError() }}>Error</Button>
                    </Grid>
                </Grid>
                <ReactTable data={this.state.cars} columns={columns}
                    filterable={true} pageSize={10} />
                <Snackbar
                    style={{ width: 300, color: 'green' }}
                    open={this.state.open} onClose={this.handleClose}
                    autoHideDuration={1500} message={this.state.message} />
            </div>
        );
    }
    handleClose = (event, reason) => {
        this.setState({ open: false });
    };

}
export default Carlist;
