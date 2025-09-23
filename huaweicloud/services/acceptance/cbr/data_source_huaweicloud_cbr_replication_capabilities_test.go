package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReplicationCapabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cbr_replication_capabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceReplicationCapabilities_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.replication_destinations.#"),
				),
			},
		},
	})
}

const testDataSourceReplicationCapabilities_basic = `
data "huaweicloud_cbr_replication_capabilities" "test" {
}
`
