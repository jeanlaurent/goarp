package main

import (
	"sync"
	"time"
)

type Device struct {
	MacAddress string
	LastIP     string
	FirstSeen  time.Time
	LastSeen   time.Time
	CName      string
	Vendor     string
}

type Devices struct {
	mutex sync.RWMutex
	list  map[string]Device
}

func (d *Devices) exist(macAddress string) bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	_, exist := d.list[macAddress]
	return exist
}

func (d *Devices) get(macAddress string) Device {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.list[macAddress]
}

func (d *Devices) put(device Device) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.list[device.MacAddress] = device
}

func (d *Devices) len() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return len(d.list)
}

func (d *Devices) all() []Device {
	values := make([]Device, 0, len(d.list))

	for _, value := range d.list {
		values = append(values, value)
	}
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return values
}

func newDevices() Devices {
	return Devices{sync.RWMutex{}, map[string]Device{}}
}
