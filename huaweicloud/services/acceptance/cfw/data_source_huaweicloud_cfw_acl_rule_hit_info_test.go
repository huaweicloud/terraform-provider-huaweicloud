package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAclRuleHitInfo_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_acl_rule_hit_info.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a firewall instance ID and an ACL rule ID.
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwAclRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAclRuleHitInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rule_hit_count"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rule_last_hit_time"),
				),
			},
		},
	})
}

func testDataSourceAclRuleHitInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_acl_rule_hit_info" "test" {
  fw_instance_id = "%s"
  rule_ids       = ["%s"]
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_ACL_RULE_ID)
}
