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
set -euo pipefail

function help() {
    cat <<EOF
usage:  ${0} --help|-h      Print this help message
        ${0}                Run the setup script to prepare a machine for running
                            the cowsay service permanently.
EOF
}

function main() {
    if [[ $# -gt 0 ]]; then
        case $1 in
            --help|-h)
                help
                exit 0
                ;;
            *)
                ;;
        esac
    fi

    set -x

    # build cowsay service binary
    go build -o /tmp/cowsay-service ./cmd/cowsay-service/...

    # fetch centos stream 9 cloud image
    wget -o ./vm/CentOS-Stream-GenericCloud-x86_64-9-latest.x86_64.qcow2 https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-x86_64-9-latest.x86_64.qcow2

    # create vm from image
    sudo virt-install --name cowsay-service-centos-stream-9-vm --memory 2048 \
      --vcpus 2 --disk ./vm/CentOS-Stream-GenericCloud-x86_64-9-latest.x86_64.qcow2 \
      --import --os-variant centos-stream9 \
      --network default --console pty,target_type=serial --graphics vnc \
      --cloud-init root-password-generate=on,disable=on,ssh-key=$HOME/.ssh/id_rsa.pub

    sudo virsh net-list --all

    vm_ip_address="$(sudo virsh net-dhcp-leases default | grep -oE '([0-9]+\.){3}[0-9]+')"
    scp /tmp/cowsay-service root@${vm_ip_address}:/usr/local/bin/
    scp ./vm/ssh-setup.sh root@${vm_ip_address}:/tmp/

    ssh root@${vm_ip_address} 'chmod +x /tmp/ssh-setup.sh /usr/local/bin/cowsay-service'
    ssh root@${vm_ip_address} '/tmp/ssh-setup.sh'

    # test vm service call
    if ! curl --fail -v ${vm_ip_address}:8080/cowsays ; then
        echo "service call failed!"
        exit 1
    fi
}

main "$@"
