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
	
	netManager.AddHandler(1, handlePing)
	netManager.AddHandler(2, handleMessage)
	netManager.AddHandler(2, handleMessage)
	
	go play(&netManager)
	
	<-netManager.UnhandledQueue//Just hold the program open, we'll never recieve an unhandled packet currently
}

func play(netManager *GnPacket.NetManager) {
	for {
		ping := NewPacketPing();
		data := ping.Write(&ping);
		netManager.Feed(data)
		
		message := NewPacketMessage("Hello World!");
		data = message.Write(&message);
		netManager.Feed(data)
		time.Sleep(time.Second)
	}
}

func handlePing(packet GnPacket.GnPacket) {
	ping := PacketPing{&packet, time.Now()}
	ping.Deserialize(packet.Data)
	fmt.Printf("%v\n", time.Now().Sub(ping.Start))
}

var i int = 0;

func handleMessage(packet GnPacket.GnPacket) {
	message := PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	fmt.Printf("%v%d\n", message.Message, i)
	i++
}

func NewPacketMessage(message string) PacketMessage {
	return PacketMessage{&GnPacket.GnPacket{2, make([]byte, 0)}, message};
}

func (packet *PacketMessage) Serialize() []byte {
	return []byte(packet.Message)
}

func (packet *PacketMessage) Deserialize(data []byte) {
	packet.Message = string(data)
}

func NewPacketPing() PacketPing {
	return PacketPing{&GnPacket.GnPacket{1, make([]byte, 0)}, time.Now()};
}

func (packet *PacketPing) Serialize() []byte {
	data, err := packet.Start.GobEncode()
	if (err != nil) {
		panic(err)
	}
	return data
}

func (packet *PacketPing) Deserialize(data []byte) {
	packet.Start.GobDecode(data)
}