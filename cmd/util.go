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

package cmd

import (
	"fmt"
	"os"
	"strconv"
)

// UintVar sets uint variable from environment variable if defined
func UintVar(v *uint, env string) {
	s := os.Getenv(env)
	if len(s) != 0 {
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Env variable %s conversion error: %s. Ignoring.\n",
				env, err)
		} else if i < 0 {
			fmt.Fprintf(os.Stderr, "Env variable %s is negative. Ignoring.\n", env)
		} else {
			*v = uint(i)
		}
	}
}

// StringVar sets string variable from environment variable if defined
func StringVar(v *string, env string) {
	s := os.Getenv(env)
	if len(s) != 0 {
		*v = s
	}
}

// CheckPortNum checks port number range
func CheckPortNum(port uint, name string) {
	if port == 0 || port > 65535 {
		fmt.Fprintf(os.Stderr, "%s value %d is out of range [1-65535]\n",
			name, port)
		os.Exit(2)
	}
}

// CheckFileExists checks filepath exists
func CheckFileExists(name string) {
	if _, err := os.Stat(name); err != nil {
		fmt.Fprintf(os.Stderr, "Required file %v not found: %v\n", name, err)
		os.Exit(2)
	}
}
