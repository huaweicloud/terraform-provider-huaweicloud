package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceIncidentHandle_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocApplicationID(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testIncidentHandle_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testIncidentHandle_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "time_sleep" "wait_10_seconds" {
  depends_on      = [huaweicloud_coc_incident.test]
  create_duration = "10s"
}

resource "huaweicloud_coc_incident_handle" "test" {
  incident_num = huaweicloud_coc_incident.test.id
  operator     = "%[2]s"
  operate_key  = "acceptedIncident1"
  depends_on   = [time_sleep.wait_10_seconds]
}`, testIncident_basic(rName), acceptance.HW_USER_ID)
}
