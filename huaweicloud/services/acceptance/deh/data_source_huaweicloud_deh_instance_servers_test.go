package deh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDehInstanceServers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_deh_instance_servers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDehDedicatedHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDehInstanceServers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "servers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.addresses.%"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.metadata.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.updated"),
				),
			},
		},
	})
}

func testDataSourceDehInstanceServers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_deh_instance_servers" "test" {
  dedicated_host_id = "%[1]s"
}
`, acceptance.HW_DEH_DEDICATED_HOST_ID)
}
