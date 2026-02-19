sudo apt install nfs-kernel-server

sudo nano /etc/exports

# edit this line into file:
#/home               client_ip(rw,sync,no_root_squash,no_subtree_check)


sudo systemctl restart nfs-kernel-server

sudo systemctl status nfs-kernel-server

# 100.127.1.1/16 <- is netbird address specific
sudo ufw allow in from 100.87.0.0/16 to any port 111
sudo ufw allow in from 100.87.0.0/16 to any port 2049
# 13025 <- this ports needs to be manually fixed, its usually dynamic
# change it in /etc/nfs.conf under [mountd]
sudo ufw allow in from 100.87.0.0/16/16 to any port 13025