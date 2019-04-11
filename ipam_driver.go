package main
/*
Логика работы:
	GetCapabilities
	GetDefaultAddressSpaces
	GetCapabilities
	GetDefaultAddressSpaces (почему 2 раза, не знаю... вероятно как припев в песне)

	RequestPool
	RequestAddress + gateway
	RequestAddress + mac (для каждого контейнера)
	... работа ...
	ReleaseAddress
	ReleaseAddress + gateway
	ReleasePool
 */

import (
	"github.com/docker/go-plugins-helpers/ipam"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/netlabel"
	"net"
	"errors"
)

const (
	GlobalDefaultAddressSpace = "global"
	LocalDefaultAddressSpace = "local"
)

type IPAMDriver struct {
}

/*
CapabilitiesResponse возвращает, требуется ли для этого IPAM предварительно созданный MAC

RequiresMACAddress:
	Это логическое значение, которое сообщает libnetwork, должен ли драйвер ipam знать MAC-адрес интерфейса для
правильной обработки вызова RequestAddress (). Если true, по запросу CreateEndpoint () libnetwork сгенерирует случайный
MAC-адрес для конечной точки (если явный MAC-адрес еще не был предоставлен пользователем) и передаст его
RequestAddress () при запросе IP-адреса внутри карты параметров. Ключом будет константа
netlabel.MacAddress: "com.docker.network.endpoint.macaddress".
 */
func (d *IPAMDriver) GetCapabilities() (*ipam.CapabilitiesResponse, error) {
	return &ipam.CapabilitiesResponse{RequiresMACAddress:true}, nil
}

/*
GetDefaultAddressSpaces возвращает имена локального и глобального адресного пространства по умолчанию для этого IPAM.
Адресное пространство - это набор непересекающихся пулов адресов, изолированных от пулов других адресных пространств.
Другими словами, один и тот же пул может существовать в N разных адресных пространствах.
Адресное пространство естественно отображается на имя арендатора. В libnetwork значение, связанное с локальным или
глобальным адресным пространством, заключается в том, что локальному адресному пространству не нужно синхронизироваться
по всему кластеру, в то время как глобальные адресные пространства делают это. Если в конфигурации IPAM не указано иное,
libnetwork будет запрашивать пулы адресов из локального или глобального адресного пространства по умолчанию в
зависимости от области создаваемой сети. Например, если в конфигурации не указано иное, libnetwork будет запрашивать пул
адресов из локального адресного пространства по умолчанию для мостовой сети, а из глобального адресного пространства по
умолчанию для оверлейной сети.
 */
func (d *IPAMDriver) GetDefaultAddressSpaces() (*ipam.AddressSpacesResponse, error) {
	return &ipam.AddressSpacesResponse{
		GlobalDefaultAddressSpace: GlobalDefaultAddressSpace,
		LocalDefaultAddressSpace: LocalDefaultAddressSpace,
	}, nil
}

/*
Этот API предназначен для регистрации пула адресов с драйвером IPAM. Несколько идентичных вызовов должны возвращать один
и тот же результат. Драйвер IPAM отвечает за ведение счетчика ссылок для пула.

* AddressSpace - пространство IP-адресов. Обозначает набор непересекающихся пулов.
* Pool - Пул адресов IPv4 или IPv6 в формате CIDR
* SubPool - Необязательное подмножество пула адресов, диапазон ip в формате CIDR
* Options - Карта специфичных для драйвера IPAM параметров
* V6 - Независимо от IPAM выбранный пул должен быть IPv6

AddressSpace является единственным обязательным полем. Если пул не указан, драйвер IPAM может выбрать возврат
самостоятельно выбранного пула адресов. В таком случае флаг V6 должен быть установлен, если вызывающий абонент хочет
пул IPv6, выбранный IPAM. Запрос с пустым Пулом и непустым SubPool должен быть отклонен как недействительный. Если пул
не указан, IPAM выделит один из пулов по умолчанию. Когда Пул не указан, флаг V6 должен быть установлен, если сети
требуется выделение адресов IPv6.

В ответе:
    PoolID является идентификатором для этого пула. Одинаковые пулы должны иметь одинаковый идентификатор пула.
    Pool - это пул в формате CIDR
    Data - это метаданные, предоставленные драйвером IPAM для этого пула.
 */
