package main

import (
    "time"
    "log"
    "io"
    "os"
)

var (
    //recF    func(
    recC    chan interface{}
    stopC   chan interface{}
)

type Recordable interface{
    Run()
}

type RecordingRequest struct{
    addrinfo        string
    name            string
    season          string
    episode         string
    title           string
}
func (rr RecordingRequest) Run(){
}

type OnetimeRecordingRequest struct{
    Start           time.Time
    Duration        time.Duration
    RecordingRequest
}

func (orr *OnetimeRecordingRequest) Enabled() bool{
    return orr.Start.IsZero() || orr.Duration.Seconds() == 0
}

func (orr OnetimeRecordingRequest) Run(){
    now := time.Now()
    then := orr.Start.Sub(now)
    log.Printf("Recording scheduled: %s %s-%s: %s from %s. Start on %s time left %s",
        orr.name,
        orr.season,
        orr.episode,
        orr.title,
        orr.addrinfo,
        orr.Start,
        then,
    )
    <-time.NewTimer(then).C
    log.Println("Request", recC, orr)
    recC <- &orr
}

func Recorder(callback func(*OnetimeRecordingRequest)){
    log.Println("Recorder start...")
    recC = make(chan interface{})
    for{
        rr := <-recC
        log.Println("Recording starts", rr)
        orr := rr.(*OnetimeRecordingRequest); if orr != nil{
            log.Printf("ORR %+v\n", orr)
            callback(orr)
        }
    }
}

type TimedWriter struct{
    t           *time.Timer
    realw       io.Writer
    path        string
    duration    time.Duration
}

func NewTimedWriter(path string, duration time.Duration) *TimedWriter{
    var err error

    w := &TimedWriter{}
    w.path = path
	w.realw, err = os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664); if err != nil{
        log.Println("ERROR opening file for writing", err)
    }
    w.t = time.NewTimer(duration)
    log.Println("Recording timer started", path, duration)

    return w
}

func (w *TimedWriter) Write(b []byte) (n int, err error){
    select {
        case <- w.t.C:
            w.realw.(io.Closer).Close()
            log.Println("Finished recording", w.path)
            return 0, io.EOF
        default:
            return w.realw.Write(b)
    }
}
