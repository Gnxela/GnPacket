package main

import (
	"time"
	"./GnPacket"
	
	"fmt"
)

type PacketPing struct {
	*GnPacket.GnPacket
	Start time.Time
}

type PacketMessage struct {
	*GnPacket.GnPacket
	Message string
}

func main() {
	netManager := GnPacket.New(100)

	go play(&netManager)
	
	for {
		packet := <- netManager.UnhandledQueue
		if (packet.Id == 1) {
			ping := PacketPing{&packet, time.Now()}
			ping.Deserialize(packet.Data)
			fmt.Printf("%v\n", time.Now().Sub(ping.Start))
		} else if (packet.Id == 2) {
			message := PacketMessage{&packet, ""}
			message.Deserialize(packet.Data)
			fmt.Printf("%v\n", message.Message)
		}
	}
}

func play(netManager *GnPacket.NetManager) {
	for {
		ping := NewPacketPing();
		data := ping.Write(ping);
		netManager.Feed(data)
		
		message := NewPacketMessage("Hello World!");
		data = message.Write(message);
		netManager.Feed(data)
	}
}

func NewPacketMessage(message string) PacketMessage {
	return PacketMessage{&GnPacket.GnPacket{2, make([]byte, 0)}, message};
}

func (packet PacketMessage) Serialize() []byte {
	return []byte(packet.Message)
}

func (packet *PacketMessage) Deserialize(data []byte) {
	packet.Message = string(data)
}

func NewPacketPing() PacketPing {
	return PacketPing{&GnPacket.GnPacket{1, make([]byte, 0)}, time.Now()};
}

func (packet PacketPing) Serialize() []byte {
	data, err := packet.Start.GobEncode()
	if (err != nil) {
		panic(err)
	}
	return data
}

func (packet *PacketPing) Deserialize(data []byte) {
	packet.Start.GobDecode(data)
}