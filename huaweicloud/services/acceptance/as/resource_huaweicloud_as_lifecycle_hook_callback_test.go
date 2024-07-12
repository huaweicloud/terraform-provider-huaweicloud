package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Because an instance with a status of **HANGING** cannot be obtained through resources and datasource.
// So using environment variables to inject parameters for testing.
func TestAccLifecycleHookCallBack_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingGroupID(t)
			acceptance.TestAccPreCheckASLifecycleActionKey(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleHookCallBack_withActionKey("ABANDON"),
			},
		},
	})
}

func testAccLifecycleHookCallBack_withActionKey(action string) string {
	return fmt.Sprintf(`
resource "huaweicloud_as_lifecycle_hook_callback" "test" {
  scaling_group_id        = "%[1]s"
  lifecycle_action_result = "%[2]s"
  lifecycle_action_key    = "%[3]s"
}
`, acceptance.HW_AS_SCALING_GROUP_ID, action, acceptance.HW_AS_LIFECYCLE_ACTION_KEY)
}

func TestAccLifecycleHookCallBack_withInstanceIdAndHookName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingGroupID(t)
			acceptance.TestAccPreCheckASINSTANCEID(t)
			acceptance.TestAccPreCheckASLifecycleHookName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleHookCallBack_withInstanceIdAndHookName("CONTINUE"),
			},
		},
	})
}

func testAccLifecycleHookCallBack_withInstanceIdAndHookName(action string) string {
	return fmt.Sprintf(`
resource "huaweicloud_as_lifecycle_hook_callback" "test" {
  scaling_group_id        = "%[1]s"
  lifecycle_action_result = "%[2]s"
  instance_id             = "%[3]s"
  lifecycle_hook_name     = "%[4]s" 
}
`, acceptance.HW_AS_SCALING_GROUP_ID, action, acceptance.HW_AS_INSTANCE_ID, acceptance.HW_AS_LIFECYCLE_HOOK_NAME)
}
