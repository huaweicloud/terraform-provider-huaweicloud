package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocPublicScriptExecute_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
			acceptance.TestAccPreCheckCocScriptID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocPublicScriptExecute_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocPublicScriptExecute_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_public_script_execute" "test" {
  script_uuid  = "%s"
  timeout      = 300
  success_rate = 100
  execute_user = "root"
  script_params {
    param_name  = "action"
    param_value = "stop"
  }
  execute_batches {
    batch_index = 1
    target_instances {
      resource_id        = "%s"
      region_id          = "cn-south-1"
      cloud_service_name = "ECS"
      custom_attributes {
        key   = "key"
        value = "value"
      }
    }
    rotation_strategy = "CONTINUE"
  }
}
`, acceptance.HW_COC_SCRIPT_ID, acceptance.HW_COC_INSTANCE_ID)
}
