resource "kilhog_network" "production" {
  name        = "production"
  description = "Production network"

  tags = {
    environment = "production"
  }
}
