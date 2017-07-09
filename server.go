package main

import "net"
import "io"
import "log"
import "bufio"
import "bytes"
import "encoding/binary"

type RRecord struct {
	Name []string
	RRecordFooter
	NameBuffer []byte
}

type RRecordFooter struct {
	Type  uint16
	Class uint16
}

type DNSPacketHeader struct {
	Id          uint16
	Flags       [2]uint8
	Questions   uint16
	Answers     uint16
	Authorities uint16
	Additionals uint16
}

type Request struct {
	DNSPacketHeader
	records []RRecord
}

func (d *DNSPacketHeader) total_records() int {
	return int(d.Questions + d.Answers + d.Authorities + d.Additionals)
}

func (record *RRecord) read_record_names(name_reader *bufio.Reader) (err error) {
	var name []byte
	for {
		var name_length byte
		if name_length, err = name_reader.ReadByte(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%+v", err)
		}
		if name, err = name_reader.Peek(int(name_length)); err != nil {
			log.Fatalf("%+v", err)
		}
		name_reader.Discard(int(name_length))
		record.Name = append(record.Name, string(name))
	}
	return
}

func (request *Request) read_record(packet_reader *bufio.Reader) (err error) {
	for record_index := 0; record_index < request.total_records(); record_index++ {
		record := RRecord{}
		var name []byte
		var err error
		if name, err = packet_reader.ReadBytes('\x00'); err != nil {
			log.Fatalf("%+v", err)
		}
		name_reader := bufio.NewReader(bytes.NewReader(name))
		record.read_record_names(name_reader)
		err = binary.Read(packet_reader, binary.BigEndian, &record.RRecordFooter)
		request.records = append(request.records, record)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
	return
}

func respond(request Request, requester net.Addr) (err error) {
	return nil
}

func serve(byte_count int, requester net.Addr, packet []byte) (err error) {
	request := Request{}
	packet_reader := bufio.NewReader(bytes.NewReader(packet))
	if err = binary.Read(packet_reader, binary.BigEndian, &request.DNSPacketHeader); err == nil {
		request.read_record(packet_reader)
		err = respond(request, net.Addr)
	}
	if err != nil {
		log.Fatalf("%+v", err)
		log.Println("binary.Read failed:", err)
	}
	return
}

func main() {
	log.SetFlags(log.Flags() | log.Llongfile)
	server, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	packet := make([]byte, 65507)
	for {
		var byte_count int
		var requester net.Addr
		byte_count, requester, err = server.ReadFrom(packet)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		serve(byte_count, requester, packet)
	}
}
