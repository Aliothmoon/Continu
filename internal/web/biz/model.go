package biz

const (
	PID       = "PID"
	RID       = "RID"
	TimeStamp = "ts"
)

// JsonModel ViewModel
type JsonModel struct {
	Code int    `json:"code"` // O - Ok
	Msg  string `json:"msg,,omitempty"`
	Data any    `json:"data,,omitempty"`
}

type Parameters []string

type Project struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Status     int32  `json:"status"`
	Branch     string `json:"branch"`
	ProjectURL string `json:"projectUrl"`
	WorkDir    string `json:"workDir"`
	IsGit      bool   `json:"isGit"`
	PrivateKey string `json:"privateKey"`
	Bin        string `json:"bin"`
	Parameters string `json:"parameters"`
}
type Log struct {
	ID        int32   `json:"Id"`
	BuildID   *int32  `json:"BuildId"`
	Content   *string ` json:"Content"`
	CreatedAt int64   ` json:"CreatedAt"`
}
