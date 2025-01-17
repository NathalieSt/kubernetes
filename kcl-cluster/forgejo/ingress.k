import shared

forgejo_ingress = shared.TailscaleIngress {
    metadata: {
        name: "forgejo-ingress"
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
                                    name = "forgejo"
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
                    "forgejo"
                ]
            }
        ]
    }
}
