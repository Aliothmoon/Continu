package build

import (
	"errors"
	"fmt"
	"github.com/Aliothmoon/Continu/internal/logger"
	"github.com/Aliothmoon/Continu/internal/repo/model"
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/bytedance/sonic"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"os/exec"
	"path"
	"time"
)

var (
	DLog = query.Log
)

const (
	user = "git"
)

type ConstructInfo struct {
	BuildID     int32
	Log         *LogWriter
	ProjectInfo *model.Project
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
	isGit := false
	if c.ProjectInfo.IsGit != nil {
		isGit = *c.ProjectInfo.IsGit != 0
	}
	if isGit {
		err = doProcessWorkDir(c)

		if err != nil {
			logger.Debugf("doProcessWorkDir Err %v", err)
			tips := "Exception occurred in emptying the working directory"
			createLog(c.BuildID, fmt.Sprintf("%v\nProcessWorkDir Err %v", tips, err))
			return
		} else {
			createLog(c.BuildID, "doProcessWorkDir Complete")
		}

		err = doProcessGit(c)
		if err != nil {
			logger.Debugf("doProcessGit Err %v", err)
			createLog(c.BuildID, fmt.Sprintf("ProcessGit Err %v", err))
			return
		} else {
			createLog(c.BuildID, "doProcessGit Complete")
		}
	}

	err = doProcessExec(c)
	if err != nil {
		logger.Debugf("doProcessExec Err %v", err)
		createLog(c.BuildID, fmt.Sprintf("ProcessExec Err %v", err))
		return
	} else {
		createLog(c.BuildID, "doProcessExec Complete")
	}
	return nil
}

func doProcessWorkDir(c *ConstructInfo) error {
	var (
		dir string
		p   = c.ProjectInfo
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

func doProcessGit(c *ConstructInfo) error {
	var (
		refer = "main"
		pem   []byte
		url   string
		p     = c.ProjectInfo
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
		return errors.New("Error of git clone & URL is Null ")
	}

	auth, err := ssh.NewPublicKeys(user, pem, "")
	if err != nil {
		return err
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:           url,
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName(refer),
	})
	if err != nil {
		return err
	}
	return nil
}

func doProcessExec(c *ConstructInfo) (err error) {
	var (
		bin  string
		para biz.Parameters
		dir  string
	)

	if c.ProjectInfo.Bin != nil {
		bin = *c.ProjectInfo.Bin
	}
	if c.ProjectInfo.Parameters != nil {
		err := sonic.Unmarshal([]byte((*c.ProjectInfo.Parameters)), &para)
		if err != nil {
			return err
		}
	}
	if c.ProjectInfo.WorkDir != nil {
		dir = *c.ProjectInfo.WorkDir
	}

	if bin == "" {
		return errors.New("Command line error ")
	}

	cmd := exec.Command(bin, para...)

	cmd.Stdout = c.Log
	cmd.Stderr = c.Log
	cmd.Dir = dir
	Config(cmd)
	err = cmd.Start()

	if err != nil {
		return err
	}
	ProcessMap.Store(c.BuildID, cmd.Process)

	err = cmd.Wait()

	ProcessMap.Delete(c.BuildID)

	return err
}

func createLog(buildID int32, lg string) {
	if lg == "" {
		return
	}
	err := DLog.Create(&model.Log{
		BuildID:   buildID,
		Content:   lg,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		logger.Warnf("Record Log An error occurred %v", err)
	}
}

func TxSuccessBuild(c *ConstructInfo) {
	var err error
	_, err = DProject.Where(DProject.ID.Eq(c.ProjectInfo.ID)).Update(DProject.Status, biz.ProjectIdle)
	if err != nil {
		logger.Errorf("Update ProjectInfo Status An error occurred %v", err)
	}
	_, err = DRecord.Where(DRecord.ID.Eq(c.BuildID)).Update(DRecord.Status, biz.BuildSuccess)
	if err != nil {
		logger.Errorf("Update Build Record Status An error occurred %v", err)
	}

}

func TxFailBuild(c *ConstructInfo) {
	var err error
	_, err = DProject.Where(DProject.ID.Eq(c.ProjectInfo.ID)).Update(DProject.Status, biz.ProjectIdle)
	if err != nil {
		logger.Errorf("Update ProjectInfo Status An error occurred %v", err)
	}
	_, err = DRecord.Where(DRecord.ID.Eq(c.BuildID)).Update(DRecord.Status, biz.BuildFailed)
	if err != nil {
		logger.Errorf("Update Build Record Status An error occurred %v", err)
	}

}
