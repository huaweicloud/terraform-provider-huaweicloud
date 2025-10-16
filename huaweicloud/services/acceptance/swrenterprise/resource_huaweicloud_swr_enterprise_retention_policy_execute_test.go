package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseRetentionPolicyExecute_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseRetentionPolicyExecute_basic(rName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccSwrEnterpriseRetentionPolicyExecute_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_retention_policy_execute" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_retention_policy.test.policy_id
  dry_run        = true
}
`, testAccSwrEnterpriseRetentionPolicy_basic(rName))
}
