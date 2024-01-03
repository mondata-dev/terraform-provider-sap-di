terraform {
  required_providers {
    sapdi = {
      source = "mondata.de/terraform/sap-di"
    }
  }
}

provider "sapdi" {
  username = "admin"
  password = "test123"
  host     = "http://localhost:8080"
}

data "sapdi_factsheet" "test" {
  metadata = {
    uri           = "/XYZ/012/ABCD"
    connection_id = "P40_XYZ"
  }
}
