package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	h "github.com/rinosukmandityo/hexagonal-login/api"
	"github.com/rinosukmandityo/hexagonal-login/helper"
	"github.com/rinosukmandityo/hexagonal-login/logic"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	set mongo_url=mongodb://localhost:27017/local
	set mongo_timeout=30
	set mongo_db=local
	set url_db=mongo
*/

func main() {
	userRepo := helper.ChooseRepo()
	userService := logic.NewUserService(userRepo)
	handler := h.NewUserHandler(userService)

	r := h.RegisterHandler(handler)

	errs := make(chan error, 2)
	go func() {
		log.Printf("Listening on port %s\n", httpPort())
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)

	}()
	log.Printf("Terminated %s", <-errs)

}

func httpPort() string {
	port := "8000"
	if os.Getenv("port") != "" {
		port = os.Getenv("port")
	}
	return fmt.Sprintf(":%s", port)
}
