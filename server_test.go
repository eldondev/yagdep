package main

import "testing"
import "net"
import "time"

func TestServe(t *testing.T) {
	test_body := []byte { 0x67,0x52,0x01,0x00,0x00,0x01,0x00,0x00,0x00,0x00,0x00,0x00,0x06,0x67,0x6f,0x6c,0x61,0x6e,0x67,0x03,0x6f,0x72,0x67,0x00,0x00,0x01,0x00,0x01 }
	test_req_port, _ := net.ListenUDP("udp",&net.UDPAddr{net.IP{127,0,0,0},0,""}) 
	test_resp_port, _ := net.ListenUDP("udp",&net.UDPAddr{})
	test_req_addr, _ := net.ResolveUDPAddr("udp", test_req_port.LocalAddr().String())
	packet := make([]byte, 65507)
	var byte_count int
	waiter := make(chan interface{}); 
	config = Config{Hosts: map[string][]net.IP{"golang.org.":{net.IP{127,0,0,1}}}}
	go func() { byte_count, _, _ = test_req_port.ReadFrom(packet) ; waiter <- nil }();
	serve(len(test_body), test_req_addr, test_body, test_resp_port)
	select {
		case  <- waiter:
			t.Logf("Response received %+v", packet[:byte_count]);
		case <- time.After(2 * time.Second):
			t.Logf("Timed out")
			t.Fail();
	}
}
