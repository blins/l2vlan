package main

import (
	"testing"
	"net"
)

func TestHashOfString(t *testing.T) {
	hash := HashOfString("hello world!")
	if hash != HashOfString("hello world!") {
		t.Error()
	}
}

func TestIncIP(t *testing.T) {
	ip := net.ParseIP("192.168.1.0")
	if IncIP(ip).String() != "192.168.1.1" {
		t.Error()
	}
}

func TestIncIP2(t *testing.T) {
	ip := net.ParseIP("192.168.1.5")
	if IncIP(ip).String() != "192.168.1.6" {
		t.Error()
	}
}

func TestIncIP3(t *testing.T) {
	ip := net.ParseIP("192.168.1.255")
	if IncIP(ip).String() != "192.168.2.0" {
		t.Error()
	}
}

func TestIncIP4(t *testing.T) {
	ip := net.ParseIP("192.168.255.255")
	if IncIP(ip).String() != "192.169.0.0" {
		t.Error()
	}
}

func TestIncNet(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.168.1.0/24")
	if NetToCIDR(IncNet(*n)) != "192.168.2.0/24" {
		t.Error()
	}
}

func TestIncNet2(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.168.255.0/24")
	if NetToCIDR(IncNet(*n)) != "192.169.0.0/24" {
		t.Error()
	}
}

func TestIncNet3(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.168.1.0/30")
	if NetToCIDR(IncNet(*n)) != "192.168.1.4/30" {
		t.Error()
	}
}

func TestIncNet4(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.168.1.0/27")
	if NetToCIDR(IncNet(*n)) != "192.168.1.32/27" {
		t.Error()
	}
}

func TestIncNet5(t *testing.T) {
	_, n, _ := net.ParseCIDR("10.50.0.0/24")
	if NetToCIDR(IncNet(*n)) != "10.50.1.0/24" {
		t.Error()
	}
}


