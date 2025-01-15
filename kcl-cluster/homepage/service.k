import .variables

homepage_service = {
    apiVersion = "v1"
    kind = "Service"
    metadata = {
        name = variables.service_name
        labels = {
            "app.kubernetes.io/name" = "homepage"
        }
        annotations = None
    }
    spec = {
        $type = "ClusterIP"
        ports = [
            {
                port = 3000
                targetPort = "http"
                $protocol = "TCP"
                name = "http"
            }
        ]
        selector = {
            "app.kubernetes.io/name" = "homepage"
        }
    }
}
