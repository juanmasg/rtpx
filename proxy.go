package main

import (
    "io"
    "strings"
    "strconv"
    "time"
    "log"
    "net"
)

type Proxy struct{
    readers map[string]io.Reader
    writers map[string][]chan []byte
    rtpseq map[string]uint16
}

func NewProxy() *Proxy{
    p := &Proxy{}
    p.readers = make(map[string]io.Reader)
    p.writers = make(map[string][]chan []byte)
    p.rtpseq = make(map[string]uint16)

    return p
}

func (p *Proxy) RegisterReader(addrinfo string){
    ip, port := addrinfoToIPPort(addrinfo)
	r, ok := p.readers[addrinfo]; if !ok{
		r = NewMulticastReader(ip, port)
		p.readers[addrinfo] = r
        log.Println("Registered new reader for", addrinfo, r)
	}
}

func (p *Proxy) RegisterReader2(addrinfo string){
    r, ok := p.readers[addrinfo]; if !ok{
	    addr, err := net.ResolveUDPAddr("udp", addrinfo); if err != nil {
	        log.Fatal(err)
	    }
	    r, err = net.ListenMulticastUDP("udp", nil, addr)
        log.Println("Registered new reader2 for", addrinfo, r)
    }
}

func (p *Proxy) RegisterWriter(addrinfo string, c chan []byte){
    _, ok := p.writers[addrinfo]; if !ok{
        p.writers[addrinfo] = make([]chan []byte, 0)
    }

    p.writers[addrinfo] = append(p.writers[addrinfo], c)
    log.Println("Registered new writer for", addrinfo, c)

    log.Println("Readers", p.readers)
    log.Println("Writers", p.writers)
}

func (p *Proxy) RemoveWriter(addrinfo string, c chan[]byte){
    for i, wc := range p.writers[addrinfo]{
        if c == wc{
            p.writers[addrinfo] = append(p.writers[addrinfo][:i], p.writers[addrinfo][i+1:]...)
            log.Println("Removed writer", c)
        }
    }
    if len(p.writers[addrinfo]) == 0{
        log.Println("No more writers left for", addrinfo, "closing reader")
        p.readers[addrinfo].(io.ReadCloser).Close()
        delete(p.readers, addrinfo)
        delete(p.writers, addrinfo)
        log.Println("Removed all read/write references to", addrinfo)
    }

    log.Println("Readers", p.readers)
    log.Println("Writers", p.writers)
}

func (p *Proxy) Loop(){
    for{
        time.Sleep(100 * time.Millisecond)
        for{
            if len(p.readers) == 0{ break }
            for addrinfo, r := range p.readers{
                rtp := ReadRTP(r)

                if rtp == nil{ continue }

                if p.rtpseq[addrinfo] + 2 != rtp.SequenceNumber - 0xffff{
                    log.Println("RTP SEQUENCE FAILED FOR", addrinfo, "expecting", p.rtpseq[addrinfo] + 2, "have", rtp.SequenceNumber, rtp.SequenceNumber - 0xffff)
                }

                for _, c := range p.writers[addrinfo]{
                    c <- rtp.Payload
                }
                p.rtpseq[addrinfo] = rtp.SequenceNumber
            }
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
