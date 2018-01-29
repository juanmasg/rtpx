package main

import (
    "net/http"
    "log"
)

func NewHTTPServer(bindstr string, maps map[string]func(http.ResponseWriter, *http.Request)){
	for prefix, f := range maps{
		http.HandleFunc(prefix, f)
	}
    err := http.ListenAndServe(bindstr, nil); if err != nil{
        log.Fatal(err)
    }
}
