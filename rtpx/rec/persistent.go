package rec

import (
    "log"
    "io/ioutil"
    "encoding/json"
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

