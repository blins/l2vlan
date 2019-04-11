package main

import (
	"testing"
	"github.com/docker/go-plugins-helpers/ipam"
)

func TestRightIPAMDriver_GetCapabilities(t *testing.T) {
	d := IPAMDriver{}
	r, err := d.GetCapabilities()
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error()
	}
	if r.RequiresMACAddress != true {
		t.Error()
	}
}

func TestRightIPAMDriver_GetDefaultAddressSpaces(t *testing.T) {
	d := IPAMDriver{}
	r, err := d.GetDefaultAddressSpaces()
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error()
	}
	if r.LocalDefaultAddressSpace != LocalDefaultAddressSpace || r.GlobalDefaultAddressSpace != GlobalDefaultAddressSpace {
		t.Error()
	}
}

func TestRightIPAMDriver_RequestPool(t *testing.T) {
	StartDatabase(TestDatabase)
	defer ShutdownDatabase()
	d := IPAMDriver{}
	request := ipam.RequestPoolRequest{}
	request.AddressSpace = LocalDefaultAddressSpace
	request.V6 = false

	res, err := d.RequestPool(&request)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error()
	}
	if res.PoolID == "" {
		t.Error()
	}
	if res.Pool == "" {
		t.Error()
	}
	if res.Pool == "0.0.0.0/0" {
		t.Error()
	}

	rp_request := ipam.ReleasePoolRequest{PoolID:res.PoolID}
	err = d.ReleasePool(&rp_request)
	if err != nil {
		t.Error(err)
	}
}

func TestRightIPAMDriver_RequestPool2(t *testing.T) {
	StartDatabase(TestDatabase)
	defer ShutdownDatabase()
	d := IPAMDriver{}
	request := ipam.RequestPoolRequest{}
	request.AddressSpace = LocalDefaultAddressSpace
	request.SubPool = "192.168.1.2/32"
	request.V6 = false

	res, err := d.RequestPool(&request)
	if err == nil {
		t.Error(err)
	}
	if res != nil {
		t.Error()
	}

}

func TestRightIPAMDriver_RequestPool3(t *testing.T) {
	StartDatabase(TestDatabase)
	defer ShutdownDatabase()
	d := IPAMDriver{}
	request := ipam.RequestPoolRequest{}
	request.AddressSpace = LocalDefaultAddressSpace
	request.V6 = false

	request.Pool = "192.168.1.1/24"

	res, err := d.RequestPool(&request)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error()
	}
	if res.PoolID == "" {
		t.Error()
	}
	if res.Pool != "192.168.1.0/24" {
		t.Error()
	}

	rp_request := ipam.ReleasePoolRequest{PoolID:res.PoolID}
	err = d.ReleasePool(&rp_request)
	if err != nil {
		t.Error(err)
	}
}

func TestRightIPAMDriver_RequestPool4(t *testing.T) {
	StartDatabase(TestDatabase)
	defer ShutdownDatabase()
	d := IPAMDriver{}
	request := ipam.RequestPoolRequest{}
	request.AddressSpace = LocalDefaultAddressSpace
	request.V6 = false

	request.Pool = "192.168.1.1/24"
	request.SubPool = "192.168.1.4/30"

	res, err := d.RequestPool(&request)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error()
	}
	if res.PoolID == "" {
		t.Error()
	}
	if res.Pool != "192.168.1.0/24" {
		t.Error()
	}

	rp_request := ipam.ReleasePoolRequest{PoolID:res.PoolID}
	err = d.ReleasePool(&rp_request)
	if err != nil {
		t.Error(err)
	}
}

func TestRightIPAMDriver_RequestAddress(t *testing.T) {
	StartDatabase(TestDatabase)
	defer ShutdownDatabase()
	d := IPAMDriver{}
	request := ipam.RequestPoolRequest{}
	request.AddressSpace = LocalDefaultAddressSpace
	request.V6 = false

	request.Pool = "192.168.1.1/24"
	request.SubPool = "192.168.1.4/30"

	pool, err := d.RequestPool(&request)

	r_addr := ipam.RequestAddressRequest{PoolID:pool.PoolID}
	resp_addr, err := d.RequestAddress(&r_addr)
	if err != nil {
		t.Error(err)
	}
	if resp_addr == nil {
		t.Error()
	}
	if resp_addr.Address != "192.168.1.4/24" {
		t.Error()
	}

	rp_addr := ipam.ReleaseAddressRequest{PoolID:pool.PoolID, Address:resp_addr.Address}
	err = d.ReleaseAddress(&rp_addr)
	if err != nil {
		t.Error(err)
	}

	rp_request := ipam.ReleasePoolRequest{PoolID:pool.PoolID}
	err = d.ReleasePool(&rp_request)
	if err != nil {
		t.Error(err)
	}
}

