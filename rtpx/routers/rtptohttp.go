package routers

import (
    "net/http"
    "strings"
    "log"
    "time"
)

func RTPToHTTP(w http.ResponseWriter, req *http.Request){

    log.Println(req)

    if req.Method != "GET"{
        return
    }

    addrinfo := strings.Split(req.URL.Path, "/")[2]
    proxy.RegisterReader(addrinfo)

    c := make(chan []byte, 32)
    //log.Println("CHANNEL_32", c)

    proxy.RegisterWriter(addrinfo, c)

    t := time.Now()

    done := false
    closed := w.(http.CloseNotifier).CloseNotify()

    for{
        //log.Println("SELECT!")
        select{
        case b, ok := <-c:
            if !ok{
                done = true
                break
            }
            n, err := w.Write(b); if err != nil{
                log.Println(err, n)
                done = true
                break
            }
            //log.Println(n, err)
        case <- closed:
            done = true
            break
        //case <- time.After(100 * time.Millisecond):
            //timeout
        }

        if done{ break }
    }

    proxy.RemoveWriter(addrinfo, c)

    log.Println("RTPtoHTTP session for", addrinfo, "lasted", time.Now().Sub(t))
}

