forgejo_release = {
    apiVersion: "helm.toolkit.fluxcd.io/v2"
    kind: "HelmRelease"
    metadata: {
        name: "forgejo"
    }
    spec: {
        releaseName: "forgejo"
        chart: {
            spec: {
                chart: "forgejo"
                sourceRef: {
                    kind: "HelmRepository"
                    name: "forgejo"
                }
            }
        }
        interval: "50m"
        install.remediation.retries: 3
        values: {
            actions.enabled = True
            actions.provisioning.enabled = True
            gitea.config.server.SSH_DOMAIN = "forgejo-ssh.tail495c5f.ts.net"
            gitea.config.server.ENABLE_PUSH_CREATE_USER = True
            "postgresql-ha": {
                enabled = False
            }
            postgresql.enabled = True
        }
    }
}
