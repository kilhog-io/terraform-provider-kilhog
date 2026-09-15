resource "kilhog_network" "production" {
  name = "production"
}

resource "kilhog_subnet" "dmz" {
  network_id  = kilhog_network.production.id
  name        = "dmz"
  description = "DMZ subnet"
  address     = "10.0.0.0"
  prefix      = 24
  type        = "ipv4"
}

resource "kilhog_subnet" "apps" {
  network_id       = kilhog_network.production.id
  parent_subnet_id = kilhog_subnet.dmz.id
  name             = "apps"
  prefix           = 25
  type             = "ipv4"
}
