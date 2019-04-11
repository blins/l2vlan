package main

import (
	"fmt"
	"net"
	"go.etcd.io/bbolt"
	"bytes"
	"encoding/gob"
	"errors"
)

func NewPoolv4() *Poolv4 {
	return &Poolv4{
		Data: make(map[string]string),
	}
}

type Poolv4 struct {
	ID string
	Network net.IPNet
	Subnet net.IPNet
	Gateway net.IP
	Data map[string]string
}

func (pool Poolv4) String() string {
	return fmt.Sprintf("Pool{ID: %s, Network: %s, Subnet: %s, Gateway: %s}", pool.ID, NetToCIDR(pool.Network), NetToCIDR(pool.Subnet), IpToCIDR(pool.Gateway, pool.Network.Mask))
}

// загружает пул по ID из БД
func (pool *Poolv4) Load(id string) error {
	pool.ID = id
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(pool.ID))
		if b != nil {
			data := b.Get([]byte("binary"))
			if data != nil {
				reader := bytes.NewReader(data)
				enc := gob.NewDecoder(reader)
				return enc.Decode(pool)
			}
		}
		return errors.New("Pool not exists")
	})
}

// сохраняет пул
func (pool *Poolv4) Save() error {
	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(pool.ID))
		if err != nil { return err}
		var writer bytes.Buffer
		dec := gob.NewEncoder(&writer)
		err = dec.Encode(pool)
		if err != nil {
			return err
		}
		return b.Put([]byte("binary"), writer.Bytes())
	})
}

// удалить пул из хранилища
func (pool *Poolv4) Delete() error {
	return db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(pool.ID))
	})
}

// разбирает адресацию сетей и подсетей
func (pool *Poolv4) ParseCIDR(network string, subnet string) error {
	_, n, err := net.ParseCIDR(network)
	if err != nil { return err}
	pool.Network = *n
	pool.Network.IP = pool.Network.IP.Mask(pool.Network.Mask)
	if subnet != "" {
		_, n, err := net.ParseCIDR(subnet)
		if err != nil { return err}
		pool.Subnet = *n
	}
	return nil
}

// помечает IP как выданный
func (pool *Poolv4) RegisterIP(ip net.IP) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(pool.ID))
		if b == nil {
			return errors.New("Pool not exists")
		}
		return b.Put([]byte(ip.String()), []byte("1"))
	})
}

// удаляет IP адрес
func (pool *Poolv4) DeregisterIP(ip net.IP) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(pool.ID))
		if b == nil {
			return errors.New("Pool not exists")
		}
		return b.Delete([]byte(ip.String()))
	})
}

//проверяет что адрес свободен
func (pool *Poolv4) IpIsFree(ip net.IP) bool {
	// сеть для проверки
	check := pool.Network
	if !pool.Subnet.IP.IsUnspecified() {
		// если определен пул выдаваемых адресов в большой сети, то начинать с него
		check = pool.Subnet
	}
	if !pool.Network.Contains(IncIP(ip)) || check.Contains(ip) {
		// проверка на broadcast
		// или что мы не вылезли за диапазон
		return false
	}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(pool.ID))
		if b == nil {
			return errors.New("Pool not exists")
		}
		// проверка на существование любого значения с названием адреса. Если значения нет, то адрес свободен
		if b.Get([]byte(ip.String())) != nil {
			return errors.New("Ip exists")
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

// возвращает первый свободный адрес
func (pool *Poolv4) GetFirstFree() (net.IP, error) {
	var res net.IP
	res = IncIP(pool.Network.IP)
	// сеть для проверки
	check := pool.Network
	if !pool.Subnet.IP.IsUnspecified() {
		// если определен пул выдаваемых адресов в большой сети, то начинать с него
		res = pool.Subnet.IP
		check = pool.Subnet
	}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(pool.ID))
		if b == nil {
			return errors.New("Pool not exists")
		}
		for check.Contains(res) {
			// проверка, что IP не шлюз по умолчанию. Он всегда посылается сначала для регистрации видать
			if pool.Gateway.Equal(res) {
				res = IncIP(res)
				continue
			}
			// проверка на существование любого значения с названием адреса. Если значения нет, то адрес свободен
			if b.Get([]byte(res.String())) == nil {
				break
			}
			// добавить еденичку
			res = IncIP(res)
		}
		return nil
	})
	if err != nil {
		return net.ParseIP("0.0.0.0"), errors.New("New ip not available")
	}
	if !check.Contains(res) {
		// или что мы не вылезли за диапазон
		return net.ParseIP("0.0.0.0"), errors.New("New ip not available")
	}
	if !pool.Network.Contains(IncIP(res)) {
		// проверка на broadcast
		return net.ParseIP("0.0.0.0"), errors.New("New ip not available")
	}
	return res, nil
}

// Helpers

