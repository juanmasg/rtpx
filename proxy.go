package main

import (
    "io"
    "strings"
    "strconv"
    "time"
    "log"
    "net"
    "sync"
)

type Proxy struct{
    sync.RWMutex
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
    p.Lock()
    ip, port := addrinfoToIPPort(addrinfo)
	r, ok := p.readers[addrinfo]; if !ok{
		r = NewMulticastReader(ip, port)
		p.readers[addrinfo] = r
        log.Println("Registered new reader for", addrinfo, r)
	}
    p.Unlock()
}

func (p *Proxy) RegisterReader2(addrinfo string){
    p.Lock()
    r, ok := p.readers[addrinfo]; if !ok{
	    addr, err := net.ResolveUDPAddr("udp", addrinfo); if err != nil {
	        log.Fatal(err)
	    }
	    r, err = net.ListenMulticastUDP("udp", nil, addr)
        r.(*net.UDPConn).SetReadBuffer(4 * 1024 * 1024)
        p.readers[addrinfo] = r
        log.Println("Registered new reader2 for", addrinfo, r)
    }
    p.Unlock()
}

func (p *Proxy) RegisterWriter(addrinfo string, c chan []byte){
    p.Lock()
    _, ok := p.writers[addrinfo]; if !ok{
        p.writers[addrinfo] = make([]chan []byte, 0)
    }

    p.writers[addrinfo] = append(p.writers[addrinfo], c)
    log.Println("Registered new writer for", addrinfo, c)

    log.Println("Readers", p.readers)
    log.Println("Writers", p.writers)

    p.Unlock()
}

func (p *Proxy) RemoveWriter(addrinfo string, c chan[]byte){
    p.Lock()

    for i, wc := range p.writers[addrinfo]{
        if c == wc{
            p.writers[addrinfo] = append(p.writers[addrinfo][:i], p.writers[addrinfo][i+1:]...)
            log.Println("Removed writer", c)
        }
    }
    if len(p.writers[addrinfo]) == 0{
        log.Println("No more writers left for", addrinfo, "closing reader")
        delete(p.writers, addrinfo)

        _, exists := p.readers[addrinfo]; if exists{
            p.readers[addrinfo].(io.ReadCloser).Close()
            delete(p.readers, addrinfo)
        }
        log.Println("Removed all read/write references to", addrinfo)
    }

    log.Println("Readers", p.readers)
    log.Println("Writers", p.writers)

    p.Unlock()
}

func (p *Proxy) RemoveReader(reader io.Reader){
    log.Println("Remove reader",  reader, p.readers)
    for addrinfo, r := range p.readers{
        if r == reader{
            log.Println("Found addrinfo for reader", addrinfo, ". Remove all writers")
            r.(io.Closer).Close()
            delete(p.readers, addrinfo)
            for _, c := range p.writers[addrinfo]{
                close(c)
            }
        }
    }
}

func (p *Proxy) Loop(){
    for{
        time.Sleep(100 * time.Millisecond)
        for{
            if len(p.readers) == 0{ break }
            p.Lock()
            for addrinfo, r := range p.readers{
                r.(*net.UDPConn).SetReadDeadline(time.Now().Add(250 * time.Millisecond)) //FIXME: ?!
                rtp, err := ReadRTP(r); if err != nil{
                    if err == io.EOF{
                        log.Println("EOF from", r)
                    }else{
                        operr := err.(*net.OpError)
                        if operr != nil && operr.Timeout(){
                            log.Println("Timeout from", r)
                        }
                    }
                    p.RemoveReader(r)
                    break
                }

                if rtp == nil{ continue }

                if p.rtpseq[addrinfo] + 2 != rtp.SequenceNumber - 0xffff{
                    log.Println("RTP SEQUENCE FAILED FOR", addrinfo, "expecting", p.rtpseq[addrinfo] + 2, "have", rtp.SequenceNumber, rtp.SequenceNumber - 0xffff)
                }

                for _, c := range p.writers[addrinfo]{
                    c <- rtp.Payload
                }
                p.rtpseq[addrinfo] = rtp.SequenceNumber
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
    for{
        
    }
}


