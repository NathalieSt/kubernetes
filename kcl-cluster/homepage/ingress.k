import shared
import .variables

homepage_ingress = shared.TailscaleIngress {
    metadata: {
        name: "homepage-ingress"
    }
    spec: {
        ingressClassName = "tailscale"
        rules = [
            {
                http = {
                    paths = [
                        {
                            path = "/"
                            pathType = "Prefix"
                            backend = {
                                service = {
                                    name = "homepage"
                                    port = {
                                        number = 3000
                                    }
                                }
                            }
                        }
                    ]
                }
            }
        ]
        tls: [
            {
                hosts: [
                    "homepage"
                ]
            }
        ]
    }
}
