resource "kilhog_network" "production" {
  name        = "production"
  description = "Production network"

  tags = [
    {
      key   = "environment"
      value = "production"
    },
  ]
}
