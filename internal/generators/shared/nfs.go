package shared

import "fmt"

// NFS Remote Config
var NFSRemoteClass = "nfs-remote"
var NFSRemoteServer = fmt.Sprintf("remote-fs.%v", NetbirdDomainBase)
var NFSRemoteShare = "/mnt/HC_Volume_103061115"

// NFS Local Config
var NFSLocalClass = "nfs-local-v2"
var NFSLocalServer = "raspberry-pi-5-0"
var NFSLocalShare = "/mnt/external_ssd"

// NFS Local Next
var NFSLocalClassNext = "nfs-local-next"
var NFSLocalServerNext = "raspberry-pi-5-0"
var NFSLocalShareNext = "/mnt/external_ssd"

// NFS Debian Config
var DebianStorageClass = "debian-storage"
var DebianServer = "debian"
var DebianShare = "/srv/nfs"
