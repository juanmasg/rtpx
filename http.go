package main

import (
    "net/http"
    "log"
    //"io"
//    "bufio"
)

/*
func handleUDP(w http.ResponseWriter, req *http.Request){
    log.Printf("%+v", req)
    path := strings.Split(req.URL.Path, "/")
    addrinfo := strings.Split(path[2], ":")

    if len(addrinfo) != 2{
        log.Println("BAD ADDRESS", addrinfo)
        return
    }

    addr := addrinfo[0]
    port, err := strconv.Atoi(addrinfo[1]); if err != nil{
        log.Println(addrinfo, err)
        return
    }

    log.Println(addr, port)

    r := NewMulticastReader(addr, port)

    c := make(chan []byte, 1024)
    go func(){
        for{
            b := make([]byte, 1500)
            r.Read(b)
            c <- b
        }
    }()

    for{
        b := <-c
		rtp := RTPPacket(b)
        w.Write(rtp.Payload)
        //log.Println(len(c))
    }


    //for{
    //    b := make([]byte, 1500)
    //    n, err := r.Read(b); if err != nil{
    //        log.Println(err, n)
    //        break
    //    }
    //    log.Println(n)
    //    w.Write(b)
    //}
    //bufr := bufio.NewReaderSize(r, 2 * 1024 * 1024)
    //io.Copy(w, bufr)
    //io.Copy(w, r)
    //defer r.(io.Closer).Close()
}

func NewHTTPServer(bindstr string){
    http.HandleFunc("/udp/", handleUDP)
    err := http.ListenAndServe(bindstr, nil); if err != nil{
        log.Fatal(err)
    }
}
*/

func NewHTTPServer(bindstr string, maps map[string]func(http.ResponseWriter, *http.Request)){
	for prefix, f := range maps{
		http.HandleFunc(prefix, f)
	}
    err := http.ListenAndServe(bindstr, nil); if err != nil{
        log.Fatal(err)
    }
}
