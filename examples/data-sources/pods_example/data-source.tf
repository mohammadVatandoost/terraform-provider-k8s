terraform {
  required_providers {
    k8s = {
      source = "registry.terraform.io/mohammadvatandoost/k8s"
    }
  }
}

provider "k8s" {
  # example configuration here
}

data "k8s_pods" "pos_test" {
  namespace = "test"
}

output "pos_test" {
  value = data.k8s_pods.pos_test
}