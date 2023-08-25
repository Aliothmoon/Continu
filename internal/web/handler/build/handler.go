package build

import (
	"github.com/Aliothmoon/Continu/internal/repo/model"
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"log"
	"os"
	"path"
)

var (
	DLog = query.Log
)

const (
	user = "git"
)

func publishTask(buildID int32, p *model.Project) {
	var dir string
	if p.WorkDir != nil {
		dir = *p.WorkDir
	}
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		createLog(buildID, "Working directory creation is not available")
		return
	}
	ets, err := os.ReadDir(dir)
	if err != nil {
		createLog(buildID, "Working directory Read DirEntry error")
		return
	}
	for i := range ets {
		e := path.Join(dir, ets[i].Name())
		err := os.RemoveAll(e)
		if err != nil {
			createLog(buildID, "Working directory delete error")
			return
		}
	}

	return
}

func processWorkDir() {

}

func processGit(dir string, p *model.Project) bool {
	var (
		refer string
		pem   []byte
		url   string
	)
	if p.PrivateKey != nil {
		pem = []byte(*p.PrivateKey)
	}
	if p.ProjectURL != nil {
		url = *p.ProjectURL
	}
	if url == "" {
		log.Println("Error of git clone")
		return false
	}

	auth, err := ssh.NewPublicKeys(user, pem, "")
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:           url,
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName(refer),
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func createLog(buildID int32, lg string) {
	err := DLog.Create(&model.Log{
		BuildID: &buildID,
		Content: &lg,
	})
	log.Println(err)
}

func failBuildRecord(buildID int32) {
	info, err := DRecord.Where(DRecord.ID.Eq(buildID)).Update(DRecord.Status, biz.BuildFailed)
	if err != nil {

	}
	if info.RowsAffected != 1 {

	}
}

func successBuildRecord(buildID int32) {
	info, err := DRecord.Where(DRecord.ID.Eq(buildID)).Update(DRecord.Status, biz.BuildSuccess)
	if err != nil {

	}
	if info.RowsAffected != 1 {

	}
}
