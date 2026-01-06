package generators

import (
	"fmt"
)

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

// Other common shared variables
var PostgresCredsSecret = "postgres-creds-secret"
var MatrixPGCredsSecret = "matrix-pg-creds-secret"
var ForgejoPGCredsSecret = "forgejo-pg-creds-secret"
var MariaDBCredsSecret = "mariadb-creds-secret"
var NetbirdSecretName = "netbird-setup-key-vault-secret"
var NetbirdAPIKeySecretName = "netbird-mgmt-api-key"
var HetznerAPITokenSecretName = "hetzner-api-token-secret"
var NetbirdDomainBase = "cloud.nathalie-stiefsohn.eu"
var SynapseSecretName = "synpase-secret"
var DiscordBridgeSecretName = "discord-bridge-secret"
var WhatsappBridgeSecretName = "whatsapp-bridge-secret"
var ElasticSearchAdminSecretName = "elastic-search-admin-secret"
var ElasticSearchVectorSecretName = "elastic-search-vector-secret"
