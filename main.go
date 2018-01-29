package main

import (
    //"github.com/Comcast/gots"
    //"github.com/Comcast/gots/packet"
	"log"
	"io"
	//"fmt"
	"net"
	"os"
    "time"
    "flag"
    "strconv"
    "strings"
	"net/http"
)

var (
	readers map[string]io.Reader
    writers map[string]io.Writer
)

func AddrinfoToIPPort(addrinfo string) (ip string, port int){
    ipport := strings.Split(addrinfo, ":")
    if len(ipport) != 2{
        log.Println("BAD ADDRESS", addrinfo)
        return
    }

    ip = ipport[0]
    port, err := strconv.Atoi(ipport[1]); if err != nil{
        log.Println(ipport, err)
        return
    }

    log.Println(ip, port)
    return
}

func GetMulticastReader(addrinfo string) io.Reader{
    ip, port := AddrinfoToIPPort(addrinfo)
	udpr, ok := readers[addrinfo]; if !ok{
		udpr = NewMulticastReader(ip, port)
		readers[addrinfo] = udpr
	}

    return udpr
}

func GetMulticastReader2(addrinfo string) io.Reader{
	addr, err := net.ResolveUDPAddr("udp", addrinfo); if err != nil {
	    log.Fatal(err)
	}
	r, err := net.ListenMulticastUDP("udp", nil, addr)

    return r
}

func RTPToHTTP(w http.ResponseWriter, req *http.Request){
    addrinfo := strings.Split(req.URL.Path, "/")[2]
    udpr := GetMulticastReader(addrinfo)

    Copy(udpr, w)
}

func RTPToFile(addrinfo string){
    udpr := GetMulticastReader(addrinfo)

	f, err := os.OpenFile(addrinfo, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664); if err != nil{
        log.Fatal(err)
    }

    Copy(udpr, f)
}

func Copy(r io.Reader, w io.Writer){
    for{
        rtp := ReadRTP(r)
        _, err := w.Write(rtp.Payload); if err != nil{
            break
        }
    }
}

func Loop(){
    for{
        time.Sleep(100 * time.Millisecond)
    }
}

func main(){

	log.SetFlags(log.LstdFlags | log.Lshortfile)

    opt_http := flag.Bool("http", false, "Enable HTTP Proxy")
    opt_testrtp := flag.String("testrtp", "", "Test a multicast RTP stream")

    flag.Parse()

	readers = make(map[string]io.Reader)

    if *opt_testrtp != ""{
        RTPToFile(*opt_testrtp)
    }

    if *opt_http{
	    NewHTTPServer(":1234", map[string]func(http.ResponseWriter, *http.Request){
	        "/udp/": RTPToHTTP,
	        "/rtp/": RTPToHTTP,
	    })
    }
}

/*
func main1(){
    for read, err := r.Read(pkt); read > 0 && err == nil; read, err = r.Read(pkt) {
//		if err != nil {
//             println(err)
//             return
//		}
//        pid, err := packet.Pid(pkt)
//        if err != nil {
//            println(err)
//            continue
//        }
//        pidSet[pid] = true
//		//pat, err := packet.IsPat(pkt)
//		cc, err := packet.ContinuityCounter(pkt)
//		if ( prevcc + 1 ) % 16 != cc{
//			//log.Println("prev", prevcc, "cc", cc)
//			panic(cc)
//		}
//		prevcc = cc
		fmt.Print(string(pkt))
        //fmt.Println("Found pid", pid, pat, cc)
    }
}
*/
