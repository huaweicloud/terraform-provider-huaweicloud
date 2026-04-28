package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRdsSubscriptionRegenerate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRdsSubscriptionRegenerate_basic(),
			},
		},
	})
}

func testAccResourceRdsSubscriptionRegenerate_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_subscriptions" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_subscription_regenerate" "test" {
  instance_id     = "%[1]s"
  subscription_id = data.huaweicloud_rds_subscriptions.test.subscriptions[0].id
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
