// Copyright 2018 Mike Jarmy. All rights reserved.  Use of this
// source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	port := flag.String("port", "8080", "port on localhost")
	dir := flag.String("dir", "../../docs", "directory to server")

	log.Printf("Serving '" + *dir + "' on port " + *port)
	http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*dir)))

}
