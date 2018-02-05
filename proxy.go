package main

import (
    "io"
    "strings"
    "strconv"
    "time"
    "log"
    "net"
    "fmt"
    "sync"
)

type ProxyGroup struct{
    port        int
    reader      MulticastReader
    writers     map[string][]chan []byte
    rtpseq      map[string]uint16
}

func newProxyGroup(port int, reader MulticastReader) ProxyGroup{
    g := ProxyGroup{}
    g.port = port
    g.reader = reader
    g.writers = make(map[string][]chan []byte)
    g.rtpseq = make(map[string]uint16)

    return g
}

func (g *ProxyGroup) String() string{
    s := fmt.Sprintf("Group %d: ", g.port)
    for addrinfo, writers := range g.writers{
        s += fmt.Sprintf("/ %s: %d", addrinfo, len(writers))
    }

    return s
}


type Proxy struct{
    sync.RWMutex
    groups map[int]ProxyGroup
}

func NewProxy() *Proxy{
    p := &Proxy{}
    p.groups = make(map[int]ProxyGroup)

    return p
}

func (p *Proxy) PrintGroups(){
    log.Println("Groups:", len(p.groups))
    for _, g := range p.groups{
        log.Printf(g.String())
    }
}

func (p *Proxy) RegisterReader(addrinfo string){
    p.Lock()
    ip, port := addrinfoToIPPort(addrinfo)
	g, ok := p.groups[port]; if !ok{
		r := NewMulticastReader("eno1", ip, port)
        log.Println("Registered new reader for", port, r)
        g = newProxyGroup(port, r)
        log.Printf("Created new proxy group for %d, with reader %s", port, r)
		p.groups[port] = g
	}else{
        log.Println("Reader already exists for", port, g, g.reader)
        _, alreadyjoined := g.writers[addrinfo]; if alreadyjoined{ //TODO: this should be handled in the reader
            return
        }
    }
    p.Unlock()
    g.reader.JoinGroup(net.ParseIP(ip))
}

func (p *Proxy) RegisterWriter(addrinfo string, c chan []byte){
    p.Lock()
    _, port := addrinfoToIPPort(addrinfo)
    g, ok := p.groups[port]; if !ok{
        p.RegisterReader(addrinfo)
    }
    _, ok = g.writers[addrinfo]; if !ok{
        g.writers[addrinfo] = make([]chan []byte, 0)
    }

    g.writers[addrinfo] = append(g.writers[addrinfo], c)
    log.Println("Registered new writer for", g.port, addrinfo, c)

    p.PrintGroups()

    p.Unlock()
}

func (p *Proxy) RemoveWriter(addrinfo string, c chan[]byte){
    p.Lock()

    ip, port := addrinfoToIPPort(addrinfo)

    g, ok := p.groups[port]; if !ok{
        log.Println("WARNING: No reader for group", port)
    }

    for i, wc := range g.writers[addrinfo]{
        if c == wc{
            g.writers[addrinfo] = append(g.writers[addrinfo][:i], g.writers[addrinfo][i+1:]...)
            log.Println("Removed writer", c)
        }
    }

    if len(g.writers[addrinfo]) == 0{
        log.Println("No more writers left for", addrinfo, "closing reader")
        delete(g.writers, addrinfo)

        g.reader.LeaveGroup(net.ParseIP(ip))

        if len(g.writers) == 0{
            log.Println("No more readers in group", port, ". Removing group")
            g.reader.Close()
            delete(p.groups, port)
        }

    }

    p.PrintGroups()

    p.Unlock()
}

func (p *Proxy) CloseGroup(port int){
    log.Println("Close group", port)
    g := p.groups[port]
    log.Println("Found group for", port, ". Remove all writers")
    g.reader.Close()
    for addrinfo, chans := range g.writers{
        for _, c := range chans{
            close(c)
            log.Println("Writer", addrinfo, "closed")
        }
    }
    delete(p.groups, port)
}

func (p *Proxy) Loop(){
    waitintvl := 100 * time.Millisecond
    for{
        time.Sleep(waitintvl)
        for{
            if len(p.groups) == 0{ break }
            p.Lock()
            for port, g := range p.groups{
                g.reader.SetReadDeadline(time.Now().Add(2500 * time.Millisecond))
                rtp, dst, err := ReadRTP(g.reader); if err != nil{
                    if err == io.EOF{
                        log.Println("EOF from", g.reader)
                    }else{
                        operr := err.(*net.OpError); if operr != nil && operr.Timeout(){
                            log.Println("Timeout from", g.reader)
                        }
                    }
                    p.CloseGroup(port)
                    break
                }

                if rtp == nil{ continue }

                addrinfo := fmt.Sprintf("%s:%d", dst, port)

                if g.rtpseq[addrinfo] + 2 != rtp.SequenceNumber - 0xffff{
                    log.Println("BAD RTPSEQ", addrinfo,
                        "expecting", g.rtpseq[addrinfo] + 2,
                        "have", rtp.SequenceNumber, rtp.SequenceNumber - 0xffff)
                }

                for _, c := range g.writers[addrinfo]{
                    c <- rtp.Payload
                }
                g.rtpseq[addrinfo] = rtp.SequenceNumber
            }
            p.Unlock()
        }
    }
}

func addrinfoToIPPort(addrinfo string) (ip string, port int){
    ipport := strings.Split(addrinfo, ":")
    if len(ipport) != 2{
        log.Println("BAD ADDRESS", addrinfo)
        return
    }

    ip = ipport[0]
    port, err := strconv.Atoi(ipport[1]); if err != nil{
        log.Println(err, ipport)
        return
    }

    return
}

func readToChannel(addrinfo string, r io.Reader, c chan []byte){
}
