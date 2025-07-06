terraform {
  required_providers {
    unbound = {
      source = "registry.terraform.io/platylab/unbound"
    }
  }
}

provider "unbound" {
  host             = "adm-dns1.adm.platylab.com:22"
  username         = "platy"
  private_key_path = "/home/areaute/test/ssh_key"
}

data "unbound_local_zone" "ipa_platylab_com" {
  name = "ipa.platylab.com"
}

output "ipa_platylab_com" {
  value = data.unbound_local_zone.ipa_platylab_com
}
