package main

import (
	"github.com/docker/go-plugins-helpers/network"
)

type NetworkDebug struct {
	Wrap network.Driver
	Debugger
}

func (self *NetworkDebug) GetCapabilities() (*network.CapabilitiesResponse, error) {
	self.Println("GetCapabilities call")
	res, err := self.Wrap.GetCapabilities()
	self.Println("GetCapabilities result", res, err)
	return res, err
}

func (self *NetworkDebug) CreateNetwork(r *network.CreateNetworkRequest) error {
	self.Println("CreateNetwork call", r)
	err := self.Wrap.CreateNetwork(r)
	self.Println("CreateNetwork result", err)
	return err
}

func (self *NetworkDebug) AllocateNetwork(r *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	self.Println("AllocateNetwork call", r)
	res, err := self.Wrap.AllocateNetwork(r)
	self.Println("AllocateNetwork result", res, err)
	return res, err
}

func (self *NetworkDebug) DeleteNetwork(r *network.DeleteNetworkRequest) error {
	self.Println("DeleteNetwork call", r)
	err := self.Wrap.DeleteNetwork(r)
	self.Println("DeleteNetwork result", err)
	return err
}

func (self *NetworkDebug) FreeNetwork(r *network.FreeNetworkRequest) error {
	self.Println("FreeNetwork call", r)
	err := self.Wrap.FreeNetwork(r)
	self.Println("FreeNetwork result", err)
	return err
}

func (self *NetworkDebug) CreateEndpoint(r *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	self.Println("CreateEndpoint call", r)
	res, err := self.Wrap.CreateEndpoint(r)
	self.Println("CreateEndpoint result", res, err)
	return res, err
}

func (self *NetworkDebug) DeleteEndpoint(r *network.DeleteEndpointRequest) error {
	self.Println("DeleteEndpoint call", r)
	err := self.Wrap.DeleteEndpoint(r)
	self.Println("DeleteEndpoint result", err)
	return err
}

func (self *NetworkDebug) EndpointInfo(r *network.InfoRequest) (*network.InfoResponse, error) {
	self.Println("EndpointInfo call", r)
	res, err := self.Wrap.EndpointInfo(r)
	self.Println("EndpointInfo result", res, err)
	return res, err
}

func (self *NetworkDebug) Join(r *network.JoinRequest) (*network.JoinResponse, error) {
	self.Println("Join call", r)
	res, err := self.Wrap.Join(r)
	self.Println("Join result", res, err)
	return res, err
}

func (self *NetworkDebug) Leave(r *network.LeaveRequest) error {
	self.Println("Leave call", r)
	err := self.Wrap.Leave(r)
	self.Println("Leave result", err)
	return err
}

func (self *NetworkDebug) DiscoverNew(r *network.DiscoveryNotification) error {
	self.Println("DiscoverNew call", r)
	err := self.Wrap.DiscoverNew(r)
	self.Println("DiscoverNew result", err)
	return err
}

func (self *NetworkDebug) DiscoverDelete(r *network.DiscoveryNotification) error {
	self.Println("DiscoverDelete call", r)
	err := self.Wrap.DiscoverDelete(r)
	self.Println("DiscoverDelete result", err)
	return err
}

func (self *NetworkDebug) ProgramExternalConnectivity(r *network.ProgramExternalConnectivityRequest) error {
	self.Println("ProgramExternalConnectivity call", r)
	err := self.Wrap.ProgramExternalConnectivity(r)
	self.Println("ProgramExternalConnectivity result", err)
	return err
}

func (self *NetworkDebug) RevokeExternalConnectivity(r *network.RevokeExternalConnectivityRequest) error {
	self.Println("RevokeExternalConnectivity call", r)
	err := self.Wrap.RevokeExternalConnectivity(r)
	self.Println("RevokeExternalConnectivity result", err)
	return err
}



