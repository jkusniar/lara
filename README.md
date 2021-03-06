# Lara

[![Build Status](https://travis-ci.org/jkusniar/lara.svg?branch=master)](https://travis-ci.org/jkusniar/lara)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://raw.githubusercontent.com/jkusniar/lara/master/LICENSE)

Veterinary practice support software - back-end services.

## Install

Requires go 1.9+

```
go install github.com/jkusniar/lara/cmd/lara
go install github.com/jkusniar/lara/cmd/lara-ctl
```

## Run

Runtime requirements:
* system timezone is used to translate date/time strings in application

## Development

Requires `make` utility. Install build dependencies if building first time:

* `make build-deps`
* `make`

Will install build dependencies and run golint/go vet, unit tests and build binaries.
Optionally install additional development dependencies using `make dev-deps`.

## Tests

```
make test
```

In order to run tests on local machine a postgresql instance 
must be running. Following environment variables are honored for tests:

* `POSTGRES_HOST` (default "localhost")
* `POSTGRES_DB` (default "lara_test")
* `POSTGRES_USER` (default "postgres")
* `POSTGRES_PASSWORD` (default "")
* `POSTGRES_PORT` (default 5432)

# License

The project license is specified in LICENSE

Lara is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Lara is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
