terraform {
  required_providers {
    unbound = {
      source = "registry.terraform.io/platylab/unbound"
    }
  }
}

provider "unbound" {
  host = "https://plbnet-dns-eu01.adm.platylab.com:8091"
}

data "unbound_local_zone" "example_com" {
  id = 1
}


output "local_zone_id" {
  value = data.unbound_local_zone.example_com.id
}

output "local_zone_name" {
  value = data.unbound_local_zone.example_com.name
}
output "local_zone_type" {
  value = data.unbound_local_zone.example_com.type
}


data "unbound_local_data" "test_example_com" {
  id = 1
}


output "local_data_id" {
  value = data.unbound_local_data.test_example_com.id
}

output "local_data_domain" {
  value = data.unbound_local_data.test_example_com.domain
}
output "local_data_type" {
  value = data.unbound_local_data.test_example_com.type
}

output "local_data_value" {
  value = data.unbound_local_data.test_example_com.value
}


resource "unbound_local_zone" "example_eu" {
  name = "example.eu."
  type = "transparent"
}

resource "unbound_local_data" "tofu_example_eu" {
  domain = "tofu.example.eu."
  type   = "A"
  value  = "3.5.45.126"
}
