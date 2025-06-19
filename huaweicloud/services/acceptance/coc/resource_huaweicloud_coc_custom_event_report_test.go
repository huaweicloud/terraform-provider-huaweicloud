package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocCustomEventReport_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocIntegrationKey(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocCustomEventReport_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocCustomEventReport_basic(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_coc_custom_event_report" "test" {
  integration_key = "%[1]s"
  alarm_id        = "%[2]s"
  alarm_name      = "%[3]s"
  alarm_level     = "Critical"
  time            = 1709118444540
  namespace       = "shanghai"
  application_id  = "%[2]s"
  alarm_desc      = "%[3]s"
  alarm_source    = "coc"
  region_id       = "cn-north-4"
  resource_name   = "%[3]s"
  resource_id     = "%[2]s"
  url             = "https://xxx.com"
  alarm_status    = "alarm"
  additional      = jsonencode({
    "key": "test"
  })
}
`, acceptance.HW_COC_INTEGRATION_KEY, randUUID, name)
}