func (d *IPAMDriver) RequestPool(request *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	res := ipam.RequestPoolResponse{Data:make(map[string]string)}
	if request.V6 {
		// TODO
	} else {
		pool := NewPoolv4()
		// считаем PoolID
		if request.Pool == "" {
			if request.SubPool != "" {
				return nil, errors.New("Invalid request")
			}
			request.Pool = NetToCIDR(NewNetwork())
		}
		// считаем ID пула
		res.PoolID = HashOfString(request.AddressSpace + request.Pool + request.SubPool)
		// присваиваем ID и делаем попытку загрузить
		err := pool.Load(res.PoolID)
		if err == nil {
			//пул существует!!!!!
			// чего с этим делать - непонятно
		}
		// заполняем его диапазонами
		pool.ParseCIDR(request.Pool, request.SubPool)
		// копируем опции
		for k, v := range request.Options {
			pool.Data[k] = v
		}
		res.Pool = NetToCIDR(pool.Network)
		// сохраняем в БД
		pool.Save()
	}
	return &res, nil
}

/*
Этот API предназначен для освобождения ранее зарегистрированного пула адресов.
 */
func (d *IPAMDriver) ReleasePool(request *ipam.ReleasePoolRequest) error {
	pool := NewPoolv4()
	pool.Load(request.PoolID)
	return pool.Delete()
}

/*
Этот API предназначен для резервирования IP-адреса.
Параметры запроса:

    PoolID - это идентификатор пула
    Address - это требуемый адрес в обычной форме IP (A.B.C.D). Если этот адрес не может быть удовлетворен, запрос не
		выполняется. Если пусто, драйвер IPAM выбирает любой доступный адрес в пуле
    Options - это параметры драйвера IPAM


В ответе:
    Address - это выделенный адрес в формате CIDR (A.B.C.D / MM)
    Data - это определенные метаданные драйвера IPAM
 */
func (d *IPAMDriver) RequestAddress(request *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {
	res := ipam.RequestAddressResponse{Data:make(map[string]string)}
	pool := NewPoolv4()
	pool.Load(request.PoolID)
	if v, ok := request.Options[ipamapi.RequestAddressType]; ok {
		if v == netlabel.Gateway {
			// шлюз по умолчанию.
			//назначаем его в пул
			pool.Gateway = net.ParseIP(request.Address)
			// формируем CIDR с маской
			res.Address = IpToCIDR(pool.Gateway, pool.Network.Mask)
			// сохраняем изменения
			pool.Save()
			return &res, nil
		}
	}
	var ip net.IP
	// остальные адреса
	if request.Address == "" {
		// если адрес не задан
		ip, _ = pool.GetFirstFree()
	} else {
		// если у нас есть указание адреса уже
		ip = net.ParseIP(request.Address)
		// проверяем его на свободность
		if !pool.IpIsFree(ip) {
			// ежели занят
			return nil, errors.New("Ip already assigned")
		}
	}
	res.Address = IpToCIDR(ip, pool.Network.Mask)
	// выдаем и регистрируем
	pool.RegisterIP(ip)
	return &res, nil
}

/*
Этот API предназначен для освобождения IP-адреса.
 */
func (d *IPAMDriver) ReleaseAddress(request *ipam.ReleaseAddressRequest) error {
	pool := NewPoolv4()
	pool.Load(request.PoolID)
	ip := net.ParseIP(request.Address)
	return pool.DeregisterIP(ip)
}

/*
go get github.com/docker/go-plugins-helpers
go get go.etcd.io/bbolt/
go get github.com/docker/libnetwork
go get github.com/coreos/go-systemd
go get github.com/docker/go-connections
 */

//      "docker.ipamdriver/1.0",