package routers

import (
    "net/http"
    "log"
    "../../rtpx"
)

var proxy *rtpx.Proxy
var autorec bool

func SetProxy(p *rtpx.Proxy){
    proxy = p
}

func SetAutorec(a bool){
    autorec = a
}

func NewHTTPServer(bindstr string, maps map[string]func(http.ResponseWriter, *http.Request)){
	for prefix, f := range maps{
		http.HandleFunc(prefix, f)
	}
    err := http.ListenAndServe(bindstr, nil); if err != nil{
        log.Fatal(err)
    }
}
