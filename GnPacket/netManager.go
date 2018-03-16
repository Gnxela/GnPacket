package GnPacket

import (
	"sync"
	"encoding/binary"
)

type NetManager struct {
	UnhandledQueue chan GnPacket
	data []byte
	mutex sync.Mutex
}

func New(queueLength int) NetManager {
	netManager := NetManager{
		make(chan GnPacket, queueLength),
		make([]byte, 0),
		sync.Mutex{},
	}
	
	return netManager
	
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
		var id int16 = int16(packetId[0]) * 255 + int16(packetId[1])
		
		packet := GnPacket{id, data[6:]}
		netManager.UnhandledQueue <- packet;
	}
	netManager.mutex.Unlock()
}

func (netManager *NetManager) HasUnhandledPacket() bool {
	return len(netManager.UnhandledQueue) > 0
}