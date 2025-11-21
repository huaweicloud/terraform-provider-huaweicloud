package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceMfaDevices_basic(t *testing.T) {
	userId := acceptance.HW_IDENTITY_CENTER_USER_ID
	rName := "data.huaweicloud_identitycenter_mfa_devices.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMfaDevices_basic(userId),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "mfa_devices.#"),
					resource.TestCheckResourceAttrSet(rName, "mfa_devices.0.device_id"),
					resource.TestCheckResourceAttrSet(rName, "mfa_devices.0.device_name"),
					resource.TestCheckResourceAttrSet(rName, "mfa_devices.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "mfa_devices.0.mfa_type"),
					resource.TestCheckResourceAttrSet(rName, "mfa_devices.0.registered_date"),
				),
			},
		},
	})
}

func testAccDatasourceMfaDevices_basic(userId string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}
 
data "huaweicloud_identitycenter_mfa_devices" "test"{
  depends_on        = [data.huaweicloud_identitycenter_instance.test]
  user_id           = "%s"
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`, userId)
}
