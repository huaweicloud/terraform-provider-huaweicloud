package ces

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMetricData_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_ces_metric_data.real_time"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)

	dataSource2 := "data.huaweicloud_ces_metric_data.aggregate"
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	earlyTime := time.Now().UTC().Add(-6 * time.Minute)
	earlyTimeString := earlyTime.Format("2006-01-02 15:04:05")
	earlyTimeRFCString := earlyTime.Format(time.RFC3339)

	currentTime := time.Now().UTC()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")

	baseConfig := testDataSourceMetricData_base(earlyTimeRFCString)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetricData_basic(baseConfig, earlyTimeString, currentTimeString),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource1, "datapoints.0.average"),
					resource.TestCheckResourceAttrSet(dataSource1, "datapoints.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSource1, "datapoints.0.unit"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
				),
			},
			{
				Config: testDataSourceMetricData_aggregate(baseConfig, earlyTimeString, currentTimeString),
				Check: resource.ComposeTestCheckFunc(
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource2, "datapoints.0.max"),
					resource.TestCheckResourceAttrSet(dataSource2, "datapoints.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSource2, "datapoints.0.unit"),
					resource.TestCheckOutput("is_aggregate_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceMetricData_basic(baseConfig, earlyTimeString, currentTimeString string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_metric_data" "real_time" {
  namespace   = "You.APP"
  metric_name = "cpu_util"
  dim_0       = "platform_id,test_platform_id"
  dim_1       = "instance_id,test_instance_id"
  dim_2       = "cpu_type,test_cpu_type"
  filter      = "average" 
  period      = 1  
  from        = "%[2]s"
  to          = "%[3]s"

  depends_on = [huaweicloud_ces_metric_data_add.point]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_metric_data.real_time.datapoints) >= 5 
}
`, baseConfig, earlyTimeString, currentTimeString)
}

func testDataSourceMetricData_aggregate(baseConfig, earlyTimeString, currentTimeString string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_metric_data" "aggregate" {
  namespace   = "You.APP"
  metric_name = "cpu_util"
  dim_0       = "platform_id,test_platform_id"
  dim_1       = "instance_id,test_instance_id"
  dim_2       = "cpu_type,test_cpu_type"
  filter      = "max"
  period      = 300
  from        = "%[2]s"
  to          = "%[3]s"

  depends_on = [huaweicloud_ces_metric_data_add.point]
}

output "is_aggregate_filter_useful" {
  value = length(data.huaweicloud_ces_metric_data.aggregate.datapoints) >= 1
}
`, baseConfig, earlyTimeString, currentTimeString)
}

func testDataSourceMetricData_base(earlyTimeString string) string {
	return fmt.Sprintf(`
variable "collect_time" {
  default = "%[1]s"
}
  
resource "huaweicloud_ces_metric_data_add" "point" {
  count = 5
  
  metric {
    namespace   = "You.APP"
    metric_name = "cpu_util"
	
    dimensions {
      name  = "platform_id"
      value = "test_platform_id"
    }
  
    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }
	
    dimensions {
      name  = "cpu_type"
      value = "test_cpu_type"
    }
  }
  
  ttl          = 600
  collect_time = formatdate("YYYY-MM-DD HH:mm:ss", timeadd(var.collect_time, format("%%ds", count.index * 60)))
  value        = 0.5 + count.index * 0.1
  unit         = "%%"
  type         = "float"
}
`, earlyTimeString)
}
