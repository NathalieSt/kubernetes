sudo apt install nfs-kernel-server

sudo nano /etc/exports

# edit this line into file:
#/home               client_ip(rw,sync,no_root_squash,no_subtree_check)


sudo systemctl restart nfs-kernel-server

sudo systemctl status nfs-kernel-server