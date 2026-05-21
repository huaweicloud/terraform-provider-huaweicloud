package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEmailTemplatesBatchAction_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_email_templates_batch_action.test"

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailTemplatesBatchAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "subscribe", "true"),
					resource.TestCheckResourceAttr(rName, "email_template_ids.#", "1"),
				),
			},
			{
				Config: testAccEmailTemplatesBatchAction_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "subscribe", "false"),
					resource.TestCheckResourceAttr(rName, "email_template_ids.#", "2"),
				),
			},
		},
	})
}

func testAccEmailTemplatesBatchAction_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_das_instance_group" "test" {
  datastore_type = "mysql"
  group_name     = "%[1]s"
  description    = "Created by terraform script"
}

resource "huaweicloud_das_email_template" "test" {
  count = 2

  datastore_type  = "mysql"
  name            = "%[1]s_${count.index}"
  groups          = [huaweicloud_das_instance_group.test.id]
  health_rank     = ["dangerous", "sub_healthy"]
  inspection_time = "00:00-00:00"
  send_time       = "08:00-10:00"
  time_zone       = "Asia/Shanghai"
  email           = "test@example.com"
}
`, name)
}

func testAccEmailTemplatesBatchAction_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_email_templates_batch_action" "test" {
  subscribe          = true
  email_template_ids = [huaweicloud_das_email_template.test[0].id]
}
`, testAccEmailTemplatesBatchAction_base(name))
}

func testAccEmailTemplatesBatchAction_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_email_templates_batch_action" "test" {
  subscribe          = false
  email_template_ids = huaweicloud_das_email_template.test[*].id

  enable_force_new = "true"
}
`, testAccEmailTemplatesBatchAction_base(name))
}
