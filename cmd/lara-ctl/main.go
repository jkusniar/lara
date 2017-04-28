package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jkusniar/lara/cmd"
	"github.com/jkusniar/lara/version"
)

var (
	printVersion = flag.Bool("v", false, "print version and exit")
	dbUser       = flag.String("dbUser", "postgres", "database user [env LARA_DB_USER]")
	dbPass       = flag.String("dbPass", "postgres", "database password [env LARA_DB_PASS]")
	dbHost       = flag.String("dbHost", "localhost", "database host [env LARA_DB_HOST]")
	dbPort       = flag.Uint("dbPort", uint(5432), "database port [env LARA_DB_PORT]")
	dbName       = flag.String("dbName", "lara", "database name [env LARA_DB_NAME]")
	dbSSLMode    = flag.String("dbSSLMode", "disable", "database connection SSL Mode [env LARA_DB_SSL_MODE]")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init flags
	flag.Usage = usage
	envVars()
	flag.Parse()

	if *printVersion {
		fmt.Println(version.String())
		return
	}

	cmd.CheckPortNum(*dbPort, "dbPort")

	if flag.NArg() == 0 {
		flag.Usage()
	}

	var err error
	switch flag.Arg(0) {
	case "register":
		err = register(*dbUser, *dbPass, *dbHost, *dbName, *dbPort, *dbSSLMode,
			flag.Args())
	case "grant":
		err = grant(*dbUser, *dbPass, *dbHost, *dbName, *dbPort, *dbSSLMode,
			flag.Args())
	case "revoke":
		err = revoke(*dbUser, *dbPass, *dbHost, *dbName, *dbPort, *dbSSLMode,
			flag.Args())
	default:
		flag.Usage()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: lara-ctl [flags] command [command arguments...]")
	fmt.Fprintln(os.Stderr, "flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "\tregister - register new user. Arguments: login password")
	fmt.Fprintln(os.Stderr, "\tgrant - grant permissions to user. Arguments: login permission1,permission2,...")
	fmt.Fprintln(os.Stderr, "\trevoke - revoke permissions from user. Arguments: login permission1,permission2,...")
	os.Exit(2)
}

func envVars() {
	cmd.StringVar(dbUser, "LARA_DB_USER")
	cmd.StringVar(dbPass, "LARA_DB_PASS")
	cmd.StringVar(dbHost, "LARA_DB_HOST")
	cmd.UintVar(dbPort, "LARA_DB_PORT")
	cmd.StringVar(dbName, "LARA_DB_NAME")
	cmd.StringVar(dbSSLMode, "LARA_DB_SSL_MODE")
}
