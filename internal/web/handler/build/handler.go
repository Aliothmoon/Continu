package build

import (
	"errors"
	"github.com/Aliothmoon/Continu/internal/logger"
	"github.com/Aliothmoon/Continu/internal/repo/model"
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/bytedance/sonic"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"log"
	"os"
	"os/exec"
	"path"
)

var (
	DLog = query.Log
)

const (
	user = "git"
)

type ConstructInfo struct {
	BuildID int32
	Log     *LogWriteCloser
	Project *model.Project
}

func PublishTask(c *ConstructInfo) {
	err := doInternal(c)
	if err != nil {
		TxFailBuild(c)
	} else {
		TxSuccessBuild(c)
	}
}

func doInternal(c *ConstructInfo) (err error) {
	logger.Info("Publish Task")
	err = doProcessWorkDir(c)
	if err != nil {
		logger.Errorf("doProcessWorkDir Err %v", err)
		return
	} else {
		logger.Info("doProcessWorkDir Complete")
	}
	err = doProcessGit(c)
	if err != nil {
		logger.Errorf("doProcessGit Err %v", err)
		return
	} else {
		logger.Info("doProcessGit Complete")
	}
	err = doProcessExec(c)
	if err != nil {
		logger.Warnf("doProcessExec Err %v", err)
		return
	} else {
		logger.Info("doProcessExec Complete")
	}
	return nil
}

func doProcessWorkDir(c *ConstructInfo) error {
	var (
		dir string
		p   = c.Project
		//bid = c.BuildID
	)

	if p.WorkDir != nil {
		dir = *p.WorkDir
	}

	err := os.MkdirAll(dir, 0666)
	if err != nil {
		//createLog(bid, "Working directory creation is not available")
		logger.Error(err)
		return err
	}
	ets, err := os.ReadDir(dir)
	if err != nil {
		//createLog(bid, "Working directory Read DirEntry error")
		logger.Error(err)
		return err
	}
	// Todo Adaptation Linux rm -r & Windows rd /r
	for i := range ets {
		e := path.Join(dir, ets[i].Name())
		err = os.RemoveAll(e)
		if err != nil {
			//createLog(bid, "Working directory delete error")
			logger.Error(err)
			return err
		}
	}
	return nil
}

func doProcessGit(c *ConstructInfo) (err error) {
	var (
		refer = "main"
		pem   []byte
		url   string
		p     = c.Project
		dir   string
	)
	if p.WorkDir != nil {
		dir = *p.WorkDir
	}

	if p.Branch != nil {
		refer = *p.Branch
	}

	if p.PrivateKey != nil {
		pem = []byte(*p.PrivateKey)
	}
	if p.ProjectURL != nil {
		url = *p.ProjectURL
	}
	if url == "" {
		log.Println("Error of git clone")
		err = errors.New("Error of git clone ")
		return
	}

	auth, err := ssh.NewPublicKeys(user, pem, "")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:           url,
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName(refer),
		Progress:      c.Log,
	})
	if err != nil {
		log.Println(err)
		return
	}
	return nil
}

func doProcessExec(c *ConstructInfo) (err error) {
	var (
		bin  string
		para biz.Parameters
		dir  string
	)

	if c.Project.Bin != nil {
		bin = *c.Project.Bin
	}
	if c.Project.Parameters != nil {
		err := sonic.Unmarshal([]byte((*c.Project.Parameters)), &para)
		if err != nil {
			return err
		}
	}
	if c.Project.WorkDir != nil {
		dir = *c.Project.WorkDir
	}

	if bin == "" {
		err = errors.New("Command line error ")
		return
	}

	cmd := exec.Command(bin, para...)
	cmd.Stdout = c.Log
	cmd.Stderr = c.Log
	cmd.Dir = dir
	err = cmd.Start()

	ProcessMap.Store(c.BuildID, cmd.Process)

	err = cmd.Wait()

	ProcessMap.Delete(c.BuildID)

	_ = c.Log.Close()
	return err
}

func createLog(buildID int32, lg string) {

	err := DLog.Create(&model.Log{
		BuildID: &buildID,
		Content: &lg,
	})
	if err != nil {
		logger.Warnf("Record Log An error occurred %v", err)
	}
}

func TxSuccessBuild(c *ConstructInfo) {
	var err error
	_, err = DProject.Where(DProject.ID.Eq(c.Project.ID)).Update(DProject.Status, biz.ProjectIdle)
	if err != nil {
		logger.Errorf("Update Project Status An error occurred %v", err)
	}
	_, err = DRecord.Where(DRecord.ID.Eq(c.BuildID)).Update(DRecord.Status, biz.BuildSuccess)
	if err != nil {
		logger.Errorf("Update Build Record Status An error occurred %v", err)
	}

}

func TxFailBuild(c *ConstructInfo) {
	var err error
	_, err = DProject.Where(DProject.ID.Eq(c.Project.ID)).Update(DProject.Status, biz.ProjectIdle)
	if err != nil {
		logger.Errorf("Update Project Status An error occurred %v", err)
	}
	_, err = DRecord.Where(DRecord.ID.Eq(c.BuildID)).Update(DRecord.Status, biz.BuildFailed)
	if err != nil {
		logger.Errorf("Update Build Record Status An error occurred %v", err)
	}

}
