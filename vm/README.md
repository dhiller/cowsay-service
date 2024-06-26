# VM sources

Base of the vm is the libvirt qcow2 image for [CentOS Stream 9]

[CentOS Stream 9]: https://cloud.centos.org/centos/9-stream/x86_64/images/

## create a VM as service outside KubeVirt

host requirements for running a virtual machine:

```bash
sudo dnf install libvirt
sudo dnf install virt-install
```

### automated setup

**NOTE: I haven't tested the script, I am just providing it to give an idea on how it should work.**

Execute the task item that calls [./vm/prepare-vm.sh](./prepare-vm.sh)

```bash
go-task vm:prepare
```

### manual setup

#### basic setup

use centos stream 9 latest qcow2 as base


[create vm from image](https://smoogespace.blogspot.com/2022/02/how-to-install-centos-stream-9-cloud.html)

```bash
sudo virt-install --name cowsay-service-centos-stream-9-vm --memory 2048 \
  --vcpus 2 --disk ./CentOS-Stream-GenericCloud-x86_64-9-latest.x86_64.qcow2 \
  --import --os-variant centos-stream9 \
  --network default --console pty,target_type=serial --graphics vnc \
  --cloud-init root-password-generate=on,disable=on,ssh-key=$HOME/.ssh/id_rsa.pub
```

#### ssh setup
```bash
$ sudo virsh net-list --all
…
$ sudo virsh net-dhcp-leases default
```

##### service setup

install epel:
```bash
dnf install -y epel-release
```

install cowsay:
```bash
dnf install -y cowsay
```

build cowsay service locally
```bash
scp cowsay-service to vm /usr/local/bin
```

##### install cowsay as a systemd service

```bash
# echo <<EOF > /lib/systemd/system/cowsay.service
[Unit]
Description=
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/cowsay-service

[Install]
WantedBy=multi-user.target
EOF

# systemctl enable cowsay.service

# systemctl start cowsay.service
```

##### Test it
locally

```bash
curl -v localhost:8080/cowsays
```

from outside the vm

```bash
curl -v <ip-address-vm>:8080/cowsays
```

