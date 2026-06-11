package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbInstanceMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_metrics.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbInstanceMetrics_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.metric"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.unit"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.datapoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.datapoints.0.datapoint_name"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.datapoints.0.datapoint_values.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.timestamps.#"),
					resource.TestCheckOutput("component_id_filter", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbInstanceMetrics_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "test_1234"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

`, common.TestVpc(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbInstanceMetrics_basic(name string) string {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startTime := startOfDay.UnixMilli()

	endOfDay := startOfDay.Add(24 * time.Hour)
	endTime := endOfDay.UnixMilli()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_all_instances_metrics" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
}

data "huaweicloud_gaussdb_metric_group_metrics" "test" {
  group_name = "CPUMEMORY"
}

data "huaweicloud_gaussdb_instance_metrics" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  start_time  = "%[2]d"
  end_time    = "%[3]d"
  metric      = data.huaweicloud_gaussdb_metric_group_metrics.test.metric_names[*].metric
  node_id     = data.huaweicloud_gaussdb_all_instances_metrics.test.instances[0].nodes[*].id
}

data "huaweicloud_gaussdb_instance_metrics" "component_id_filter" {
  instance_id  = huaweicloud_gaussdb_instance.test.id
  start_time   = "%[2]d"
  end_time     = "%[3]d"
  metric       = data.huaweicloud_gaussdb_metric_group_metrics.test.metric_names[*].metric
  node_id      = data.huaweicloud_gaussdb_all_instances_metrics.test.instances[0].nodes[*].id
  component_id = data.huaweicloud_gaussdb_all_instances_metrics.test.instances[0].nodes[0].component_ids
}
output "component_id_filter" {
  value = length(data.huaweicloud_gaussdb_instance_metrics.component_id_filter.metrics) > 0
}
`, testDataSourceGaussdbInstanceMetrics_base(name), startTime, endTime)
}
