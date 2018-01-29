package main

import (
    "golang.org/x/net/ipv4"
    "fmt"
    "log"
    "net"
    "io"
)

type MulticastReader struct{
    p *ipv4.PacketConn
    c net.PacketConn
}

func (r MulticastReader) Read(b []byte) (n int, err error){
    n, _, _, err = r.p.ReadFrom(b);
    return
}

func (r MulticastReader) Close(){
    log.Println("CLOSE!")
    r.p.Close()
    r.c.Close()
}

func (r *MulticastReader) WriteTo(w io.Writer) (ntot int64, err error){

    var nr, nw int

    for{
        b := make([]byte, 1500)
        nr, err = r.Read(b); if err != nil{
            break
        }
        nw, err = w.Write(b); if err != nil{
            break
        }
        ntot += int64(nw)
    }

	if nr == 0{}

    return ntot, err
}

func NewMulticastReader(addr string, port int) io.Reader{

    var iface *net.Interface
    var group net.IP
    var conn net.PacketConn
    var err error

    ifacename := "wlp1s0"
    group = net.ParseIP(addr)
    iface, err = net.InterfaceByName(ifacename); check(err)
	listenaddr := fmt.Sprintf("%s:%d", addr, port)

	conn, err = net.ListenPacket("udp4", listenaddr); check(err)

    p := ipv4.NewPacketConn(conn)

    err = p.JoinGroup(iface, &net.UDPAddr{IP: group}); check(err)

    return MulticastReader{p, conn}
}

func check(err error){
    if err != nil{
        log.Fatal(err)
    }
}
