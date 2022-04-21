package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	error := false
	defer func() {
		if !error {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("playing..."))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}()

	fmt.Println(*req)

	if req.ParseForm() != nil {
		error = true
		return
	}

	params, ok1 := req.Form["play_url"]
	params = append(params, "--force-window")
	args, ok2 := req.Form["mpv_args"]

	fmt.Println("MPV ARGS:", args)
	fmt.Println("URL:", params[0])
	fmt.Println("")

	if !ok1 || !ok2 {
		error = true
		return
	}

	params = append(params, args...)
	c := exec.Command("mpv", params...)
	c.Stdout = os.Stdout

	err := c.Start()

	if err != nil {
		error = true
		fmt.Println("Cannot launch mpv: ", err.Error())
	}
}

func main() {
	port := flag.String("port", "7531", "The port to listen on")
	public := flag.Bool("public", true, "Accept traffic from other computers")
	flag.Parse()
	hostname := "localhost"
	if *public {
		hostname = "0.0.0.0"
	}
	fmt.Println("Serving on: ", hostname+":"+*port)

	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(hostname+":"+(*port), nil)
}
