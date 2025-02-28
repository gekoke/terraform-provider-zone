terraform {
  required_providers {
    zone = {
      source = "gekoke/zone"
    }
  }
}

provider "zone" {
  username = var.zone_username
  api_key  = var.zone_api_key
}

variable "zone_username" {
  type = string
}

variable "zone_api_key" {
  type      = string
  sensitive = true
}

resource "zone_record_url" "this" {
  domain      = "grigorjan.ee"
  name        = "foo.grigorjan.ee"
  destination = "email.example.com"
  type        = 301
}
