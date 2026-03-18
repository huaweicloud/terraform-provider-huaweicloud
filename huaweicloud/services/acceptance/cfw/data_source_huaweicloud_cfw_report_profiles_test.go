package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReportProfiles_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_report_profiles.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceReportProfiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "total"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.profile_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.last_time"),
				),
			},
		},
	})
}

func testDataSourceReportProfiles_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "%[1]s"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "daily"
  name              = "%[1]s"
  topic_urn         = huaweicloud_smn_topic.test.topic_urn
  send_period       = "3"
  subscription_type = "0"
  status            = "0"
  description       = "test description"
}
`, name, acceptance.HW_CFW_INSTANCE_ID)
}

func testDataSourceReportProfiles_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cfw_report_profiles" "test" {
  fw_instance_id = "%[2]s"
  category       = huaweicloud_cfw_report_profile.test.category
}
`, testDataSourceReportProfiles_base(name), acceptance.HW_CFW_INSTANCE_ID)
}
