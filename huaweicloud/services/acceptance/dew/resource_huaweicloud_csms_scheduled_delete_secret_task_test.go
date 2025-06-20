package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCsmsScheduledDeleteSecretTask_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCsmsScheduledDeleteSecretTask_basic(name),
			},
		},
	})
}

func testAccCsmsScheduledDeleteSecretTask_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name        = "%s"
  secret_text = "this is a password"
}

resource "huaweicloud_csms_scheduled_delete_secret_task" "test1" {
  secret_name = huaweicloud_csms_secret.test.name
  action      = "create"
}

resource "huaweicloud_csms_scheduled_delete_secret_task" "test2" {
  secret_name = huaweicloud_csms_secret.test.name
  action      = "cancel"

  depends_on = [huaweicloud_csms_scheduled_delete_secret_task.test1]
}
`, name)
}
