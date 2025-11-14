package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseReplicationPolicyExecute_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseReplicationPolicyExecute_basic(rName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccSwrEnterpriseReplicationPolicyExecute_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_replication_policy_execute" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  policy_id   = huaweicloud_swr_enterprise_replication_policy.test.policy_id
}

resource "huaweicloud_swr_enterprise_replication_policy_execution_stop" "test" {
  instance_id  = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  execution_id = huaweicloud_swr_enterprise_replication_policy_execute.test.execution_id
}
`, testAccSwrEnterpriseReplicationPolicy_basic(rName))
}
