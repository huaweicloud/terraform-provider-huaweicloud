package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIpBlacklist_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_ip_blacklist.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIpBlacklist_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.records.0.effect_scope.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.records.0.import_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.records.0.import_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.total"),
				),
			},
		},
	})
}

func testDataSourceIpBlacklist_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_ip_blacklist" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
