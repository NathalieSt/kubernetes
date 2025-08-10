import ....lib.istio
import ....lib.k8s_wrapper
import ...shared
import yaml

out_dir = shared.cluster_path + "istio/security/"
# Filenames
kustomization_file = "kustomization.yaml"
peer_authenthication_file = "peer-authenthication.yaml"
# Configs
service_name = "istio-config"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    peer_authenthication_file
])

peer_authenthication: istio.PeerAuthenthication = {
    metadata: {
        name: "enforce-mtls-in-mesh"
    }
    spec: {
        mtls: {
            mode: "STRICT"
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(peer_authenthication, out_dir + peer_authenthication_file)
