# Get a factsheet by URI and Connection ID.
data "sapdi_factsheet" "test" {
  metadata = {
    uri           = "/XYZ/012/ABCD"
    connection_id = "P40_XYZ"
  }
}
