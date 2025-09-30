package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSettingTwoFactorLoginConfig_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testSettingTwoFactorLoginConfig_basic(name),
			},
		},
	})
}

func testSettingTwoFactorLoginConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "display_%[1]s"
}

resource "huaweicloud_hss_setting_two_factor_login_config" "test" {
  enabled               = true
  auth_type             = "sms"
  host_id_list          = ["%[2]s"]
  topic_display_name    = huaweicloud_smn_topic.test.display_name
  topic_urn             = huaweicloud_smn_topic.test.id
  enterprise_project_id = "0"
}
`, name, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
