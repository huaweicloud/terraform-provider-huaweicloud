package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to the lack of test scenarios, this test case only verifies that the API request is normal, not the response result.
func TestAccCcAttackProtectionQPSDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_cc_attack_protection_qps.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCcAttackProtectionQps_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

const testCcAttackProtectionQps_basic = `
data "huaweicloud_aad_cc_attack_protection_qps" "test" {
  recent = "1month"
}
`
