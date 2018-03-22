package GnPacket

import (
	"encoding/binary"
	
	"fmt"
)

type NetManager struct {
	UnhandledQueue chan GnPacket
	handlers map[uint16][]func(packet GnPacket) bool
}

func New(queueLength int) NetManager {
	netManager := NetManager{
		make(chan GnPacket, queueLength),
		make(map[uint16][]func(packet GnPacket) bool),
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

func (netManager *NetManager) ReadData(data *[]byte) {
	for {
		packetLength := (*data)[:4];
		length := binary.LittleEndian.Uint32(packetLength)
		if (len(*data) >= int(length) && length > 0) {//If we have a completed packet

			packetData := (*data)[:length]
			*data = (*data)[length:]//Remove the data of the packet
			
			packetId := packetData[4:6];
			var id uint16 = uint16(packetId[0]) * 255 + uint16(packetId[1])
			
			packet := GnPacket{id, packetData[6:]}
						
			if handlers, ok := netManager.handlers[id]; ok {
				for _, handler := range handlers {
					if !handler(packet) {
						break
					}
				}
			} else {
				netManager.UnhandledQueue <- packet;
			}
		} else {
			break
		}
	}
}

func (netManager *NetManager) HasUnhandledPacket() bool {
	return len(netManager.UnhandledQueue) > 0
}