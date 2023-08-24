package model

type Project struct {
	Name       string `json:"Name"`
	Status     int32  `json:"Status"`
	Branch     string `json:"Branch"`
	ProjectURL string `json:"ProjectUrl"`
	PrivateKey string `json:"PrivateKey"`
	Script     string `json:"Script"`
}
