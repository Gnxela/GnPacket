package GnPacket

import (
	"sync"
	"encoding/binary"
	
	"fmt"
)

type NetManager struct {
	UnhandledQueue chan GnPacket
	handlers map[uint16][]func(packet GnPacket) bool
	data []byte
	mutex sync.Mutex
}

func New(queueLength int) NetManager {
	netManager := NetManager{
		make(chan GnPacket, queueLength),
		make(map[uint16][]func(packet GnPacket) bool),
		make([]byte, 0),
		sync.Mutex{},
	}
	
	return netManager
	
}

func (netManager *NetManager) AddHandler(id uint16, handler func(packet GnPacket) bool) {
	netManager.handlers[id] = append(netManager.handlers[id], handler);
}

func (netManager *NetManager) RemoveHandler(id uint16, handler func(packet GnPacket) bool) {
	if handlers, ok := netManager.handlers[id]; ok {
		for i, handle := range handlers {
			if fmt.Sprintf("%v", handle) == fmt.Sprintf("%v", handler) {
				netManager.handlers[id] = append(netManager.handlers[id][:i], netManager.handlers[id][i + 1:]...);
			}
		}
	}
}

func (netManager *NetManager) Feed(data []byte) {
	netManager.mutex.Lock()
	
	netManager.data = append(netManager.data, data...)
	
	packetLength := netManager.data[:4];
	length := binary.LittleEndian.Uint32(packetLength)
	
	if (len(netManager.data) >= int(6 + length)) {
		//We have a completed packet
		data := netManager.data[:6 + length]
		netManager.data = netManager.data[6 + length:]//Remove the data of the packet
		
		packetId := data[4:6];
		var id uint16 = uint16(packetId[0]) * 255 + uint16(packetId[1])
		
		packet := GnPacket{id, data[6:]}
		
		if handlers, ok := netManager.handlers[id]; ok {
			for _, handler := range handlers {
				if !handler(packet) {
					break
				}
			}
		} else {
			netManager.UnhandledQueue <- packet;
		}
	}
	netManager.mutex.Unlock()
}

func (netManager *NetManager) HasUnhandledPacket() bool {
	return len(netManager.UnhandledQueue) > 0
}