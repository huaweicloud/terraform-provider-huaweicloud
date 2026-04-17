package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAclRuleImportStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_acl_rule_import_status.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAclRuleImportStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
				),
			},
		},
	})
}

func testDataSourceAclRuleImportStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}

data "huaweicloud_cfw_acl_rule_import_status" "test" {
  depends_on = [data.huaweicloud_cfw_firewalls.test]
  object_id  = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
