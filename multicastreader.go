package main

import (
    "golang.org/x/net/ipv4"
    "fmt"
    "log"
    "net"
    "io"
    "time"
)

type MulticastReader struct{
    p *ipv4.PacketConn
    c net.PacketConn
    Port    int
}

func (r MulticastReader) ReadFrom(b []byte) (n int, cm *ipv4.ControlMessage, src net.Addr, err error){
    n, cm, src, err = r.p.ReadFrom(b);
    return
}

func (r MulticastReader) Close() (err error){
    log.Println("CLOSE!")
    r.p.Close()
    r.c.Close()
    return nil
}

func (r *MulticastReader) WriteTo(w io.Writer) (ntot int64, err error){

    var nr, nw int

    for{
        b := make([]byte, 1500)
        nr, _, _, err = r.ReadFrom(b); if err != nil{
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

func (r *MulticastReader) SetDeadline(t time.Time) error{
    return r.p.SetDeadline(t)
}
func (r *MulticastReader) SetReadDeadline(t time.Time) error{
    return r.p.SetReadDeadline(t)
}

func NewMulticastReader(addr string, port int) MulticastReader{

    var iface *net.Interface
    var group net.IP
    var conn net.PacketConn
    var err error

    ifacename := "eno1"
    group = net.ParseIP(addr)
    iface, err = net.InterfaceByName(ifacename); check(err)
	listenaddr := fmt.Sprintf("%s:%d", addr, port)

    log.Println("Listen", listenaddr)

	conn, err = net.ListenPacket("udp4", listenaddr); check(err)

    log.Println("Conn", conn, conn.LocalAddr())

    p := ipv4.NewPacketConn(conn)

    err = p.JoinGroup(iface, &net.UDPAddr{IP: group}); check(err); if err != nil{
        log.Println("Cannot join group", group, err)
    }
    err = p.SetControlMessage(ipv4.FlagDst, true); if err != nil {
        log.Println("Cannot SetControlMessage FlagDst")
    }

    return MulticastReader{p, conn, port}
}

func check(err error){
    if err != nil{
        log.Fatal(err)
    }
}
