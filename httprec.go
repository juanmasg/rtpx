package main

import (
    "os"
    "fmt"
    "log"
    "io"
    "time"
    "strconv"
    "net/http"
    "strings"
    "encoding/base64"
)

func httprecCallback(orr *OnetimeRecordingRequest){

    recdir := "rec"
    os.MkdirAll(recdir, 0755)

    filepath := fmt.Sprintf("%s/%s %s.%s - %s.ts", recdir, orr.Name, orr.Season, orr.Episode, orr.Title)

    w := NewTimedWriter(filepath, orr.Duration)

    proxy.RegisterReader2(orr.Addrinfo)

    c := make(chan []byte, 1024)
    proxy.RegisterWriter(orr.Addrinfo, c)

    done := false
    for{
        select{
        case b, stillopen := <-c:

            if ! stillopen{
                done = true; break
            }

            n, err := w.Write(b); if err != nil{
                log.Println(err, n)
                if err == io.EOF{
                    log.Println("Done recording", filepath, orr)
                    done = true; break
                }
            }
        }
        if done{ break }
    }
    proxy.RemoveWriter(orr.Addrinfo, c)
    RemoveSchedule(orr)
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
    }else{
        log.Println("Periodic request received but EPG is empty")
    }

    log.Printf("%+v", rec)

    go rec.(Recordable).Run(true)
}
