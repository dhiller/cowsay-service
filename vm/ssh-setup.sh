#!/usr/bin/env bash
#
# MIT License
#
# Copyright (c) 2024 Daniel Hiller
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
#

echo "WARNING: THIS SCRIPT IS NOT INTENDED TO RUN ON A HOST MACHINE!"

set -euo pipefail
set -x

dnf install -y epel-release
dnf install -y cowsay
cat <<EOF > /lib/systemd/system/cowsay.service
[Unit]
Description=
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/cowsay-service

[Install]
WantedBy=multi-user.target
EOF
systemctl enable cowsay.service
systemctl start cowsay.service
