package networkreudp

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"

	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/snuk182/reudp"
)

const (
	bufferSize       = 8192
	portKey          = "udpPort"
	serverAddressKey = "serverAddress"
	defaultPort      = 8787
)

var (
	header                = binary.LittleEndian.Uint32([]byte{8, 8, 8, 8})
	version               = binary.LittleEndian.Uint32([]byte{0, 0, 0, 1})
	timeout time.Duration = 4000
	portTkv               = api.TypeKeyValue{
		Type:  reflect.Int,
		Key:   portKey,
		Name:  "A port to run server on (if being a server) or to connect to (if being a client).",
		Value: defaultPort,
	}
	hostTkv = api.TypeKeyValue{
		Type: reflect.String,
		Key:  serverAddressKey,
		Name: "An address of server to connect to. No server address means being the server itself.",
	}
)

type udpPacket struct {
	who  *net.UDPAddr
	what []byte
}

type reudpNetwork struct {
	mxSend, mxConnections, mxPackets sync.RWMutex

	port        int
	conn        *reudp.Reudp
	address     *net.UDPAddr
	socket      *net.UDPConn
	packets     []*udpPacket
	writeBuffer [bufferSize]byte
	closed      bool
}

func newRaw() *reudpNetwork {
	return &reudpNetwork{
		port: defaultPort,
	}
}

func NewServer() modules.Network {
	return newRaw()
}

func NewClient(host string) modules.Network {
	hostTkv.Value = host

	n := newRaw()
	n.SetConfiguration(hostTkv)

	return n
}

func (this *reudpNetwork) ID() string {
	return "reudpNetwork"
}

func (this *reudpNetwork) Start() (err error) {
	var addr *net.UDPAddr
	if this.address != nil {
		log.Debugf("Start client connection to %v", this.address)
		addr = &net.UDPAddr{
			IP:   this.address.IP,
			Port: 0,
		}
	} else {
		log.Debugf("Start server connection at %v", this.port)
		addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", this.port))
		if err != nil {
			return err
		}
	}

	mh := func(who *net.UDPAddr, data []byte) { this.onMessage(who, data) }
	eh := func(who *net.UDPAddr, err error) { this.onError(who, err) }

	this.conn = reudp.New(mh, eh, onLog)

	err = this.conn.Listen(addr)
	if err != nil {
		return err
	}

	log.Infof("Listening on %s", this.conn.LocalAddr())

	this.closed = false

	return nil
}

func (this *reudpNetwork) onMessage(who *net.UDPAddr, data []byte) {
	log.Tracef("from %v : %v bytes", who, len(data))

	hdr := binary.LittleEndian.Uint32(data[:4])
	ver := binary.LittleEndian.Uint32(data[4:8])

	if hdr == header {
		switch ver {
		case version:
			p := new(udpPacket)
			p.who = who
			p.what = make([]byte, len(data)-16)
			copy(p.what, data[16:])

			this.mxPackets.Lock()
			this.packets = append(this.packets, p)
			this.mxPackets.Unlock()
		}
	} else {
		log.Warnf("not a packet: %v", data[:4])
		return
	}
}

func (this *reudpNetwork) onError(who *net.UDPAddr, err error) {
	log.Errorf("from %v: %v", err)
}

func onLog(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (this *reudpNetwork) Stop() {
	this.closed = true

	if this.conn != nil {
		this.conn.Close()
	}

	if this.socket != nil {
		log.Debugf("reudpNetwork.Stop %v: error: %v", this.port, this.socket.Close())
		//this.socket = nil
	}
}

func (this *reudpNetwork) Where() interface{} {
	//log.Tracef("reudpNetwork.Where: %#v", this.conn.LocalAddr())

	if this.address != nil {
		return this.address
	} else {
		return this.conn.LocalAddr()
	}
}

func (this *reudpNetwork) Configuration() []api.TypeKeyValue {
	hostTkv.Value = this.address.String()
	portTkv.Value = this.port

	tkvs := []api.TypeKeyValue{hostTkv, portTkv}
	log.Tracef("%#v", tkvs)
	return tkvs
}

func (this *reudpNetwork) SetConfiguration(values ...api.TypeKeyValue) error {
	log.Tracef("%#v", values)

	addr := ""
	for _, v := range values {
		if v.Value == nil {
			continue
		}

		switch v.Key {
		case portKey:
			this.port = v.Value.(int)
		case serverAddressKey:
			addr = v.Value.(string)
		}
	}

	var e error

	if addr != "" {
		this.address, e = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, this.port))
	}

	return e
}

func (this *reudpNetwork) Send(where interface{}, b []byte, reliable bool) error {
	log.Tracef("to %v: %v bytes", where, len(b))

	if addr, ok := where.(*net.UDPAddr); ok {
		this.mxSend.Lock()

		packetLength := len(b) + 16

		var buffer []byte
		if len(this.writeBuffer) < packetLength {
			buffer = make([]byte, packetLength)
		} else {
			buffer = this.writeBuffer[:packetLength]
		}

		binary.LittleEndian.PutUint32(buffer[:4], header)
		binary.LittleEndian.PutUint32(buffer[4:8], version)
		copy(buffer[16:], b)
		err := this.conn.Send(addr, buffer)
		log.Debugf("sent %d bytes to %v (%v)", len(buffer), addr, buffer[:16])

		/*if packetLength > 500 {
			<-time.Tick(time.Millisecond * time.Duration(packetLength / 15))
		}*/

		this.mxSend.Unlock()
		return err
	} else {
		return errors.New(fmt.Sprintf("unsupported address type: %#v", where))
	}
}

func (this *reudpNetwork) frame() {
	this.packets = make([]*udpPacket, 0, 8192)
}

func (this *reudpNetwork) Poll() ([]interface{}, [][]byte, error) {
	//defer this.frame()

	this.mxPackets.Lock()

	addrs := make([]interface{}, len(this.packets))
	data := make([][]byte, len(this.packets))

	for i, v := range this.packets {
		addrs[i] = v.who
		data[i] = v.what
	}

	this.frame()

	this.mxPackets.Unlock()

	return addrs, data, nil
}
