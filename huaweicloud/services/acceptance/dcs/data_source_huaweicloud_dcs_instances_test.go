package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDcsInstance_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcsInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.port", "6388"),
					resource.TestCheckResourceAttr(rName, "instances.0.flavor", "redis.ha.xu1.tiny.r2.128"),
				),
			},
		},
	})
}

func testAccDatasourceDcsInstance_base(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  maintain_begin     = "22:00:00"
  maintain_end       = "02:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, instanceName)
}

func testAccDatasourceDcsInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instances" "test" {
  name     = huaweicloud_dcs_instance.test.name
  status   = "RUNNING"
}
`, testAccDatasourceDcsInstance_base(name))
}
