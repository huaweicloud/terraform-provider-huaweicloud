package ces

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceBatchMetricData_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_batch_metric_data.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBatchMetricData_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.0.dimensions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.0.dimensions.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "data_points.0.unit"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceBatchMetricData_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 22.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceBatchMetricData_basic(name string) string {
	earlyTime := time.Now().UTC().Add(-2 * time.Minute)
	earlyTimeString := earlyTime.Format("2006-01-02 15:04:05")
	currentTime := time.Now().UTC().Add(3 * time.Minute)
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_batch_metric_data" "test" {
  depends_on = [huaweicloud_compute_instance.test]

  namespace        = "SYS.ECS"
  metric_name      = "cpu_util"
  metric_dimension = "instance_id"
}

data "huaweicloud_ces_batch_metric_data" "time_filter" {
  depends_on = [huaweicloud_compute_instance.test]

  namespace        = "SYS.ECS"
  metric_name      = "cpu_util"
  metric_dimension = "instance_id"
  from             = "%s"
  to               = "%s"
}
output "time_filter_is_useful" {
  value = length(data.huaweicloud_ces_batch_metric_data.time_filter.data_points) > 0 
}
`, testDataSourceBatchMetricData_base(name), earlyTimeString, currentTimeString)
}
