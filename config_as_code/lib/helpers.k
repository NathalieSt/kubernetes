resources = {
    requests: {
        cpu: "50m"
        memory: "64Mi"
    }
    limits:{
        cpu: "100m"
        memory: "128Mi"
    }
}

create_nginx_sidecar = lambda config_volume_name: str, tls_volume_name: str {
    sidecar = {
        name: "nginx"
        image: "nginx:alpine"
        ports: [
            { containerPort: 80},
            { containerPort: 443}
        ]
        volumeMounts: [
            { 
                name: config_volume_name
                mountPath: "/etc/nginx/conf.d"
            },
            { 
                name: tls_volume_name
                mountPath: "/etc/nginx/ssl"
                readOnly: True
            }
        ]
        resources: resources
    }
    sidecar
}

create_netbird_sidecar = lambda setup_key_secret_name: str, setup_key_secret_key: str, extra_dns_labels: str {
    sidecar = {
        name: "netbird-agent"
        image: "netbirdio/netbird:latest"
        env: [
            {
                name: "NB_SETUP_KEY"
                valueFrom: {
                    secretKeyRef: {
                        name: setup_key_secret_name
                        key: setup_key_secret_key 
                    }
                }
            },
            {
                name: "NB_HOSTNAME"
                value: "cluster"
            },
            {
                name: "NB_MANAGEMENT_URL"
                value: "https://netbird.nathalie-stiefsohn.eu"
            },
            {
                name: "NB_EXTRA_DNS_LABELS"
                value: extra_dns_labels
            }
        ]
        resources: resources
        securityContext: {
          privileged: True
        }
    }
}