# VM sources

Base of the vm is the libvirt qcow2 image for [CentOS Stream 9]

[CentOS Stream 9]: https://cloud.centos.org/centos/9-stream/x86_64/images/

#### create a VM as service outside KubeVirt


##### basic setup

use centos stream 9 latest qcow2 as base

$ sudo dnf install libvirt

$ sudo dnf install virt-install

[create vm from image](https://smoogespace.blogspot.com/2022/02/how-to-install-centos-stream-9-cloud.html)

```bash
$ sudo virt-install --name cowsay-service-centos-stream-9-vm --memory 2048 \

  --vcpus 2 --disk ./CentOS-Stream-GenericCloud-x86_64-9-latest.x86_64.qcow2 \

  --import --os-variant centos-stream9 \

  --network default --console pty,target_type=serial --graphics vnc \

  --cloud-init root-password-generate=on,disable=on,ssh-key=/home/dhiller/.ssh/id_rsa.pub
```

##### ssh setup
```bash
$ sudo virsh net-list --all
…
$ sudo virsh net-dhcp-leases
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
