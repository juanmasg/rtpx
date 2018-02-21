package rec

import (
    "time"
    "log"
    "io"
    "os"
)

var (
    recC                chan interface{}
    stopC               chan interface{}
    Requests            []*RecordingRequest
    OnetimeRequests     []*OnetimeRecordingRequest
)

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
