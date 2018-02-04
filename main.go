package main

import (
    //"github.com/Comcast/gots"
    //"github.com/Comcast/gots/packet"
	"log"
	//"fmt"
    "flag"
	"net/http"
)

var (
    proxy *Proxy
)


func main(){

	log.SetFlags(log.LstdFlags | log.Lshortfile)

    opt_http := flag.String("http", "", "Enable HTTP Proxy to listen in this port")
    opt_testrtp := flag.String("testrtp", "", "Test a multicast RTP stream")

    flag.Parse()

    proxy = NewProxy()
    go proxy.Loop()

    go Recorder(httprecCallback)

    if *opt_testrtp != ""{
        RTPToFile(*opt_testrtp)
    }

    if *opt_http != ""{
	    NewHTTPServer(*opt_http, map[string]func(http.ResponseWriter, *http.Request){
	        "/udp/": RTPToHTTP,
	        "/rtp/": RTPToHTTP,
	        "/rec/": HTTPRec,
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
