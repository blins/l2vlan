package main
/*
Логика.
	После enable:
		- GetCapabilities
		- DiscoverNew c адресом ноды и указанием что это она и есть

	После docker network create -d rightipam --subnet=192.168.1.0/24 --gateway=192.168.1.1 --ip-range=192.168.1.4/32 -o vlan_id=2000 -o ext_if=eno1 -o bridge_name=vlan2000 net1
		- CreateNetwork {
			deec618526f0cc2c5ddb7b66acb749b6b0ebf056d2ed0df7a3c20e9baaa5f182
			map[
				com.docker.network.enable_ipv6:false
				com.docker.network.generic:map[
					bridge_name:vlan2000
					ext_if:eno1
					vlan_id:2000
				]
			]
			[0xc00009b940] IPAMData IPv4
			[]
			}

	После docker network rm net1
		- DeleteNetwork

	После docker network create -d rightipam:latest --scope swarm --subnet=192.168.1.0/24 --gateway=192.168.1.1 --ip-range=192.168.1.4/32 -o vlan_id=2000 -o ext_if=eno1 -o bridge_name=vlan2000 net1
		- ничего но параметры драйвера не передаются. Они теряются. Надо через конфиг

	После docker stack deploy
		- CreateEndpoint {
			o38vug9jkt0ia3a8dcz3lwsa0
			f85b2a0c98f126ea3bbf4fcea0e6a4f02a3ddba5e66607a73ebeed2130981c3d
			0xc000170d80
			map[
				com.docker.network.endpoint.exposedports:[
					map[
						Proto:6 Port:80
						]
				]
				com.docker.network.portmap:[]
			]
			}

 */

import (
	"github.com/docker/go-plugins-helpers/network"
	"github.com/docker/libnetwork/netlabel"
	"github.com/milosgajdos83/tenus"
	"net"
	"github.com/docker/libcontainer/netlink"
	"errors"
	"log"
)


var (
	FeatureNotAvailableErr = errors.New("Feature not Available")
)

type NetworkDriver struct {
	Wrap network.Driver
}

/*
Значение «Scope» должно быть «local» или «global», что указывает, может ли выделение ресурсов для сети этого
драйвера выполняться только локально для узла или глобально через кластер узлов.
Любое другое значение не выполнит регистрацию драйвера и вернет ошибку вызывающей стороне.

Аналогично, значение «ConnectivityScope» должно быть либо «локальным», либо «глобальным», что указывает, может ли сеть
драйвера обеспечивать подключение только локально к этому узлу или глобально через кластер узлов.
Если значение отсутствует, libnetwork установит для него значение «Scope»

 */
func (self *NetworkDriver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	return &network.CapabilitiesResponse{
		Scope: network.LocalScope,
		ConnectivityScope:network.GlobalScope,
	}, nil
}

/*
NetworkID генерируется LibNetwork, который представляет собой уникальную сеть.

Options является произвольная карта, предоставленная прокси-сервером LibNetwork.

IPv4Data и IPv6Data - это данные IP-адресации, настроенные пользователем и управляемые драйвером IPAM.
	Ожидается, что сетевой драйвер будет учитывать данные IP-адресации, предоставленные драйвером IPAM.
	Данные включают в себя,
		AddressSpace: уникальная строка представляет изолированное пространство для IP-адресации
		Pool: диапазон IP-адресов, представленных в формате адреса / маски CIDR.
			Поскольку драйвер IPAM отвечает за распределение IP-адресов контейнера,
			сетевой драйвер может использовать эту информацию для целей сетевого подключения.
		Gateway: При желании драйвер IPAM может предоставить IP-адрес шлюза в формате CIDR для подсети,
			представленной Пулом. Сетевой драйвер может использовать эту информацию для целей сетевого подключения.
		AuxAddresses: список предварительно назначенных IP-адресов со связанным идентификатором, предоставленным
			пользователем, чтобы помочь сетевому драйверу, если для его работы требуются определенные IP-адреса.
*/
func (self *NetworkDriver) CreateNetwork(r *network.CreateNetworkRequest) error {

	// Берем параметры и проверяем на вшивость
	igData := r.Options[netlabel.GenericData]
	if igData == nil {
		return errors.New("Driver need params vlan_id, bridge_name, ext_if")
	}
	gData := igData.(map[string]interface{})
	if gData == nil {
		return errors.New("Invalid options format")
	}
	iVlanId := gData["vlan_id"]
	iBridgeName := gData["bridge_name"]
	iExtIf := gData["ext_if"]
	if iVlanId == nil { return errors.New("Driver need vlan_id option") }
	if iBridgeName == nil { return errors.New("Driver need bridge_name option") }
	if iExtIf == nil { return errors.New("Driver need ext_if option") }

	vlanId := AnyVal{iVlanId}
	bridgeName := AnyVal{iBridgeName}
	extIf := AnyVal{iExtIf}
	// кончили возиться с параметрами

	n := &Network{}
	if n.Load(r.NetworkID) == nil {
	}
	defer n.Save()

	n.BridgeName = bridgeName.String()
	if r.IPv4Data != nil {
		if len(r.IPv4Data) > 0 {
			n.Gateway = r.IPv4Data[0].Gateway
		}
	}


	br, err := n.GetOrCreateBridge()
	if err != nil {
		return err
	}
	createVlanAndAddToBridge(br, extIf.String(), vlanId.Int())
	return nil
}

func (self *NetworkDriver) AllocateNetwork(r *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	return &network.AllocateNetworkResponse{}, nil
}

