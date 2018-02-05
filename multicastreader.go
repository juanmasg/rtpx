package main

import (
    "golang.org/x/net/ipv4"
    "fmt"
    "log"
    "net"
    //"io"
    "time"
)

type MulticastReader struct{
    p *ipv4.PacketConn
    c net.PacketConn
    Port    int
    iface   *net.Interface
}

func (r MulticastReader) String() string{
    return fmt.Sprintf("MulticastReader/%s/0.0.0.0:%d", r.iface.Name, r.Port)
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

func (r *MulticastReader) SetDeadline(t time.Time) error{
    return r.p.SetDeadline(t)
}
func (r *MulticastReader) SetReadDeadline(t time.Time) error{
    return r.p.SetReadDeadline(t)
}

func (r *MulticastReader) JoinGroup(group net.IP) error{
    return r.p.JoinGroup(r.iface, &net.UDPAddr{IP:group})
}

func (r *MulticastReader) LeaveGroup(group net.IP) error{
    return r.p.LeaveGroup(r.iface, &net.UDPAddr{IP:group})
}

func NewMulticastReader(ifacename string, addr string, port int) MulticastReader{

    var err error

    r := MulticastReader{}
    r.Port = port

    r.iface, err = net.InterfaceByName(ifacename); check(err)
	listenaddr := fmt.Sprintf("%s:%d", addr, port)

    log.Println("Listen", listenaddr)

	r.c, err = net.ListenPacket("udp4", listenaddr); check(err)

    r.c.(*net.UDPConn).SetReadBuffer(4 * 1024 * 1024)

    log.Println("Conn", r.c, r.c.LocalAddr())

    r.p = ipv4.NewPacketConn(r.c)

    err = r.p.SetControlMessage(ipv4.FlagDst, true); if err != nil {
        log.Println("Cannot SetControlMessage FlagDst")
    }

    return r
}

func check(err error){
    if err != nil{
        log.Fatal(err)
    }
}
