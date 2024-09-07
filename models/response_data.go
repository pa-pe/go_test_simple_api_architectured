package models

type ProcessingInfo struct {
    TimeTaken         string `json:"TimeTaken"`
    DuplicatesRemoved int    `json:"DuplicatesRemoved"`
}

type ResponseData struct {
    Name          string         `json:"Name"`
    Last          string         `json:"Last"`
    Addresses     []Address      `json:"Addresses"`
    ProcessingInfo ProcessingInfo `json:"ProcessingInfo"`
}
