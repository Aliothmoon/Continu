package repo

import (
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"log"
	"testing"
)

func TestQuery(t *testing.T) {
	logs, err := query.Log.Find()
	if err != nil {
		t.Error(err)
	}
	for i := range logs {
		log.Printf("%v", logs[i])
	}

}
