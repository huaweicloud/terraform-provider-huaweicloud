package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchDeleteAclRules_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting the firewall ID, ACL rule ID, and enterprise project ID for CFW,
			// and ensuring that they are all under the same enterprise project.
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwAclRuleId(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBatchDeleteAclRules_basic(),
			},
		},
	})
}

func testBatchDeleteAclRules_base() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testBatchDeleteAclRules_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_batch_delete_acl_rules" "test" {
  object_id             = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  rule_ids              = ["%[2]s"]
  fw_instance_id        = "%[3]s"
  enterprise_project_id = "%[4]s"
}
`, testBatchDeleteAclRules_base(), acceptance.HW_CFW_ACL_RULE_ID, acceptance.HW_CFW_INSTANCE_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
