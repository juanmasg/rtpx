package routers

import (
    "os"
    "log"
)

func RTPToFile(addrinfo string){
    proxy.RegisterReader(addrinfo)

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

