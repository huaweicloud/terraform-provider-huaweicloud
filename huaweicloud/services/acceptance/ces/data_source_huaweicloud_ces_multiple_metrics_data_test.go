package ces

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMultipleMetricsData_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_ces_multiple_metrics_data.real_time"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)

	dataSource2 := "data.huaweicloud_ces_multiple_metrics_data.aggregate"
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	earlyTime := time.Now().UTC().Add(-6 * time.Minute)
	earlyTimeString := earlyTime.Format("2006-01-02 15:04:05")
	earlyTimeRFCString := earlyTime.Format(time.RFC3339)

	currentTime := time.Now().UTC()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")

	baseConfig := testDataSourceMultipleMetricsData_base(earlyTimeRFCString)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMultipleMetricsData_basic(baseConfig, earlyTimeString, currentTimeString),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource1, "data.0.namespace", "YOU.APP"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.metric_name", "cpu_util"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.dimensions.0.name", "cpu_type"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.dimensions.0.value", "test_cpu_type"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.dimensions.1.name", "instance_id"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.dimensions.1.value", "test_instance_id"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.dimensions.2.name", "platform_id"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.dimensions.2.value", "test_platform_id"),
					resource.TestCheckResourceAttr(dataSource1, "data.0.unit", "%"),
					resource.TestCheckResourceAttr(dataSource1, "data.1.namespace", "MINE.APP"),
					resource.TestCheckResourceAttr(dataSource1, "data.1.metric_name", "mem_util"),
					resource.TestCheckResourceAttr(dataSource1, "data.1.dimensions.0.name", "instance_id"),
					resource.TestCheckResourceAttr(dataSource1, "data.1.dimensions.0.value", "test_instance_id"),
					resource.TestCheckResourceAttr(dataSource1, "data.1.dimensions.1.name", "memory_type"),
					resource.TestCheckResourceAttr(dataSource1, "data.1.dimensions.1.value", "test_memory_type"),
					resource.TestCheckResourceAttrSet(dataSource1, "data.0.datapoints.0.average"),
					resource.TestCheckResourceAttrSet(dataSource1, "data.0.datapoints.0.timestamp"),
					resource.TestCheckOutput("is_raw_filter_useful", "true"),
				),
			},
			{
				Config: testDataSourceMultipleMetricsData_aggregate(baseConfig, earlyTimeString, currentTimeString),
				Check: resource.ComposeTestCheckFunc(
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource2, "data.0.namespace", "YOU.APP"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.metric_name", "cpu_util"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.dimensions.0.name", "cpu_type"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.dimensions.0.value", "test_cpu_type"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.dimensions.1.name", "instance_id"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.dimensions.1.value", "test_instance_id"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.dimensions.2.name", "platform_id"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.dimensions.2.value", "test_platform_id"),
					resource.TestCheckResourceAttr(dataSource2, "data.0.unit", "%"),
					resource.TestCheckResourceAttr(dataSource2, "data.1.namespace", "MINE.APP"),
					resource.TestCheckResourceAttr(dataSource2, "data.1.metric_name", "mem_util"),
					resource.TestCheckResourceAttr(dataSource2, "data.1.dimensions.0.name", "instance_id"),
					resource.TestCheckResourceAttr(dataSource2, "data.1.dimensions.0.value", "test_instance_id"),
					resource.TestCheckResourceAttr(dataSource2, "data.1.dimensions.1.name", "memory_type"),
					resource.TestCheckResourceAttr(dataSource2, "data.1.dimensions.1.value", "test_memory_type"),
					resource.TestCheckResourceAttrSet(dataSource2, "data.0.datapoints.0.sum"),
					resource.TestCheckResourceAttrSet(dataSource2, "data.0.datapoints.0.timestamp"),
					resource.TestCheckOutput("is_aggregate_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceMultipleMetricsData_basic(baseConfig, earlyTimeString, currentTimeString string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_multiple_metrics_data" "real_time" {
  metrics {
    namespace   = "YOU.APP"
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

  metrics {
    namespace   = "MINE.APP"
    metric_name = "mem_util"

    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }

    dimensions {
      name  = "memory_type"
      value = "test_memory_type"
    }
  }

  from   = "%[2]s"
  to     = "%[3]s"
  period = 1
  filter = "average"

  depends_on = [
    huaweicloud_ces_metric_data_add.YourData,
    huaweicloud_ces_metric_data_add.MyData
  ]
}

output "is_raw_filter_useful" {
  value = alltrue([
    for v in data.huaweicloud_ces_multiple_metrics_data.real_time.data[*].datapoints : length(v) >= 2
  ])
}
`, baseConfig, earlyTimeString, currentTimeString)
}

func testDataSourceMultipleMetricsData_aggregate(baseConfig, earlyTimeString, currentTimeString string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_multiple_metrics_data" "aggregate" {
  metrics {
    namespace   = "YOU.APP"
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
  
  metrics {
    namespace   = "MINE.APP"
    metric_name = "mem_util"
  
    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }
  
    dimensions {
      name  = "memory_type"
      value = "test_memory_type"
    }
  }
  
  from   = "%[2]s"
  to     = "%[3]s"
  period = 300
  filter = "sum"
  
  depends_on = [
    huaweicloud_ces_metric_data_add.YourData,
    huaweicloud_ces_metric_data_add.MyData
  ]
}
  
output "is_aggregate_filter_useful" {
  value = alltrue([
    for v in data.huaweicloud_ces_multiple_metrics_data.aggregate.data[*].datapoints : length(v) >= 1
  ])
}
`, baseConfig, earlyTimeString, currentTimeString)
}

func testDataSourceMultipleMetricsData_base(earlyTimeString string) string {
	return fmt.Sprintf(`
variable "collect_time" {
  default = "%[1]s"
}
  
resource "huaweicloud_ces_metric_data_add" "YourData" {
  count = 2
  
  metric {
    namespace   = "YOU.APP"
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

resource "huaweicloud_ces_metric_data_add" "MyData" {
  count = 2
  
  metric {
    namespace   = "MINE.APP"
    metric_name = "mem_util"
  
    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }
    
    dimensions {
      name  = "memory_type"
      value = "test_memory_type"
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
