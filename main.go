package main

import (
    //"github.com/Comcast/gots"
    //"github.com/Comcast/gots/packet"
	"log"
	//"fmt"
	"os"
    "time"
    "flag"
    "strings"
	"net/http"
    "strconv"
    "encoding/base64"
)

var (
    proxy *Proxy
)

func RTPToHTTP(w http.ResponseWriter, req *http.Request){

    log.Println(req)

    if req.Method != "GET"{
        return
    }

    addrinfo := strings.Split(req.URL.Path, "/")[2]
    proxy.RegisterReader2(addrinfo)

    c := make(chan []byte, 1024)
    proxy.RegisterWriter(addrinfo, c)

    done := false
    closed := w.(http.CloseNotifier).CloseNotify()

    for{
        //log.Println("SELECT!")
        select{
        case b := <-c:
            n, err := w.Write(b); if err != nil{
                log.Println(err, n)
                done = true
                break
            }
            //log.Println(n, err)
        case <- closed:
            done = true
            break
        }
        if done{ break }
    }

    proxy.RemoveWriter(addrinfo, c)
}

func HTTPRec(w http.ResponseWriter, req *http.Request){
    log.Println(req)

    args := strings.Split(req.URL.Path, "/")
    addrinfo, _ := base64.URLEncoding.DecodeString(args[2])
    name, _ := base64.URLEncoding.DecodeString(args[3])
    season, _ := base64.URLEncoding.DecodeString(args[4]) // '' for any
    episode, _ := base64.URLEncoding.DecodeString(args[5]) // '' for any
    title, _ := base64.URLEncoding.DecodeString(args[6]) // '' for any

    var start       time.Time
    var duration    time.Duration
    var rec         interface{}

    if len(args) > 8{
        startint, err := strconv.ParseInt(args[7], 10, 64); if err != nil{
            log.Println(err)
        }
        start = time.Unix(startint, 0)
        duration, err = time.ParseDuration(args[8]); if err != nil{
            log.Println(err)
        }
    }

    rec = RecordingRequest{
        string(addrinfo),
        string(name),
        string(season),
        string(episode),
        string(title),
    }

    if ! start.IsZero() && duration.Seconds() > 0{
        rec = OnetimeRecordingRequest{start,duration,rec.(RecordingRequest)}
    }

    log.Printf("%+v", rec)

    go rec.(Recordable).Run()

}

func RTPToFile(addrinfo string){
    proxy.RegisterReader2(addrinfo)

	w, err := os.OpenFile(addrinfo, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664); if err != nil{
        log.Fatal(err)
    }

    c := make(chan []byte, 1024)
    proxy.RegisterWriter(addrinfo, c)

    done := false

    for{
        select{
        case b := <-c:
            n, err := w.Write(b); if err != nil{
                log.Println(err, n)
                done = true
                break
            }
        }
        if done{ break }
    }
    proxy.RemoveWriter(addrinfo, c)
}

func main(){

	log.SetFlags(log.LstdFlags | log.Lshortfile)

    opt_http := flag.Bool("http", false, "Enable HTTP Proxy")
    opt_testrtp := flag.String("testrtp", "", "Test a multicast RTP stream")

    flag.Parse()

    proxy = NewProxy()
    go proxy.Loop()

    go Recorder()

    if *opt_testrtp != ""{
        RTPToFile(*opt_testrtp)
    }

    if *opt_http{
	    NewHTTPServer(":1235", map[string]func(http.ResponseWriter, *http.Request){
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
