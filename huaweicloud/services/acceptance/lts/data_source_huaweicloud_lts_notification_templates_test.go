package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNotificationTemplates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_lts_notification_templates.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNotificationTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "templates.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.locale"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.source"),
					resource.TestMatchResourceAttr(dataSource, "templates.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "templates.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testDataSourceNotificationTemplates_expectError(),
				// The domain ID not exist.
				ExpectError: regexp.MustCompile("Failed to query notification template."),
			},
		},
	})
}

func testDataSourceNotificationTemplates_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_notification_template" "test" {
  name        = "%[1]s"
  source      = "LTS"
  locale      = "en-us"
  description = "Created by terraform script"

  templates {
    sub_type = "sms"
    content  = "This content of sub-template."
  }
}

data "huaweicloud_lts_notification_templates" "test" {
  depends_on = [
    "huaweicloud_lts_notification_template.test"
  ]

  domain_id = "%[2]s"
}
`, name, acceptance.HW_DOMAIN_ID)
}

func testDataSourceNotificationTemplates_expectError() string {
	return `
data "huaweicloud_lts_notification_templates" "not_found_domain_id" {
  domain_id = "not_found_domain_id"
}
`
}
