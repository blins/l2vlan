package main

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strconv"
	"math/rand"
	"time"
)

/*
SHA256 для string
 */
func HashOfString(str string) string {
	sum := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sum)
}

// inc increments `ip` by one address, returning it in `ret`
// e.g. 10.0.2.100   -> 10.0.2.101
//      10.0.2.255   -> 10.0.3.0
//      10.0.255.255 -> 10.1.0.0
// `ip` is not modified.
// стянул у кого-то из интернета и переобозвал.
func IncIP(ip net.IP) (ret net.IP) {
	ret = make(net.IP, len(ip))
	copy(ret, ip)

	for j := len(ret) - 1; j >= 0; j-- {
		ret[j]++
		if ret[j] > 0 {
			break
		}
	}
	return
}

/*
Возвращает следующую сеть с той же маской.
 */
func IncNet(n net.IPNet) (res net.IPNet) {
	// copying network
	res.IP = make(net.IP, len(n.IP))
	copy(res.IP, n.IP)
	// mask is ref
	res.Mask = n.Mask
	// сколько бит в маске и длинна адреса
	countBits, ipLen := res.Mask.Size()
	// стартовый байт в массиве c которого начинать отсчет сетей
	startByte := countBits / 8
	// пограничный случай
	if countBits % 8 == 0 {
		startByte --
	}
	firstInc := byte(1 << (uint(ipLen - countBits) % 8))
	for j := startByte; j >= 0; j-- {
		res.IP[j] += firstInc
		if res.IP[j] > 0 {
			break
		}
		firstInc = 1
	}
	return
}

/*
Возвращает строковое представление IP адреса с маской
 */
func IpToCIDR(ip net.IP, mask net.IPMask) string {
	if ip == nil {
		return ""
	}
	postfix, _ := mask.Size()
	return ip.String() + "/" + strconv.Itoa(postfix)
}

/*
Возвращает строковое представление IP сети с маской
В случае, если сеть 0.0.0.0/0, то возвращается пустая строка
 */
func NetToCIDR(n net.IPNet) string {
	if n.IP.IsUnspecified() { return "" }
	return IpToCIDR(n.IP, n.Mask)
}

/*
Генерирует мак случайным образом
 */
func GenerateMac() net.HardwareAddr {
	hw := make(net.HardwareAddr, 6)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Read(hw)
	hw[0] = (hw[0] | 2) & 0xfe // Set local bit, ensure unicast address
	return hw
}
