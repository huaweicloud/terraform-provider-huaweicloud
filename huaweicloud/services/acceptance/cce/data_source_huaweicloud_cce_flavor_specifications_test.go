package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceFlavorSpecifications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_flavor_specifications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCceFlavorSpecifications_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.node_capacity"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.is_sold_out"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.is_support_multi_az"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.available_master_flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.available_master_flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_flavor_specs.0.available_master_flavors.0.azs.#"),
				),
			},
		},
	})
}

const testDataSourceCceFlavorSpecifications_basic = `
data "huaweicloud_cce_flavor_specifications" "test" {
  cluster_type = "VirtualMachine"
}
`
