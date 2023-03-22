package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.k8s_pods.pos_test", "namespace", "test"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
terraform {
	required_providers {
	  k8s = {
		source = "mohammadvatandoost/k8s"
	  }
	}
  }

data "k8s_pods" "pos_test" {
	namespace = "test"
}
`
