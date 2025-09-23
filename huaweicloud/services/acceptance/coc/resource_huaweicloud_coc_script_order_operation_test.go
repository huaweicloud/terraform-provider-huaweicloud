package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocScriptOrderOperation_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocScriptOrderOperation_basic(rName, acceptance.HW_COC_INSTANCE_ID),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocScriptOrderOperation_base(name, instanceID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script" "test" {
  name        = "%[1]s"
  description = "a new demo script"
  risk_level  = "MEDIUM"
  version     = "1.0.1"
  type        = "SHELL"

  content = <<EOF
#! /bin/bash
echo "hello $${name}@$${company}!"
sleep 2m
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
  parameters {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
}

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%[2]s"
  timeout      = 600
  execute_user = "root"

  parameters {
    name  = "name"
    value = "somebody"
  }
  parameters {
    name  = "company"
    value = "HuaweiCloud"
  }

  timeouts {
    create = "1m"
  }
}`, name, instanceID)
}

func testCocScriptOrderOperation_basic(name, instanceId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_coc_script_order_operation" "test" {
  execute_uuid   = huaweicloud_coc_script_execute.test.id
  batch_index    = 1
  instance_id    = 1
  operation_type = "CANCEL_ORDER"
}
`, testCocScriptOrderOperation_base(name, instanceId))
}
