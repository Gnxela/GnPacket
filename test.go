package main

import (
	"time"
	"./GnPacket"
	
	"fmt"
	"math/rand"
)

type PacketPing struct {
	*GnPacket.GnPacket
	Start time.Time
}

func main() {
	netManager := GnPacket.New(100)

	go play(&netManager)
	
	for {
		packet := <- netManager.UnhandledQueue
		if (packet.Id == 1) {
			ping := PacketPing{&packet, time.Now()}
			ping.Deserialize(packet.Data)
			fmt.Printf("%v\n", time.Now().Sub(ping.Start));
		}
	}
}

func play(netManager *GnPacket.NetManager) {
	for {
		packet := NewPacketPing();
		fmt.Printf("%v\n", packet.Start);
		data := packet.Write(packet);
		
		length := len(data)
		cut := rand.Intn(length + 1)
		
		fmt.Printf("Cutting at %d, length %d\n", cut, length)
		netManager.Feed(data[:cut])

		netManager.Feed(data[cut:])
		time.Sleep(time.Second / 4)
	}
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
	packet.Start.GobDecode(data[6:])
}