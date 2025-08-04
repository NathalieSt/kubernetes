import ..shared

get_larger_body_size_if_required = lambda service: shared.ClusterService -> str {
    result = ""
    if service.gateway_config.requires_larger_body == True:
        result = """
            # Increase client body size for file uploads
            client_max_body_size 100M;
        """

    result
}

get_websocket_support_if_required = lambda service: shared.ClusterService -> str {
    result = ""
    if service.gateway_config.requires_websockets == True:
        result = """
            # WebSocket support
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";

            # Use the upgrade mapping variable
            proxy_set_header Connection $connection_upgrade;
        """

    result
}

create_conf_from_services = lambda services: [shared.ClusterService] -> str {
    configmap_upstreams = "\n".join([""" 
    upstream ${service.name} {
        server ${service.cluster_url}:${service.port};
    }
    """ for service in services])

    configmap_server_names = ",".join(["${service.gateway_config.dns_name}.netbird.selfhosted" for service in services])

    configmap_servers = "\n".join([""" 
    server {
        listen 443 ssl;
        http2 on;
        server_name ${service.gateway_config.dns_name}.netbird.selfhosted;
        
        ssl_certificate /etc/nginx/ssl/tls.crt;
        ssl_certificate_key /etc/nginx/ssl/tls.key;
        
        # Security headers
        add_header Strict-Transport-Security "max-age=63072000" always;
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";

        ${get_larger_body_size_if_required(service)}
        
        location / {
            proxy_pass http://${service.name};
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-Host $server_name;
            proxy_set_header X-Forwarded-Port $server_port;

            ${get_websocket_support_if_required(service)}
        }
    }
    """ for service in services])

    """
        # WebSocket upgrade mapping
        map $http_upgrade $connection_upgrade {
            default upgrade;
            '' close;
        }
        ${configmap_upstreams}
        server {
            listen 80;
            server_name ${configmap_server_names};
            return 301 https://$server_name$request_uri;
        }
        ${configmap_servers}
    """
}
