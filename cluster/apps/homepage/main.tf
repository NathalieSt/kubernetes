terraform {
  required_providers {
    kubernetes = {
      source  = "opentofu/kubernetes"
      version = "2.35.1"
    }
  }
}

resource "kubernetes_deployment" "homepage" {
  metadata {
    name      = "homepage"
    namespace = "homepage"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "homepage"
      }
    }

    template {
      metadata {
        labels = {
          app = "homepage"
        }
      }

      spec {
        container {
          name  = "homepage"
          image = "gitea-web.tail495c5f.ts.net/gitea_admin/custom-homer:latest"

          port {
            container_port = 8080
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "homepage" {
  metadata {
    name      = "homepage"
    namespace = "homepage"
  }

  spec {
    selector = {
      app = "homepage"
    }

    port {
      port        = 8080
      target_port = 8080
    }
  }
}

resource "kubernetes_ingress_v1" "homepage" {
  metadata {
    name      = "homepage-ingress"
    namespace = "homepage"
  }

  spec {
    ingress_class_name = "tailscale"

    rule {
      http {
        path {
          path      = "/"
          path_type = "Prefix"

          backend {
            service {
              name = kubernetes_service.homepage.metadata[0].name
              port {
                number = 8080
              }
            }
          }
        }
      }
    }
  }
}
