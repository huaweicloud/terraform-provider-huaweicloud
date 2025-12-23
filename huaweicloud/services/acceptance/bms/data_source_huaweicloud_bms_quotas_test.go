package bms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBmsQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_bms_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBmsQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.#"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_total_instances"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_server_group_members"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.total_server_groups_used"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_security_groups"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_image_meta"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.total_cores_used"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_total_keypairs"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_personality"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_security_group_rules"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_total_floating_ips"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.total_instances_used"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_total_ram_size"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_server_meta"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_personality_size"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_server_groups"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.total_floating_ips_used"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.max_total_cores"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.total_ram_used"),
					resource.TestCheckResourceAttrSet(dataSource, "absolute.0.total_security_groups_used"),
				),
			},
		},
	})
}

func testDataSourceBmsQuotas_basic() string {
	return `
data "huaweicloud_bms_quotas" "test" {}
`
}
