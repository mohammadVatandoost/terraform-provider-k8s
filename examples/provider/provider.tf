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