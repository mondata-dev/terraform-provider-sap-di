terraform {
  required_providers {
    sapdi = {
      source = "mondata.de/terraform/sap-di"
    }
  }
}

provider "sapdi" {
}

data "sapdi_factsheet" "test" {}
