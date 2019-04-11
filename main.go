package main

import (
	"net"
	"log"
)

const (
	DefaultPoolNameIPv4 = "10.50.0.0/24"
	DriverName = "l2vlan"
)

var (
	nets *net.IPNet
)


func NewNetwork() net.IPNet {
	if nets == nil {
		_, nets, _ = net.ParseCIDR(DefaultPoolNameIPv4)
	}
	res := IncNet(*nets)
	copy(nets.IP, res.IP)
	return res
}


func main() {
	err := StartDatabase(DriverName + ".db")
	if err != nil {
		log.Panicln("Unable to start", DriverName, ":", err)
	}
	defer ShutdownDatabase()

	/*///
	// debug version
	ipamDriver := &IPAMDebug{Wrap: &IPAMDriver{}}
	networkDriver := &NetworkDebug{Wrap: &NetworkDriver{}}

	ipamDriver.SetLogger(log.New(os.Stderr, "IPAM l2vlan", log.LstdFlags))
	networkDriver.SetLogger(log.New(os.Stderr, "NETWORK l2vlan ", log.LstdFlags))
	/*///
	ipamDriver := &IPAMDriver{}
	networkDriver := &NetworkDriver{}
	//*///

	handler := NewNetworkIpamHandler(networkDriver, ipamDriver)
	handler.ServeUnix(DriverName, 0666)
}


