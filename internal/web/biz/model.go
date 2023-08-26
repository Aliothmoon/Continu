package biz

const (
	PID = "PID"
	RID = "RID"
)

// JsonModel ViewModel
type JsonModel struct {
	Code int    // O - Ok
	Msg  string `json:"Msg,,omitempty"`
	Data any    `json:"Data,,omitempty"`
}

type Parameters []string

type Project struct {
	ID         int32  `json:"Id"`
	Name       string `json:"Name"`
	Status     int32  `json:"Status"`
	Branch     string `json:"Branch"`
	ProjectURL string `json:"ProjectUrl"`
	PrivateKey string `json:"PrivateKey"`
	Bin        string `json:"Bin"`
	Parameters string `json:"Parameters"`
}
