terraform {
  required_providers {
    kubernetes = {
      source  = "opentofu/kubernetes"
      version = "2.35.1"
    }
  }
}

resource "kubernetes_namespace" "homepage" {
  metadata {
    name = "homepage"
  }
}
