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
      <ul>
    { this.state.devices.map(device => <li key={device.MacAddress}>{device.CName} - {device.LastIP} ({device.Vendor})</li>)}
      </ul>
    )
  }
}