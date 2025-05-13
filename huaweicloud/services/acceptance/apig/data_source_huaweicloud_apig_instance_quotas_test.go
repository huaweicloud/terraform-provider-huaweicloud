package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_apig_instance_quotas.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceQuotas_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_used_set_and_valid", "true"),
					resource.TestCheckOutput("is_config_value_set_and_valid", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.config_id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.config_name"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.remark"),
					resource.TestMatchResourceAttr(dataSource, "quotas.0.config_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceInstanceQuotas_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

data "huaweicloud_apig_instance_quotas" "test" {
  instance_id = "%[1]s"

  depends_on = [huaweicloud_apig_group.test]
}

# 'APIGROUP_NUM_LIMIT' means the number of API groups is limited.
locals {
  app_group_quota = try([for v in data.huaweicloud_apig_instance_quotas.test.quotas : v if v.config_name == "APIGROUP_NUM_LIMIT"][0], {})
}

output "is_used_set_and_valid" {
  value = lookup(local.app_group_quota, "used", 0) > 0
}

output "is_config_value_set_and_valid" {
  value = lookup(local.app_group_quota, "config_value", "") != ""
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
