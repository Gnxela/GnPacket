package GnPacket

import (
	"encoding/binary"
)

type GnPacket struct {
	Id int16
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
	binary.LittleEndian.PutUint32(length, uint32(len(data)))
	
	return append(append(length, packetId[:]...), data[:]...)
}

func (packet *GnPacket) Read(data []byte) {
	packetLength := data[:4];
	length := binary.LittleEndian.Uint32(packetLength)
	
	packetId := data[4:6];
	id := packetId[0] * 255 + packetId[1]
	
	packetData := data[6:6+length];
	
	_, _ = packetData, id
}

type PacketWritable interface {
	Serialize() []byte
}

type PacketReadable interface {
	Deserialize(data []byte)
}