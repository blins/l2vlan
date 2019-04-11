package main

import (
	"log"
	"strconv"
	"net"
	"github.com/milosgajdos83/tenus"
	"strings"
	"errors"
	"go.etcd.io/bbolt"
	"bytes"
	"encoding/gob"
)

func createVlanAndAddToBridge(br tenus.Bridger, parent string, vlanId int) {
	// checking vlan exists
	ifName := strings.Join([]string{parent, strconv.Itoa(vlanId)}, ".")
	vlanif, err := net.InterfaceByName(ifName)
	if err != nil {
		// link to external network by vlan
		vlanif, err := tenus.NewVlanLinkWithOptions(parent, tenus.VlanOptions{Id: uint16(vlanId), Dev: ifName, MacAddr: GenerateMac().String()})
		if err != nil {
			log.Fatalln("error creating vlan:", err)
		}
		vlanif.SetLinkUp()
		br.AddSlaveIfc(vlanif.NetInterface())
	} else {
		log.Println("Vlan interface", ifName, "exists")
		err = br.AddSlaveIfc(vlanif)
		if err != nil {
			log.Println("Error adding interface", ifName, "to bridge")
		}
	}
}

type Network struct {
	ID string
	BridgeName string
	VethName string
	Gateway string
}

func (n *Network) Load(id string) error {
	n.ID = id
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(n.ID))
		if b != nil {
			data := b.Get([]byte("binary"))
			if data != nil {
				reader := bytes.NewReader(data)
				enc := gob.NewDecoder(reader)
				return enc.Decode(n)
			}
		}
		return errors.New("Network not exists")
	})
}

func (n *Network) Save() error {
	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(n.ID))
		if err != nil { return err}
		var writer bytes.Buffer
		dec := gob.NewEncoder(&writer)
		err = dec.Encode(n)
		if err != nil {
			return err
		}
		return b.Put([]byte("binary"), writer.Bytes())
	})
}

func (n *Network) Delete() error {
	return db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(n.ID))
	})
}

func (n *Network) GetOrCreateBridge() (tenus.Bridger, error) {
	br, err := tenus.BridgeFromName(n.BridgeName)
	if err != nil {
		br, err = tenus.NewBridgeWithName(n.BridgeName)
		br.SetLinkUp()
		if err != nil {
			log.Fatalln("error on creating bridge:", err)
			return nil, err
		}
	} else {
		log.Println("Bridge", n.BridgeName, "exists")
	}
	return br, err
}

func (n *Network) HostIfName() string {
	return n.VethName + "h"
}

func (n *Network) ContainerIfName() string {
	return n.VethName + "c"
}
