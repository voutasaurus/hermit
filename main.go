package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/voutasaurus/env"
)

func main() {
	s := &server{
		log:  log.New(os.Stderr, "hermit: ", log.LstdFlags|log.Llongfile|log.LUTC),
		addr: env.Get("HERMIT_ADDR").WithDefault(":8080"),
	}

	s.log.Println("starting...")

	m := http.NewServeMux()
	m.HandleFunc("/", s.defaultHandler)
	m.HandleFunc("/cmd", s.cmdHandler)

	s.log.Println("serving on", s.addr)
	s.log.Fatal(http.ListenAndServe(s.addr, m))
}

type server struct {
	log  *log.Logger
	addr string
}

func (s *server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Hermit: endpoint not found", 404)
}

func (s *server) cmdHandler(w http.ResponseWriter, r *http.Request) {
	var v struct {
		Cmd string `json:"cmd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "Hermit: bad json: "+err.Error(), 400)
	}
	if err := s.run(r.Context(), v.Cmd); err != nil {
		http.Error(w, "Hermit: bad run: "+err.Error(), 500)
	}
}

func (s *server) run(ctx context.Context, command string) error {
	s.log.Println("received command:", command)
	cc := strings.Split(command, " ")
	cmd := exec.CommandContext(ctx, cc[0], cc[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
