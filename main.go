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

/*
      "CAP_NET_RAW",
      "CAP_NET_BIND_SERVICE",
      "CAP_AUDIT_READ",
      "CAP_AUDIT_WRITE",
      "CAP_DAC_OVERRIDE",
      "CAP_SETFCAP",
      "CAP_SETPCAP",
      "CAP_SETGID",
      "CAP_SETUID",
      "CAP_MKNOD",
      "CAP_CHOWN",
      "CAP_FOWNER",
      "CAP_FSETID",
      "CAP_KILL",
      "CAP_SYS_CHROOT",

      "CAP_NET_BROADCAST",
      "CAP_SYS_MODULE",
      "CAP_SYS_RAWIO",
      "CAP_SYS_PACCT",
      "CAP_SYS_ADMIN",
      "CAP_SYS_NICE",
      "CAP_SYS_RESOURCE",
      "CAP_SYS_TIME",
      "CAP_SYS_TTY_CONFIG",
      "CAP_AUDIT_CONTROL",
      "CAP_MAC_OVERRIDE",
      "CAP_MAC_ADMIN",
      "CAP_NET_ADMIN",
      "CAP_SYSLOG",
      "CAP_DAC_READ_SEARCH",
      "CAP_LINUX_IMMUTABLE",
      "CAP_IPC_LOCK",
      "CAP_IPC_OWNER",
      "CAP_SYS_PTRACE",
      "CAP_SYS_BOOT",
      "CAP_LEASE",
      "CAP_WAKE_ALARM",
      "CAP_BLOCK_SUSPEND"


 */
