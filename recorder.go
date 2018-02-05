package main

import (
    "time"
    "log"
    "io"
    "os"
    "io/ioutil"
    "encoding/json"
)

var (
    recC                chan interface{}
    stopC               chan interface{}
    Requests            []*RecordingRequest
    OnetimeRequests     []*OnetimeRecordingRequest
)

func SaveSchedule(){
    log.Println("SAVE!", Requests, OnetimeRequests)
    jr, err := json.Marshal(&Requests); if err != nil{
        log.Println(err)
    }

    jor, err := json.Marshal(&OnetimeRequests); if err != nil{
        log.Println(err)
    }

    ioutil.WriteFile("rr.json", jr, 0644)
    ioutil.WriteFile("jor.json", jor, 0644)
}

func LoadSchedule() (error){
    jr, err := ioutil.ReadFile("rr.json"); if err != nil{
        return err
    }
    err = json.Unmarshal(jr, &Requests); if err != nil{
        return err
    }

    jor, err := ioutil.ReadFile("jor.json"); if err != nil{
        return err
    }
    err = json.Unmarshal(jor, &OnetimeRequests); if err != nil{
        return err
    }

    log.Println("Loaded schedule")
    log.Println("Requests", Requests)
    log.Println("OnetimeRequests", OnetimeRequests)

    RestartAll()

    return nil
}

func RestartAll(){
    for _, r := range Requests{
        go r.Run(false)
    }
    for _, r := range OnetimeRequests{
        go r.Run(false)
    }
}

func RemoveSchedule(orr *OnetimeRecordingRequest){
    for i, r := range OnetimeRequests{
        if r.Addrinfo == orr.Addrinfo && r.Start == orr.Start && r.Duration == orr.Duration{
            OnetimeRequests = append(OnetimeRequests[:i], OnetimeRequests[i+1:]...)
            log.Println("Removed request from schedule.", orr)
            SaveSchedule()
            break
        }
    }
}

type Recordable interface{
    Run(save bool)
}

type RecordingRequest struct{
    Addrinfo        string
    Name            string
    Season          string
    Episode         string
    Title           string
}

func (rr RecordingRequest) Run(save bool){

    if save{
        Requests = append(Requests, &rr)
        SaveSchedule()
    }

    // Find recurrent titles in epg
    // Schedule multiple one-time requests
}

type OnetimeRecordingRequest struct{
    Start           time.Time
    Duration        time.Duration
    RecordingRequest
}

func (orr *OnetimeRecordingRequest) MarshalJSON() ([]byte, error){
//    log.Println("MARSHAL", orr, orr.RecordingRequest)
    return json.Marshal(&struct{
        Start               int64
        Duration            float64
        RecordingRequest
    }{
        orr.Start.Unix(),
        orr.Duration.Seconds(),
        orr.RecordingRequest,
    })
}

func (orr *OnetimeRecordingRequest) UnmarshalJSON(data []byte) (error){
    aux := &struct{
        Start               int64
        Duration            float64
        RecordingRequest
    }{}
    err := json.Unmarshal(data, aux); if err != nil{
        log.Println(err)
        return err
    }

    orr.Start = time.Unix(aux.Start, 0)
    orr.Duration = time.Duration(aux.Duration) * 1000000000
    orr.RecordingRequest = aux.RecordingRequest

    log.Println(orr)

    return nil
}

func (orr *OnetimeRecordingRequest) Enabled() bool{
    return orr.Start.IsZero() || orr.Duration.Seconds() == 0
}

func (orr OnetimeRecordingRequest) Run(save bool){

    if save{
        OnetimeRequests = append(OnetimeRequests, &orr)
        SaveSchedule()
    }

    now := time.Now()
    then := orr.Start.Sub(now)

    log.Printf("Recording scheduled: %s %s-%s: %s from %s. Start on %s time left %s",
        orr.Name,
        orr.Season,
        orr.Episode,
        orr.Title,
        orr.Addrinfo,
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

    err := LoadSchedule(); if err != nil{
        log.Println(err)
        Requests = make([]*RecordingRequest, 0)
        OnetimeRequests = make([]*OnetimeRecordingRequest, 0)
    }

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
