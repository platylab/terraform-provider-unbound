
resource "unbound_local_data" "tofu_example_eu" {
  domain = "tofu.example.eu."
  type   = "A"
  value  = "3.5.45.126"
}
