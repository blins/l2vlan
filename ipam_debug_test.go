package main

import (
	"github.com/docker/go-plugins-helpers/ipam"
	"errors"
	"testing"
)

type IPAMTester struct {
	lastcall     string
	isNil        bool
	returnedNNil bool
}

func (self *IPAMTester) GetCapabilities() (*ipam.CapabilitiesResponse, error) {
	self.lastcall = "GetCapabilities"
	self.isNil = true
	if self.returnedNNil {
		return &ipam.CapabilitiesResponse{}, nil
	}
	return nil, nil
}

func (self *IPAMTester) GetDefaultAddressSpaces() (*ipam.AddressSpacesResponse, error) {
	self.lastcall = "GetDefaultAddressSpaces"
	self.isNil = true
	if self.returnedNNil {
		return &ipam.AddressSpacesResponse{}, nil
	}
	return nil, nil
}

func (self *IPAMTester) RequestPool(r *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	self.lastcall = "RequestPool"
	self.isNil = r == nil
	if self.returnedNNil {
		return &ipam.RequestPoolResponse{}, nil
	}
	return nil, nil
}

func (self *IPAMTester) ReleasePool(r *ipam.ReleasePoolRequest) error {
	self.lastcall = "ReleasePool"
	self.isNil = r == nil
	if self.returnedNNil {
		return errors.New("")
	}
	return nil
}

func (self *IPAMTester) RequestAddress(r *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {
	self.lastcall = "RequestAddress"
	self.isNil = r == nil
	if self.returnedNNil {
		return &ipam.RequestAddressResponse{}, nil
	}
	return nil, nil
}

func (self *IPAMTester) ReleaseAddress(r *ipam.ReleaseAddressRequest) error {
	self.lastcall = "ReleaseAddress"
	self.isNil = r == nil
	if self.returnedNNil {
		return errors.New("")
	}
	return nil
}

func PrepareIpamTest() (*IPAMTester, ipam.Ipam) {
	var res IPAMTester
	var d DebuggerTest
	i := &IPAMDebug{Wrap:&res}
	i.SetLogger(&d)
	return &res, i
}

func TestIPAMDebug_GetCapabilities(t *testing.T) {
	tt, i := PrepareIpamTest()
	r, err := i.GetCapabilities()
	if tt.lastcall != "GetCapabilities" && tt.isNil != true && r != nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_GetCapabilities2(t *testing.T) {
	tt, i := PrepareIpamTest()
	tt.returnedNNil = true
	r, err := i.GetCapabilities()
	if tt.lastcall != "GetCapabilities" && tt.isNil != true && r == nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_GetDefaultAddressSpaces(t *testing.T) {
	tt, i := PrepareIpamTest()
	r, err := i.GetDefaultAddressSpaces()
	if tt.lastcall != "GetDefaultAddressSpaces" && tt.isNil != true && r != nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_GetDefaultAddressSpaces2(t *testing.T) {
	tt, i := PrepareIpamTest()
	tt.returnedNNil = true
	r, err := i.GetDefaultAddressSpaces()
	if tt.lastcall != "GetDefaultAddressSpaces" && tt.isNil != true && r == nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_ReleaseAddress(t *testing.T) {
	tt, i := PrepareIpamTest()
	err := i.ReleaseAddress(nil)
	if tt.lastcall != "ReleaseAddress" && tt.isNil != true && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_ReleaseAddress2(t *testing.T) {
	tt, i := PrepareIpamTest()
	err := i.ReleaseAddress(&ipam.ReleaseAddressRequest{})
	if tt.lastcall != "ReleaseAddress" && tt.isNil != false && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_ReleaseAddress3(t *testing.T) {
	tt, i := PrepareIpamTest()
	tt.returnedNNil = true
	err := i.ReleaseAddress(nil)
	if tt.lastcall != "ReleaseAddress" && tt.isNil != true && err == nil {
		t.Error()
	}
}

func TestIPAMDebug_ReleasePool(t *testing.T) {
	tt, i := PrepareIpamTest()
	err := i.ReleasePool(nil)
	if tt.lastcall != "ReleasePool" && tt.isNil != true && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_ReleasePool2(t *testing.T) {
	tt, i := PrepareIpamTest()
	err := i.ReleasePool(&ipam.ReleasePoolRequest{})
	if tt.lastcall != "ReleasePool" && tt.isNil != false && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_ReleasePool3(t *testing.T) {
	tt, i := PrepareIpamTest()
	tt.returnedNNil = true
	err := i.ReleasePool(nil)
	if tt.lastcall != "ReleasePool" && tt.isNil != true && err == nil {
		t.Error()
	}
}

func TestIPAMDebug_RequestAddress(t *testing.T) {
	tt, i := PrepareIpamTest()
	r, err := i.RequestAddress(nil)
	if tt.lastcall != "RequestAddress" && tt.isNil != true && r != nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_RequestAddress2(t *testing.T) {
	tt, i := PrepareIpamTest()
	r, err := i.RequestAddress(&ipam.RequestAddressRequest{})
	if tt.lastcall != "RequestAddress" && tt.isNil != false && r != nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_RequestAddress3(t *testing.T) {
	tt, i := PrepareIpamTest()
	tt.returnedNNil = true
	r, err := i.RequestAddress(&ipam.RequestAddressRequest{})
	if tt.lastcall != "RequestAddress" && tt.isNil != false && r == nil && err != nil {
		t.Error()
	}
}


func TestIPAMDebug_RequestPool(t *testing.T) {
	tt, i := PrepareIpamTest()
	r, err := i.RequestPool(nil)
	if tt.lastcall != "RequestPool" && tt.isNil != true && r != nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_RequestPool2(t *testing.T) {
	tt, i := PrepareIpamTest()
	r, err := i.RequestPool(&ipam.RequestPoolRequest{})
	if tt.lastcall != "RequestPool" && tt.isNil != false && r != nil && err != nil {
		t.Error()
	}
}

func TestIPAMDebug_RequestPool3(t *testing.T) {
	tt, i := PrepareIpamTest()
	tt.returnedNNil = true
	r, err := i.RequestPool(&ipam.RequestPoolRequest{})
	if tt.lastcall != "RequestPool" && tt.isNil != false && r == nil && err != nil {
		t.Error()
	}
}

