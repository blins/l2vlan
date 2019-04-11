package main

import (
	"github.com/docker/go-plugins-helpers/ipam"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/docker/go-plugins-helpers/sdk"
	"net/http"
	"log"
)

const (
	manifest = `{"Implements": ["NetworkDriver", "IpamDriver"]}`

	ipamcapabilitiesPath   = "/IpamDriver.GetCapabilities"
	addressSpacesPath  = "/IpamDriver.GetDefaultAddressSpaces"
	requestPoolPath    = "/IpamDriver.RequestPool"
	releasePoolPath    = "/IpamDriver.ReleasePool"
	requestAddressPath = "/IpamDriver.RequestAddress"
	releaseAddressPath = "/IpamDriver.ReleaseAddress"

	networkcapabilitiesPath    = "/NetworkDriver.GetCapabilities"
	allocateNetworkPath = "/NetworkDriver.AllocateNetwork"
	freeNetworkPath     = "/NetworkDriver.FreeNetwork"
	createNetworkPath   = "/NetworkDriver.CreateNetwork"
	deleteNetworkPath   = "/NetworkDriver.DeleteNetwork"
	createEndpointPath  = "/NetworkDriver.CreateEndpoint"
	endpointInfoPath    = "/NetworkDriver.EndpointOperInfo"
	deleteEndpointPath  = "/NetworkDriver.DeleteEndpoint"
	joinPath            = "/NetworkDriver.Join"
	leavePath           = "/NetworkDriver.Leave"
	discoverNewPath     = "/NetworkDriver.DiscoverNew"
	discoverDeletePath  = "/NetworkDriver.DiscoverDelete"
	programExtConnPath  = "/NetworkDriver.ProgramExternalConnectivity"
	revokeExtConnPath   = "/NetworkDriver.RevokeExternalConnectivity"
)

// ErrorResponse is a formatted error message that libnetwork can understand
type ErrorResponse struct {
	Err string
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}


type NetworkIpamHandler struct {
	ipam ipam.Ipam
	driver network.Driver
	sdk.Handler
}

func NewNetworkIpamHandler(n network.Driver, i ipam.Ipam) *NetworkIpamHandler {
	h := &NetworkIpamHandler{driver: n, ipam: i, Handler: sdk.NewHandler(manifest)}
	h.initIpam()
	h.initDriver()
	return h
}

func (h *NetworkIpamHandler) initIpam() {
	h.HandleFunc(ipamcapabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.ipam.GetCapabilities()
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(addressSpacesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.ipam.GetDefaultAddressSpaces()
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(requestPoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ipam.RequestPoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.ipam.RequestPool(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(releasePoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ipam.ReleasePoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.ipam.ReleasePool(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(requestAddressPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ipam.RequestAddressRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.ipam.RequestAddress(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(releaseAddressPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ipam.ReleaseAddressRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.ipam.ReleaseAddress(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
}


func (h *NetworkIpamHandler) initDriver() {
	h.HandleFunc(networkcapabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.driver.GetCapabilities()
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		if res == nil {
			sdk.EncodeResponse(w, NewErrorResponse("Network driver must implement GetCapabilities"), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(createNetworkPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers createnetwork")
		req := &network.CreateNetworkRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.CreateNetwork(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(allocateNetworkPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers allocatenetwork")
		req := &network.AllocateNetworkRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.AllocateNetwork(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(deleteNetworkPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.DeleteNetworkRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DeleteNetwork(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(freeNetworkPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.FreeNetworkRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.FreeNetwork(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(createEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.CreateEndpointRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.CreateEndpoint(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(deleteEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.DeleteEndpointRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DeleteEndpoint(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(endpointInfoPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.InfoRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.EndpointInfo(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(joinPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.JoinRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.Join(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})
	h.HandleFunc(leavePath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.LeaveRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.Leave(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(discoverNewPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.DiscoveryNotification{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DiscoverNew(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(discoverDeletePath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.DiscoveryNotification{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.DiscoverDelete(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(programExtConnPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.ProgramExternalConnectivityRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.ProgramExternalConnectivity(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
	h.HandleFunc(revokeExtConnPath, func(w http.ResponseWriter, r *http.Request) {
		req := &network.RevokeExternalConnectivityRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.RevokeExternalConnectivity(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})
}
