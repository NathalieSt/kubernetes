package shared

import "fmt"

// NFS Remote Config
var NFSRemoteClass = "nfs-remote"
var NFSRemoteServer = fmt.Sprintf("remote-fs.%v", NetbirdDomainBase)
var NFSRemoteShare = "/mnt/HC_Volume_103061115"

// NFS Local
var NFSLocalClass = "nfs-local"
var NFSLocalServer = "raspberry-pi-5-0"
var NFSLocalShare = "/mnt/external_ssd"

// NFS Debian Config
var DebianStorageClass = "debian-storage"
var DebianServer = "debian"
var DebianShare = "/srv/nfs"
