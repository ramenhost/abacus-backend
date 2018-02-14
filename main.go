package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var errlog *log.Logger

func main() {
	a := App{}

	var dbhost, dbuser, dbpass, dbname string
	var port string

	flag.StringVar(&port, "port", "8080", "change default port")
	flag.Parse()

	dbhost = os.Getenv("DBHOST")
	dbuser = os.Getenv("DBUSER")
	dbpass = os.Getenv("DBPASS")
	dbname = os.Getenv("DBNAME")

	logfile, _ := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errlog = log.New(logfile, "Error:", log.Lshortfile|log.LstdFlags)

	a.Initialize(dbhost, dbuser, dbpass, dbname)

	a.Run(fmt.Sprintf(":%s", port))
}
