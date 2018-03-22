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
	netManager.AddHandler(2, handleMessage2)
	
	go play(&netManager)
	
	time.Sleep(time.Second * 5)
	
	netManager.RemoveHandler(2, handleMessage)
	
	<-netManager.UnhandledQueue//Just hold the program open, we'll never recieve an unhandled packet currently
}

func play(netManager *GnPacket.NetManager) {
	var data []byte
	for {
		fmt.Printf("Data length: %d\n", len(data))
	
		ping := NewPacketPing()
		data = append(data, ping.Write(&ping)[:]...)
		
		message := NewPacketMessage("Hello World!")
		data = append(data, message.Write(&message)[:]...)

		netManager.ReadData(&data)

		time.Sleep(time.Second)
	}
}

func handlePing(packet GnPacket.GnPacket) bool {
	ping := PacketPing{&packet, time.Now()}
	ping.Deserialize(packet.Data)
	fmt.Printf("%v\n", time.Now().Sub(ping.Start))
	return true
}

var i int = 0;

func handleMessage(packet GnPacket.GnPacket) bool {
	message := PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	fmt.Printf("%v%d\n", message.Message, i)
	i++
	return false
}

func handleMessage2(packet GnPacket.GnPacket) bool {
	message := PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	fmt.Printf("%v appended this!\n", message.Message)
	return true
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