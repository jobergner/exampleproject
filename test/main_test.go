package test

import (
	"exampleproject/db"
	"exampleproject/server"
	"os"
	"path"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func init() {
	// NOTE: we pretend we're running from root directory
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	os.Setenv("POSTGRES_USER", "TEST_USER")
	os.Setenv("POSTGRES_PASSWORD", "TEST_PASSWORD")
	os.Setenv("POSTGRES_DB", "TEST_DB")
}

func TestMain(t *testing.T) {
	cancel, err := StartPostgres()
	if err != nil {
		panic(err)
	}
	defer cancel()

	if err := db.Connect(); err != nil {
		panic(err)
	}

	if err := db.MigrateUp(); err != nil {
		panic(err)
	}

	signalShutdown, signalSigInt := server.Listen()

	time.Sleep(time.Second * 3)

	signalSigInt <- syscall.SIGINT

	if err := <-signalShutdown; err != nil {
		panic(err)
	}
}
