package apig

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceMetricData_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_apig_instance_metric_data.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstanceMetricData_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "datapoints.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dcName, "datapoints.0.timestamp"),
					resource.TestCheckResourceAttrSet(dcName, "datapoints.0.unit"),
					resource.TestCheckResourceAttrSet(dcName, "datapoints.0.average"),
				),
			},
		},
	})
}

func testAccDataInstanceMetricData_basic() string {
	now := time.Now()
	to := now.UnixMilli()
	from := now.Add(-24 * time.Hour).UnixMilli()

	return fmt.Sprintf(`
data "huaweicloud_apig_instance_metric_data" "test" {
  instance_id = "%[1]s"
  dim         = "inbound_eip"
  metric_name = "upstream_bandwidth"
  from        = "%[2]d"
  to          = "%[3]d"
  period      = 300
  filter      = "average"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, from, to)
}
