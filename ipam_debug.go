package main

import (
	"github.com/docker/go-plugins-helpers/ipam"
)

type IPAMDebug struct {
	Debugger
	Wrap ipam.Ipam
}

func (self *IPAMDebug) GetCapabilities() (*ipam.CapabilitiesResponse, error) {
	self.Println("GetCapabilities call")
	res, err := self.Wrap.GetCapabilities()
	self.Println("GetCapabilities result", res, err)
	return res, err
}

func (self *IPAMDebug) GetDefaultAddressSpaces() (*ipam.AddressSpacesResponse, error) {
	self.Println("GetDefaultAddressSpaces call")
	res, err := self.Wrap.GetDefaultAddressSpaces()
	self.Println("GetDefaultAddressSpaces result", res, err)
	return res, err
}

func (self *IPAMDebug) RequestPool(r *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	self.Println("RequestPool call", r)
	res, err := self.Wrap.RequestPool(r)
	self.Println("RequestPool result", res, err)
	return res, err
}

func (self *IPAMDebug) ReleasePool(r *ipam.ReleasePoolRequest) error {
	self.Println("ReleasePool call", r)
	err := self.Wrap.ReleasePool(r)
	self.Println("ReleasePool result", err)
	return err
}

func (self *IPAMDebug) RequestAddress(r *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {
	self.Println("RequestAddress call", r)
	res, err := self.Wrap.RequestAddress(r)
	self.Println("RequestAddress result", res, err)
	return res, err

}

func (self *IPAMDebug) ReleaseAddress(r *ipam.ReleaseAddressRequest) error {
	self.Println("ReleaseAddress call", r)
	err := self.Wrap.ReleaseAddress(r)
	self.Println("ReleaseAddress result", err)
	return err
}

