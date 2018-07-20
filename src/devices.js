import React from 'react';

import axios from 'axios';

export default class Devices extends React.Component {
  state = {
    devices: []
  }

  componentDidMount() {
    axios.get(`http://localhost:8080/api/devices`)
      .then(res => {
        const devices = res.data;
        console.log(devices)
        this.setState({ devices });
      })
  }

  render() {
    return (
      <table className="table table-bordered table-striped table-hover">
        <thead>
          <tr>
            <th>Name</th>
            <th>IP</th>
            <th>MacAddress</th>
            <th>Vendor</th>
            <th>Last Seen</th>
          </tr>
        </thead>
        <tbody>
    {       
      this.state.devices.map(device => {
          var identifier = device.CName;
          if (device.CName === "") {
            identifier = device.SurName;
          }
          if (identifier == "" ) {
            identifier = "-";
          }
          return <tr key={device.MacAddress}>
            <td>{identifier}</td>
            <td>{device.LastIP}</td>
            <td>{device.MacAddress}</td>
            <td>{device.Vendor}</td>
            <td>{device.LastSeen}</td>
              </tr>
        }
      )
    }
      </tbody>
      </table>
    )
  }
}