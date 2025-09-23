package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCfwIpsRuleModeChange_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceCfwIpsRuleModeChange_basic(),
			},
		},
	})
}

func testResourceCfwIpsRuleModeChange_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_ips_rule_mode_change" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  ips_ids   = [340710, 340922]
  status    = "CLOSE"
}
`, testAccDatasourceFirewalls_basic())
}
