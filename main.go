package main

import (
    //"github.com/Comcast/gots"
    //"github.com/Comcast/gots/packet"
	"log"
	//"fmt"
    "flag"
	"net/http"
    "golang.org/x/net/webdav"
    "./rtpx"
    "./rtpx/rec"
    "./rtpx/routers"
    _ "net/http/pprof"
)

var (
    proxy *rtpx.Proxy
)


func main(){

	log.SetFlags(log.LstdFlags | log.Lshortfile)

    go http.ListenAndServe("localhost:6060", nil)

    opt_http := flag.String("http", "", "Enable HTTP Proxy to listen in this port")
    opt_testrtp := flag.String("testrtp", "", "Test a multicast RTP stream")

    flag.Parse()

    proxy = rtpx.NewProxy()
    go proxy.Loop()

    routers.SetProxy(proxy)

    go rec.Recorder(routers.HttprecCallback)

    if *opt_testrtp != ""{
        routers.RTPToFile(*opt_testrtp)
    }

    if *opt_http != ""{

        dav := webdav.Handler{Prefix: "/dav/"}
        dav.Logger = func(r *http.Request, err error){
            log.Println(err, r)
        }
        dav.FileSystem = webdav.Dir("rec/")
        dav.LockSystem = webdav.NewMemLS()

	    routers.NewHTTPServer(*opt_http, map[string]func(http.ResponseWriter, *http.Request){
	        "/udp/": routers.RTPToHTTP,
	        "/rtp/": routers.RTPToHTTP,
	        "/rec/": routers.HTTPRec,
            "/dav/": dav.ServeHTTP,
            "/rtpdav/": routers.RTPToHTTPViaDAV,
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
