/*
   Copyright (C) 2016-2017 Contributors as noted in the AUTHORS file

   This file is part of lara, veterinary practice support software.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jkusniar/lara/cmd"
	"github.com/jkusniar/lara/crypto"
	"github.com/jkusniar/lara/http"
	"github.com/jkusniar/lara/postgres"
	"github.com/jkusniar/lara/version"
)

var (
	printVersion = flag.Bool("v", false, "print version and exit")
	httpsPort    = flag.Uint("httpsPort", uint(8443), "https server port [env LARA_HTTPS_PORT]")
	hostname     = flag.String("hostname", "localhost", "hostname/addr to listen on [env LARA_HOSTNAME]")
	dbUser       = flag.String("dbUser", "postgres", "database user [env LARA_DB_USER]")
	dbPass       = flag.String("dbPass", "postgres", "database password [env LARA_DB_PASS]")
	dbHost       = flag.String("dbHost", "localhost", "database host [env LARA_DB_HOST]")
	dbPort       = flag.Uint("dbPort", uint(5432), "database port [env LARA_DB_PORT]")
	dbName       = flag.String("dbName", "lara", "database name [env LARA_DB_NAME]")
	dbSSLMode    = flag.String("dbSSLMode", "disable", "database connection SSL Mode [env LARA_DB_SSL_MODE]")
	rsaPriv      = flag.String("rsaPriv", "lara.rsa", "private RSA key [env LARA_RSA_PRIV]")
	rsaPub       = flag.String("rsaPub", "lara.rsa.pub", "public RSA key [env LARA_RSA_PUB]")
	tlsKey       = flag.String("tlsKey", "key.pem", "TLS private key [env LARA_TLS_KEY]")
	tlsCert      = flag.String("tlsCert", "cert.pem", "TLS certificate [env LARA_TLS_CERT]")
	wwwRoot      = flag.String("wwwRoot", "static", "Directory containing web client [env LARA_WWW_ROOT]")
)

/*
	NOTE ON TLS KEYS:

Generate self signed TLS keypair:
	$ cd $GOROOT/src/crypto/tls && go build generate_cert.go
	$ ./generate_cert --host localhost --ca
	$ mv *.pem /directory/for/keys/

dump:
	$ openssl x509 -in cert.pem -text
*/

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

	cmd.CheckPortNum(*httpsPort, "httpsPort")
	cmd.CheckPortNum(*dbPort, "dbPort")
	cmd.CheckFileExists(*rsaPriv)
	cmd.CheckFileExists(*rsaPub)
	cmd.CheckFileExists(*tlsKey)
	cmd.CheckFileExists(*tlsCert)
	cmd.CheckFileExists(*wwwRoot)

	// load encryption keys
	jwt, err := crypto.NewJWTToken(*rsaPriv, *rsaPub, *hostname)
	if err != nil {
		log.Fatalf("FATAL: inicializing JWT token provider failed: %+v\n", err)
	}

	// connect to DB
	db, err := postgres.Open(*dbUser, *dbPass, *dbHost, *dbName, *dbPort,
		*dbSSLMode)
	if err != nil {
		log.Fatalf("FATAL: connecting to database failed: %+v\n", err)
	}
	defer db.Close()

	// server
	sls := postgres.SimpleLovService{DB: db}
	srv := &http.Server{
		Token:          jwt,
		TitleService:   &sls,
		UnitService:    &sls,
		GenderService:  &sls,
		SpeciesService: &sls,
		BreedService:   &sls,
		AddressService: &postgres.AddressService{DB: db},
		SearchService:  &postgres.SearchService{DB: db},
		OwnerService:   &postgres.OwnerService{DB: db},
		PatientSevice:  &postgres.PatientService{DB: db},
		RecordService:  &postgres.RecordService{DB: db},
		ProductService: &postgres.ProductService{DB: db},
		ReportService:  &postgres.ReportService{DB: db, Loc: time.Local},
		UserService:    &postgres.UserService{DB: db, Pass: crypto.NewPassword()},
		TagService:     &postgres.TagService{DB: db},
		WWWRoot:        *wwwRoot,
	}

	// shutdown signal handler
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-quit
		if err := srv.Shutdown(); err != nil {
			log.Fatalf("FATAL: shutdown server failed: %+v\n", err)
		}
		db.Close()

		log.Println("Server gracefully stopped")
		os.Exit(0)
	}()

	// showtime
	if err := srv.Serve(*hostname, *httpsPort, *tlsCert, *tlsKey); err != nil {
		db.Close()
		log.Fatalf("FATAL: starting server failed: %+v\n", err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: lara [flags]")
	fmt.Fprintln(os.Stderr, "flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func envVars() {
	cmd.UintVar(httpsPort, "LARA_HTTPS_PORT")
	cmd.StringVar(hostname, "LARA_HOSTNAME")
	cmd.StringVar(dbUser, "LARA_DB_USER")
	cmd.StringVar(dbPass, "LARA_DB_PASS")
	cmd.StringVar(dbHost, "LARA_DB_HOST")
	cmd.UintVar(dbPort, "LARA_DB_PORT")
	cmd.StringVar(dbName, "LARA_DB_NAME")
	cmd.StringVar(dbSSLMode, "LARA_DB_SSL_MODE")
	cmd.StringVar(rsaPriv, "LARA_RSA_PRIV")
	cmd.StringVar(rsaPub, "LARA_RSA_PUB")
	cmd.StringVar(tlsKey, "LARA_TLS_KEY")
	cmd.StringVar(tlsCert, "LARA_TLS_CERT")
	cmd.StringVar(wwwRoot, "LARA_WWW_ROOT")
}
