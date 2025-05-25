# Steps for setting up the vps
## Create a non-root sudo user
```sh
adduser nathi
```
```sh
usermod -aG sudo nathi
```
```sh
getent group sudo
```
## Add a SSH Key for the new user
Add the generated ssh public key to `~/.ssh/authorized_keys`
```sh
ssh-copy-id -i ./ssh_key.pub nathi@vps-server-ip
```
### Test authenthication via SSH Key
Connect with private key
```sh
ssh -i ./ssh_key nathi@vps-server-ip
```
### Disable password access
```sh
sudo nano /etc/ssh/sshd_config
```
Set `PasswordAuthentication` to `no`

## Disable root user access
### Login on server
```sh
sudo nano /etc/passwd
```
Change
`root:x:0:0:root:/root:/bin/bash`
to
`root:x:0:0:root:/root:/sbin/nologin`
### SSH Access
```sh
sudo nano /etc/ssh/sshd_config
```
Set line `PermitRootLogin` to `no`
```sh
sudo systemctl restart sshd
```
## Install ufw
```sh
sudo bash ufw.sh
```
## Install fail2ban
```sh
sudo bash fail2ban.sh
```
## Install netbird
```sh
sudo bash netbird.sh
```