package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsInstanceExpiredKeyScanHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_instance_expired_key_scan_histories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsInstanceExpiredKeyScanHistories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.scan_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.num"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.started_at"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.finished_at"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcsInstanceExpiredKeyScanHistories_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

resource "huaweicloud_dcs_instance_expired_key_scan_task" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, name)
}

func testDataSourceDcsInstanceExpiredKeyScanHistories_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_expired_key_scan_histories" "test" {
  depends_on = [huaweicloud_dcs_instance_expired_key_scan_task.test]

  instance_id = huaweicloud_dcs_instance.test.id
}

locals{
  status = "success"
}
data "huaweicloud_dcs_instance_expired_key_scan_histories" "status_filter" {
  depends_on = [huaweicloud_dcs_instance_expired_key_scan_task.test]

  instance_id = huaweicloud_dcs_instance.test.id
  status      = "success"
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dcs_instance_expired_key_scan_histories.status_filter.records) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_instance_expired_key_scan_histories.status_filter.records[*].status : v == local.status]  
  )
}
`, testDataSourceDcsInstanceExpiredKeyScanHistories_base(name))
}
