package workspace

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppServerMetricData_basic(t *testing.T) {
	var (
		byAverage   = "data.huaweicloud_workspace_app_server_metric_data.average"
		dcByAverage = acceptance.InitDataSourceCheck(byAverage)

		byMax   = "data.huaweicloud_workspace_app_server_metric_data.max"
		dcByMax = acceptance.InitDataSourceCheck(byMax)

		byMin   = "data.huaweicloud_workspace_app_server_metric_data.min"
		dcByMin = acceptance.InitDataSourceCheck(byMin)

		bySum   = "data.huaweicloud_workspace_app_server_metric_data.sum"
		dcBySum = acceptance.InitDataSourceCheck(bySum)

		byVariance   = "data.huaweicloud_workspace_app_server_metric_data.variance"
		dcByVariance = acceptance.InitDataSourceCheck(byVariance)

		currentTime = time.Now()
		startTime   = currentTime.Add(-1 * time.Hour).Format(time.RFC3339)
		endTime     = currentTime.Format(time.RFC3339)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppServerMetricData_serverIdNotFound(startTime, endTime),
				ExpectError: regexp.MustCompile("The cloud application server requested by the client was not found, and '" +
					".+" + "' is a non-existing cloud application server."),
			},
			{
				Config: testAccDataAppServerMetricData_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dcByAverage.CheckResourceExists(),
					resource.TestMatchResourceAttr(byAverage, "metrics.#", regexp.MustCompile(`[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byAverage, "metrics.0.metric_name"),
					resource.TestMatchResourceAttr(byAverage, "metrics.0.datapoints.#", regexp.MustCompile(`[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byAverage, "metrics.0.datapoints.0.average"),
					resource.TestCheckResourceAttrSet(byAverage, "metrics.0.datapoints.0.collection_time"),
					resource.TestCheckResourceAttrSet(byAverage, "metrics.0.datapoints.0.unit"),
					dcByMax.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byMax, "metrics.0.datapoints.0.max"),
					dcByMin.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byMin, "metrics.0.datapoints.0.min"),
					dcBySum.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(bySum, "metrics.0.datapoints.0.sum"),
					dcByVariance.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byVariance, "metrics.0.datapoints.0.variance"),
				),
			},
		},
	})
}

func testAccDataAppServerMetricData_serverIdNotFound(startTime, endTime string) string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_server_metric_data" "test" {
  server_id   = "%[1]s"
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = "%[2]s"
  to          = "%[3]s"
  period      = 1
  filter      = "average"
}
`, randomId, startTime, endTime)
}

func testAccDataAppServerMetricData_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_server_metric_data" "average" {
  server_id   = "%[1]s"
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = "%[2]s"
  to          = "%[3]s"
  period      = 1
  filter      = "average"
}

data "huaweicloud_workspace_app_server_metric_data" "max" {
  server_id   = "%[1]s"
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = "%[2]s"
  to          = "%[3]s"
  period      = 1
  filter      = "max"
}

data "huaweicloud_workspace_app_server_metric_data" "min" {
  server_id   = "%[1]s"
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = "%[2]s"
  to          = "%[3]s"
  period      = 1
  filter      = "min"
}

data "huaweicloud_workspace_app_server_metric_data" "sum" {
  server_id   = "%[1]s"
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = "%[2]s"
  to          = "%[3]s"
  period      = 1
  filter      = "sum"
}

data "huaweicloud_workspace_app_server_metric_data" "variance" {
  server_id   = "%[1]s"
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = "%[2]s"
  to          = "%[3]s"
  period      = 1
  filter      = "variance"
}
`, acceptance.HW_WORKSPACE_APP_SERVER_ID, startTime, endTime)
}