func (self *NetworkDriver) DeleteNetwork(r *network.DeleteNetworkRequest) error {
	n := &Network{ID: r.NetworkID}
	return n.Delete()
}

func (self *NetworkDriver) FreeNetwork(r *network.FreeNetworkRequest) error {
	return nil
}

/*
Если удаленному процессу было передано непустое значение в интерфейсе, он должен ответить пустым значением интерфейса.
LibNetwork будет воспринимать это как ошибку, если она предоставит непустое значение, вернет непустое значение и
откатит операцию.
 */
func (self *NetworkDriver) CreateEndpoint(r *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	// подгружаем информацию о сети
	n := &Network{}
	err := n.Load(r.NetworkID)
	defer n.Save()
	if err != nil {
		return nil, err
	}
	// создаем бридж, если он не создан
	br, err := n.GetOrCreateBridge()
	if err != nil {
		return nil, err
	}
	// генерирую основу имени пары интерфейсов
	n.VethName = "veth" + r.NetworkID[0:3] + r.EndpointID[0:3]
	// сздаем парный интерфейс
	vether, err := tenus.NewVethPairWithOptions(n.HostIfName(), tenus.VethOptions{PeerName: n.ContainerIfName()})
	if err != nil {
		log.Println("Error creating veth pair:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			n.Save()
			//remove interfaces by error
			vether.DeletePeerLink()
			vether.DeleteLink()
		}
	}()

	hostIf, err := net.InterfaceByName(n.HostIfName())
	if err != nil {
		return nil, err
	}
	err = br.AddSlaveIfc(hostIf)

	//containerIf, err := net.InterfaceByName(n.ContainerIfName())
	if err != nil {
		log.Println("Error adding veth to bridge:", err)
		return nil, err
	}

	err = vether.SetLinkUp()
	if err != nil {
		log.Println("Error veth host up:", err)
		return nil, err
	}
	err = vether.SetPeerLinkUp()
	if err != nil {
		log.Println("Error veth peer up:", err)
		return nil, err
	}

	//r.Interface.MacAddress = containerIf.HardwareAddr.String()

	return &network.CreateEndpointResponse{Interface: nil}, nil
}

func (self *NetworkDriver) DeleteEndpoint(r *network.DeleteEndpointRequest) error {
	// подгружаем информацию о сети
	n := &Network{}
	err := n.Load(r.NetworkID)
	defer n.Save()
	if err != nil { return err }
	// берем интерфейс
	hostIf, err := net.InterfaceByName(n.HostIfName())
	if err != nil { return err }
	// ьерем бридж
	br, err := n.GetOrCreateBridge()
	if err != nil { return err }
	// удаляем из бриджа хостовый конец
	err = br.RemoveSlaveIfc(hostIf)
	if err != nil { return err }

	// удаляем хостовую часть. Есть подозрение что автоматически удалится и слейв. но надо будет проверить
	err = netlink.NetworkLinkDel(hostIf.Name)
	return err
}

func (self *NetworkDriver) EndpointInfo(r *network.InfoRequest) (*network.InfoResponse, error) {
	return &network.InfoResponse{Value:make(map[string]string)},nil
}

func (self *NetworkDriver) Join(r *network.JoinRequest) (*network.JoinResponse, error) {
	// готовим информацию для контейнера
	res := &network.JoinResponse{StaticRoutes:make([]*network.StaticRoute, 0)}
	n := &Network{}
	err := n.Load(r.NetworkID)
	if err != nil { return nil, err }
	defer n.Save()
	res.InterfaceName.SrcName = n.ContainerIfName()
	res.InterfaceName.DstPrefix = "eth"
	res.DisableGatewayService = false
	ipGateway, _, _ := net.ParseCIDR(n.Gateway)
	res.Gateway = ipGateway.String()
	//res.StaticRoutes = append(res.StaticRoutes, &network.StaticRoute{Destination: "0.0.0.0/0", RouteType: 0, NextHop: ipGateway.String()})

	return res, nil
}

func (self *NetworkDriver) Leave(r *network.LeaveRequest) error {
	return nil
}

/*
Уведомление о данных. Вызывается после GetCapabilities
// DiscoveryType represents the type of discovery element the DiscoverNew function is invoked on
type DiscoveryType int

const (
	// NodeDiscovery represents Node join/leave events provided by discovery
	NodeDiscovery = iota + 1
	// DatastoreConfig represents an add/remove datastore event
	DatastoreConfig
	// EncryptionKeysConfig represents the initial key(s) for performing datapath encryption
	EncryptionKeysConfig
	// EncryptionKeysUpdate represents an update to the datapath encryption key(s)
	EncryptionKeysUpdate
)
*/
func (self *NetworkDriver) DiscoverNew(r *network.DiscoveryNotification) error {
	return nil
}

func (self *NetworkDriver) DiscoverDelete(r *network.DiscoveryNotification) error {
	return nil
}

func (self *NetworkDriver) ProgramExternalConnectivity(r *network.ProgramExternalConnectivityRequest) error {
	return nil
}

func (self *NetworkDriver) RevokeExternalConnectivity(r *network.RevokeExternalConnectivityRequest) error {
	return nil
}


/*
docker network create -d rightipam:latest --subnet=192.168.1.0/24 --gateway=192.168.1.1 --ip-range=192.168.1.4/32 -o vlan_id=2000 -o ext_if=eno1 -o bridge_name=vlan2000 --config-only net1conf
docker network create -d rightipam:latest --scope swarm --config-from net1conf net1
 */

 //      "docker.networkdriver/1.0"