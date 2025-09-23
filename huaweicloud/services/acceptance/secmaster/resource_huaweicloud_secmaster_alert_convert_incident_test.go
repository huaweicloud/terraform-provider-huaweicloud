package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAlertConvertIncident_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterAlertId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertConvertIncident_basic(name),
			},
		},
	})
}

func testAccAlertConvertIncident_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert_convert_incident" "test" {
  workspace_id = "%[1]s"
  ids          = ["%[2]s"]
  title        = "%[3]s"

  incident_type {
    category      = "DDoS"
    incident_type = "ACK Flood"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_ALERT_ID, name)
}
