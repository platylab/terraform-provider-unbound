
resource "unbound_local_zone" "example_eu" {
  name = "example.eu."
  type = "transparent"
}
