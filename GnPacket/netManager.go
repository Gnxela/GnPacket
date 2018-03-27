package GnPacket

import (
	"encoding/binary"
	
	"fmt"
)

// NetManager helps manage the reading and writing of packets.
// Should only be used over lossless and order garenteed protocols.
type NetManager struct {
	UnhandledQueue chan GnPacket
	handlers map[uint16][]func(packet GnPacket) bool
}


// New returns a NetManager instance. 
// Ready to read data and fire listeners without further setup (adding listeners required)
func New(queueLength int) NetManager {
	netManager := NetManager{
		make(chan GnPacket, queueLength),
		make(map[uint16][]func(packet GnPacket) bool),
	}
	
	return netManager
}

// AddHandler adds a new handler to netManager listening for packet id.
// The same handler can be added multible times, to the same or different IDs.
// If a handler returns false, all subsiquent handlers will not fire.
func (netManager *NetManager) AddHandler(id uint16, handler func(packet GnPacket) bool) {
	netManager.handlers[id] = append(netManager.handlers[id], handler);
}

// RemoveHandler removes a handler from netManager.
// If the handler is not present nothing happens.
// If a handler has been added multible times, all instances listening for id are removed.
func (netManager *NetManager) RemoveHandler(id uint16, handler func(packet GnPacket) bool) {
	if handlers, ok := netManager.handlers[id]; ok {
		for i, handle := range handlers {
			if fmt.Sprintf("%v", handle) == fmt.Sprintf("%v", handler) {
				netManager.handlers[id] = append(netManager.handlers[id][:i], netManager.handlers[id][i + 1:]...);
			}
		}
	}
}

// ReadData reads the data and parses any packets inside it. 
// Any parsed packets are removed from the original data.
// ReadData automatically fires relivant listeners from netManager when a packet is recieved.
// If there are no listeners for a packet, it is added to netManager.UnhandledQueue.
func (netManager *NetManager) ReadData(data *[]byte) {
	for {
		if len(*data) < 4 {
			break
		}
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

// HasUnhandledPacket returns true if netManager.UnhandledQueue has unhandled packets.
func (netManager *NetManager) HasUnhandledPacket() bool {
	return len(netManager.UnhandledQueue) > 0
}