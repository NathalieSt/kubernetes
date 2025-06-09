import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api

create_nginx_config_map = lambda name: str, service_name: str, service_namespace: str, service_port: int -> core_api.ConfigMap {
    core_api.ConfigMap {
        metadata: meta.ObjectMeta {
            name: name
        },
        data: {
            "default.conf": """\
                upstream ${service_name} {
                    server localhost:${service_port};
                }

                # HTTP server - redirect to HTTPS
                server {
                    listen 80;
                    server_name ${service_name}.netbird.selfhosted;
                    return 301 https://$server_name$request_uri;
                }

                # HTTPS server for Uptime Kuma
                server {
                    listen 443 ssl;
                    http2 on;
                    server_name ${service_name}.netbird.selfhosted;
                    
                    ssl_certificate /etc/nginx/ssl/tls.crt;
                    ssl_certificate_key /etc/nginx/ssl/tls.key;
                    
                    # Security headers
                    add_header Strict-Transport-Security "max-age=63072000" always;
                    add_header X-Frame-Options DENY;
                    add_header X-Content-Type-Options nosniff;
                    add_header X-XSS-Protection "1; mode=block";
                    
                    location / {
                        proxy_pass http://${service_name};
                        proxy_set_header Host $host;
                        proxy_set_header X-Real-IP $remote_addr;
                        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                        proxy_set_header X-Forwarded-Proto $scheme;
                        proxy_set_header X-Forwarded-Host $server_name;
                        proxy_set_header X-Forwarded-Port $server_port;
                        
                        # WebSocket support
                        proxy_http_version 1.1;
                        proxy_set_header Upgrade $http_upgrade;
                        proxy_set_header Connection "upgrade";
                    }
                }
            
            """ 
            }
    }
}