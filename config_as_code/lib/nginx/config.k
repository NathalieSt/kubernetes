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
                    server ${service_name}.${service_namespace}.svc.cluster.local:${service_port};
                }

                # HTTP server - redirect to HTTPS
                server {
                    listen 80;
                    server_name ${service_name}.cluster.netbird.selfhosted;
                    return 301 https://$server_name$request_uri;
                }
            
            """ 
            }
    }
}