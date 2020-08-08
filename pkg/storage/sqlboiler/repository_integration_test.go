// +build integration

package sqlboiler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	var code int
	ctx := context.Background()
	dbContainer, err := testdata.NewDBContainer(ctx)
	if err != nil {
		log.Fatalln(err)
		return 1
	}
	defer func() {
	}()
	dbConnStr = dbContainer.ConnStr
	if err := config.InitDB(dbConnStr); err != nil {
		log.Fatalln("init db: ", err)
		return 1
	}
	defer func() {
		if err := (*dbContainer.Container).Terminate(ctx); err != nil {
			log.Fatalf("terminate dbContainer: %s\n", err)
		}
	}()
	if code > 0 {
		return code
	}
	code = m.Run()
	return code
}
