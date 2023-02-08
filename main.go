package main

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var records = map[string]string{
	"go.go":      "127.0.0.1",
	"amazon.com": "176.32.103.205",
}

func serve(u *net.UDPConn, addr net.Addr, request *layers.DNS) {
	log.Printf("query params: [name: %s, type: %s, class: %s] callee (%s)%s\n",
		request.Questions[0].Name, request.Questions[0].Type, request.Questions[0].Class, addr.Network(), addr.String())
	response := request
	response.QR = true
	response.ANCount = 1
	response.AA = true
	response.OpCode = layers.DNSOpCodeQuery
	response.QDCount = request.QDCount
	ip, ok := records[string(request.Questions[0].Name)]
	if !ok {
		response.ResponseCode = layers.DNSResponseCodeServFail
	} else {
		response.ResponseCode = layers.DNSResponseCodeNoErr
		response.Answers = []layers.DNSResourceRecord{
			{
				Name:  request.Questions[0].Name,
				Type:  request.Questions[0].Type,
				Class: request.Questions[0].Class,
				TTL:   3600,
				IP:    net.ParseIP(ip),
			},
		}
	}
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, response)
	u.WriteTo(buf.Bytes(), addr)
}

func main() {
	u, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("127.0.0.1"),
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("start server...")
		for {
			buf := make([]byte, 1024)
			_, addr, _ := u.ReadFrom(buf)
			pakcet := gopacket.NewPacket(buf, layers.LayerTypeDNS, gopacket.Default)
			request := pakcet.Layer(layers.LayerTypeDNS).(*layers.DNS)
			serve(u, addr, request)
		}
	}()

	// go func() {
	// 	glink()
	// }()

	select {}
}

// var glinks = map[string]string{
// 	"meet": "https://meet.google.com/lookup/xxxx",
// }

// func glink() {
// 	r := gin.New()
// 	r.GET("/:s", func(c *gin.Context) {
// 		link := c.Param("s")
// 		if v, ok := glinks[link]; ok {
// 			c.Redirect(301, v)
// 		} else {
// 			c.String(404, "not found")
// 		}
// 	})
// 	r.Run(":80")
// }
