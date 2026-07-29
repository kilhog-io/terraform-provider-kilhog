resource "terraform_data" "example" {
  input = "fake-string"

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.kilhog_action.example]
    }
  }
}

action "kilhog_action" "example" {
  config {
    configurable_attribute = "value"
  }
}
