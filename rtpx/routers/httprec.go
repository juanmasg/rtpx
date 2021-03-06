package routers

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
    "../rec"
)

func HttprecCallback(orr *rec.OnetimeRecordingRequest){

    recdir := "rec"
    seriedir := fmt.Sprintf("%s/%s", recdir, orr.Name)
    os.MkdirAll(seriedir, 0755)

    filepath := fmt.Sprintf("%s/%s S%02sE%02s", seriedir, orr.Name, orr.Season, orr.Episode)
    if orr.Title != ""{
        filepath = fmt.Sprintf("%s - %s.ts", filepath, orr.Title)
    }else{
        filepath += ".ts"
    }

    w := rec.NewTimedWriter(filepath, orr.Duration)

    proxy.RegisterReader(orr.Addrinfo)

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
    rec.RemoveSchedule(orr)
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
    var recreq      interface{}

    if len(args) > 8{
        startint, err := strconv.ParseInt(args[7], 10, 64); if err != nil{
            log.Println(err)
        }
        start = time.Unix(startint, 0)
        duration, err = time.ParseDuration(args[8]); if err != nil{
            log.Println(err)
        }
    }

    recreq = rec.RecordingRequest{
        string(addrinfo),
        string(name),
        string(season),
        string(episode),
        string(title),
    }

    if ! start.IsZero() && duration.Seconds() > 0{
        recreq = rec.OnetimeRecordingRequest{start,duration,recreq.(rec.RecordingRequest)}
    }else{
        log.Println("Periodic request received but EPG is empty")
    }

    log.Printf("%+v", recreq)

    go recreq.(rec.Recordable).Run(true)
}
