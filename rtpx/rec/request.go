package rec

import (
    "time"
    "log"
    "encoding/json"
)

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
