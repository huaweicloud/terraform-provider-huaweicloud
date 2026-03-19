package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEipCount_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_eip_count.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEipCount_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "eip_protected"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_protected_self"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_total"),
				),
			},
		},
	})
}

func testDataSourceEipCount_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cfw_eip_count" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}
`, testAccDatasourceFirewalls_basic())
}
