package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcHostedConnects_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_dc_hosted_connects.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcHostedConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcHostedConnects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.hosting_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.vlan"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.provider"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.provider_status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.port_type"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.location"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.peer_location"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.apply_time"),
					resource.TestCheckResourceAttrSet(dataSource, "hosted_connects.0.create_time"),

					resource.TestCheckOutput("hosting_id_filter_is_useful", "true"),
					resource.TestCheckOutput("hosted_connect_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcHostedConnects_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_hosted_connect" "test" {
  name               = "%s"
  description        = "This is a demo"
  resource_tenant_id = "%s"
  hosting_id         = "%s"
  vlan               = 441
  bandwidth          = 10
}
`, name, acceptance.HW_DC_RESOURCE_TENANT_ID, acceptance.HW_DC_HOSTTING_ID)
}

func testDataSourceDcHostedConnects_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_hosted_connects" "test" {
  depends_on = [huaweicloud_dc_hosted_connect.test]
}

locals {
  hosting_id = data.huaweicloud_dc_hosted_connects.test.hosted_connects.0.hosting_id
}
data "huaweicloud_dc_hosted_connects" "hosting_id_filter" {
  hosting_id = [data.huaweicloud_dc_hosted_connects.test.hosted_connects.0.hosting_id]
}

output "hosting_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_hosted_connects.hosting_id_filter.hosted_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_hosted_connects.hosting_id_filter.hosted_connects[*].hosting_id : 
      v == local.hosting_id]
  )
}

locals {
  hosted_connect_id = data.huaweicloud_dc_hosted_connects.test.hosted_connects.0.id
}
data "huaweicloud_dc_hosted_connects" "hosted_connect_id_filter" {
  hosted_connect_id = [data.huaweicloud_dc_hosted_connects.test.hosted_connects.0.id]
}
output "hosted_connect_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_hosted_connects.hosted_connect_id_filter.hosted_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_hosted_connects.hosted_connect_id_filter.hosted_connects[*].id : 
      v == local.hosted_connect_id]
  )
}

locals {
  name = data.huaweicloud_dc_hosted_connects.test.hosted_connects.0.name
}
data "huaweicloud_dc_hosted_connects" "name_filter" {
  name = [data.huaweicloud_dc_hosted_connects.test.hosted_connects.0.name]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dc_hosted_connects.name_filter.hosted_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_hosted_connects.name_filter.hosted_connects[*].name : 
      v == local.name]
  )
}

data "huaweicloud_dc_hosted_connects" "sort_filter" {
  depends_on = [huaweicloud_dc_hosted_connect.test]

  sort_key = "name"
  sort_dir = ["desc"]
}
output "sort_filter_is_useful" {
  value = length(data.huaweicloud_dc_hosted_connects.sort_filter.hosted_connects) > 0
}
`, testDataSourceDcHostedConnects_base(name))
}
