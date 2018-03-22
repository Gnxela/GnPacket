package GnPacket

import (
	"encoding/binary"
)

type GnPacket struct {
	Id uint16
	Data []byte
}

func (packet *GnPacket) Write(writable PacketWritable) []byte {
	data := writable.Serialize();
	
	if (len(data) > 4294967295) {
		panic("Data too large!")
	}
	
	packetId := make([]byte, 2);
	packetId[0] = byte(packet.Id / 255)
	packetId[1] = byte(packet.Id % 255)
	
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data) + 6))
	
	return append(append(length, packetId[:]...), data[:]...)
}

type PacketWritable interface {
	Serialize() []byte
}

type PacketReadable interface {
	Deserialize(data []byte)
}