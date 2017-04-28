#!/bin/sh
#
#   Copyright (C) 2016-2017 Contributors as noted in the AUTHORS file
#
#   This file is part of lara, veterinary practice support software.
#
#   This program is free software: you can redistribute it and/or modify
#   it under the terms of the GNU General Public License as published by
#   the Free Software Foundation, either version 3 of the License, or
#   (at your option) any later version.
#
#   This program is distributed in the hope that it will be useful,
#   but WITHOUT ANY WARRANTY; without even the implied warranty of
#   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#   GNU General Public License for more details.
#
#   You should have received a copy of the GNU General Public License
#   along with this program.  If not, see <http://www.gnu.org/licenses/>.
#

# create user _lara
useradd -L daemon -d /var/lara -s /sbin/nologin -g =uid _lara

# create dist dir
mkdir /var/lara
chown _lara:_lara /var/lara

# copy files to dist dir
cp lara.rsa /var/lara/lara.rsa
chown _lara:_lara /var/lara/lara.rsa

cp lara.rsa.pub /var/lara/lara.rsa.pub
chown _lara:_lara /var/lara/lara.rsa.pub

cp key.pem /var/lara/key.pem
chown _lara:_lara /var/lara/key.pem

cp cert.pem /var/lara/cert.pem
chown _lara:_lara /var/lara/cert.pem

cp lara /var/lara/lara
chmod +x /var/lara/lara
chown _lara:_lara /var/lara/lara

cp lara-ctl /var/lara/lara-ctl
chmod +x /var/lara/lara-ctl
chown _lara:_lara /var/lara/lara-ctl

# copy rc.d script
cp lara.rc /etc/rc.d/lara
chmod +x /etc/rc.d/lara
