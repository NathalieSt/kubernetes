schema HelmCharts:
    gatus: str
    csi_driver_nfs: str
    vault: str

schema DockerImages:

docker_images: DockerImages = {}

helm_charts: HelmCharts = {
    gatus: "1.3.1"
    csi_driver_nfs: "4.11.0"
    vault: "0.30.1"
}
