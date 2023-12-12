package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFactsheetsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sapdi_factsheet" "test" {
					metadata = {
						uri = "/XYZ/012/ABCD"
						connection_id = "P40_XYZ"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify metadata
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.connection_id", "P40_XYZ"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.name", "ABCD"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.uri", "/XYZ/012/ABCD"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.descriptions.#", "1"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.descriptions.0.origin", "REMOTE"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.descriptions.0.type", "SHORT"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "metadata.descriptions.0.value", "Characteristic"),

					// Verify number of columns returned
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.#", "2"),
					// Verify the first column to ensure all attributes are set
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.0.name", "MANDT"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.0.type", "STRING"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.0.descriptions.#", "1"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.0.descriptions.0.origin", "REMOTE"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.0.descriptions.0.type", "SHORT"),
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "columns.0.descriptions.0.value", "Client"),

					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.sapdi_factsheet.test", "id", "placeholder"),
				),
			},
		},
	})
}
