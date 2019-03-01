package main

import (
	"flag"
	"log"
	"net"

	"github.com/tritonuas/god/go-mavlink/mavlink"
)

//////////////////////////////////////
//
// mavlink udp rx example
//
// listen to ardupilot SITL and prints received msgs, more info:
// http://dev.ardupilot.com/wiki/simulation-2/sitl-simulator-software-in-the-loop/setting-up-sitl-on-linux/
//
// run via `go run main.go`
//
//////////////////////////////////////

var rxaddr = flag.String("addr", ":5760", "address to listen on")

func main() {

	flag.Parse()

	listenAndServe(*rxaddr)
}

func listenAndServe(addr string) {

	//udpAddr, err := net.ResolveUDPAddr("udp", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	//conn, listenerr := net.ListenUDP("udp", udpAddr)
	//if listenerr != nil {
	//	log.Fatal(listenerr)
	//}

	//log.Println("listening on", udpAddr)

	dec := mavlink.NewDecoder(conn)
	dec.Dialects.Add(mavlink.DialectArdupilotmega)

	enc := mavlink.NewEncoder(conn)
	p := mavlink.MessageInterval{
	    IntervalUs: 10,
	    MessageId: 33,
	}

	if err := enc.Encode(0x0, 0x0, &p); err != nil {
	    log.Fatal("Encode fail:", err)
	}

	for {
		pkt, err := dec.Decode()
		if err != nil {
			//log.Println("Decode fail:", err)
			//if pkt != nil {
			//	log.Println(*pkt)
			//}
			continue
		}

		//log.Println("msg rx:", pkt.MsgID, pkt.Payload)
		switch pkt.MsgID {
		    case mavlink.MSG_ID_HEARTBEAT:
			var pv mavlink.Heartbeat
			if err := pv.Unpack(pkt); err == nil {
			    // handle param value
			    log.Println("Heartbeat");
			}
		case mavlink.MSG_ID_GLOBAL_POSITION_INT:
			var gv mavlink.GlobalPositionInt
			if err := gv.Unpack(pkt); err == nil {
				log.Println("GPS");
			}
		case mavlink.MSG_ID_VFR_HUD:
			var vh mavlink.VfrHud
			if err := vh.Unpack(pkt);err==nil {
				log.Println("VFR");
			}		
		case mavlink.MSG_ID_SYS_STATUS:
			var ss mavlink.SysStatus
			if err := ss.Unpack(pkt); err == nil{
				log.Println("SS");
			}
		}
	}
}
