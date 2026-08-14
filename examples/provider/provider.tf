terraform {
  required_providers {
    kilhog = {
      source = "kilhog-io/kilhog"
    }
  }
}

provider "kilhog" {
  base_url = "http://localhost:8080"
  # api_key = var.kilhog_api_key
}
