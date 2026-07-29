data "kilhog_network" "production" {
  name = "production"
}

data "kilhog_subnet" "dmz" {
  network_id = data.kilhog_network.production.id
  name       = "dmz"
}
