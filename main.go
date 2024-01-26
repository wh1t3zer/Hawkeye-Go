package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye-Go/router"
)

func main() {
	lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis", "micro"})
	defer lib.Destroy()

	router.HTTPServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HTTPServerStop()
}
