package models

type Address struct {
    Country string `json:"Country"`
    City    string `json:"City"`
}

type RequestData struct {
    Name      string    `json:"Name"`
    Last      string    `json:"Last"`
    Addresses []Address `json:"Addresses"`
}
