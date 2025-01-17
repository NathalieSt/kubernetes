import shared

forgejo_web_ingress = shared.TailscaleIngress {
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
                    "homepage"
                ]
            }
        ]
    }
}
