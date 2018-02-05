package main

import (
    "net/http"
    "strings"
    "log"
)

func RTPToHTTP(w http.ResponseWriter, req *http.Request){

    log.Println(req)

    if req.Method != "GET"{
        return
    }

    addrinfo := strings.Split(req.URL.Path, "/")[2]
    proxy.RegisterReader(addrinfo)

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

