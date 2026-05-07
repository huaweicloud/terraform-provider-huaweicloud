package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsConnections_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_connections.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsConnections_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.connection_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.ssl.0.ssl_link"),
				),
			},
		},
	})
}

const testAccDataSourceDrsConnections_basic = `
data "huaweicloud_drs_connections" "test" {
}
`
