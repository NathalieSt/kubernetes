import shared
import .variables

homepage_ingress = shared.TailscaleIngress {
    metadata: {
        name: "homepage_ingress"
    }
    spec: {
        defaultBackend: {
            service.name = variables.service_name
            service.port.number = 80
        }
        tls: [
            {
                hosts: [
                    "homepage"
                ]
            }
        ]
    }
}
