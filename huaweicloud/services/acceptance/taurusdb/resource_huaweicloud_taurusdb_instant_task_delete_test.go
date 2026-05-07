package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBInstantTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBInstantTaskDelete_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccTaurusDBInstantTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_instant_task_delete" "test" {
  job_id = "%[1]s"
}`, acceptance.HW_TAURUSDB_JOB_ID)
}
