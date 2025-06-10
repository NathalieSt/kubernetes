import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api

create_nginx_config_map = lambda name: str, service_name: str, service_namespace: str, service_port: int -> core_api.ConfigMap {
    core_api.ConfigMap {
        metadata: meta.ObjectMeta {
            name: name
        },
        data: {
            "default.conf": """\
                # WebSocket upgrade mapping
                map $http_upgrade $connection_upgrade {
                    default upgrade;
                    '' close;
                }
                
                upstream ${service_name} {
                    server ${service_name}.${service_namespace}.svc.cluster.local:${service_port};
                }

                # HTTP server - redirect to HTTPS
                server {
                    listen 80;
                    server_name ${service_name}.netbird.selfhosted;
                    return 301 https://$server_name$request_uri;
                }

                # HTTPS server
                server {
                    listen 443 ssl;
                    http2 on;
                    server_name ${service_name}.netbird.selfhosted;
                    
                    ssl_certificate /etc/nginx/ssl/tls.crt;
                    ssl_certificate_key /etc/nginx/ssl/tls.key;
                    
                    # Security headers (relaxed for Authentik)
                    add_header Strict-Transport-Security "max-age=63072000" always;
                    add_header X-Content-Type-Options nosniff;
                    add_header X-XSS-Protection "1; mode=block";
                    
                    # Increase client body size for file uploads
                    client_max_body_size 100M;
                    
                    # Timeout settings for long-running requests
                    proxy_connect_timeout 600s;
                    proxy_send_timeout 600s;
                    proxy_read_timeout 600s;
                    
                    location / {
                        proxy_pass http://${service_name};
                        
                        # Essential headers for Authentik
                        proxy_set_header Host $http_host;
                        proxy_set_header X-Real-IP $remote_addr;
                        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                        proxy_set_header X-Forwarded-Proto $scheme;
                        proxy_set_header X-Forwarded-Host $http_host;
                        proxy_set_header X-Forwarded-Port $server_port;
                        
                        # Additional headers for proper authentication flow
                        proxy_set_header X-Original-URL $scheme://$http_host$request_uri;
                        proxy_set_header X-Forwarded-Uri $request_uri;
                        
                        # WebSocket support
                        proxy_http_version 1.1;
                        proxy_set_header Upgrade $http_upgrade;
                        proxy_set_header Connection $connection_upgrade;
                        
                        # Buffering settings for better performance
                        proxy_buffering off;
                        proxy_request_buffering off;
                        
                        # CORS headers for API requests
                        add_header Access-Control-Allow-Origin $http_origin always;
                        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS, PATCH" always;
                        add_header Access-Control-Allow-Headers "Authorization, Content-Type, Accept, Origin, User-Agent, DNT, Cache-Control, X-Mx-ReqToken, Keep-Alive, X-Requested-With, If-Modified-Since, X-CSRF-Token" always;
                        add_header Access-Control-Allow-Credentials true always;
                        
                        # Handle preflight requests
                        if ($request_method = 'OPTIONS') {
                            add_header Access-Control-Allow-Origin $http_origin always;
                            add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS, PATCH" always;
                            add_header Access-Control-Allow-Headers "Authorization, Content-Type, Accept, Origin, User-Agent, DNT, Cache-Control, X-Mx-ReqToken, Keep-Alive, X-Requested-With, If-Modified-Since, X-CSRF-Token" always;
                            add_header Access-Control-Allow-Credentials true always;
                            add_header Access-Control-Max-Age 86400;
                            add_header Content-Length 0;
                            add_header Content-Type text/plain;
                            return 204;
                        }
                    }
                }
            
            """ 
            }
    }
}