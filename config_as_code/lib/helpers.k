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