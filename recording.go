package main

import (
    "time"
    "log"
)

var (
    recF    func(
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

func Recorder(){
    log.Println("Recorder start...")
    recC = make(chan interface{})
    for{
        rr := <-recC
        log.Println("Recording starts", rr)
        orr := rr.(*OnetimeRecordingRequest); if orr != nil{
            log.Printf("ORR %+v\n", orr)
            continue
        }

        prr := rr.(*RecordingRequest); if prr != nil{
            log.Printf("PRR %+v\n", orr)
            continue
        }

    }
}

type TimedWriter struct{
}

func (w *TimedWriter) Write(b []byte) (n int, err error){
    select {
    }
}
