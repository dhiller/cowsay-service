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

set -x
set -euo pipefail

SCRIPT_PATH="$(
    cd "$(dirname "$BASH_SOURCE[0]")/"
    echo "$(pwd)/"
)"

image_file=$(xq -r '.domain.devices.disk[0].source."@file"' vm/cowsay-service-centos-stream-9-vm.xml)

# port forward kubevirtci cdi-uploadproxy to enable uploading
kubectl port-forward -n cdi service/cdi-uploadproxy-nodeport 18443:443 &
trap 'pkill -f "kubectl port-forward"' SIGINT SIGTERM

# upload image (note that we want a block volume with RWX access for live migration)
kubectl virt image-upload pvc cowsay-service-vm-disk \
    --image-path=${image_file} --size=18Gi \
    --volume-mode=block \
    --access-mode=ReadWriteMany \
    --insecure \
    --uploadproxy-url=https://127.0.0.1:18443
