package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccResourceGaussDbAspCollect_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGaussdbAspCollect_basic(rName),
			},
		},
	})
}

func testAccResourceGaussdbAspCollect_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "test_123456"
  port              = "9000"
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccResourceGaussdbAspCollect_basic(name string) string {
	startTime := time.Now().UTC().Format("2024-12-31T23:59:59+0800")
	endTime := time.Now().UTC().Add(1 * time.Hour).Format("2024-12-31T23:59:59+0800")
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_asp_collect" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, testAccResourceGaussdbAspCollect_base(name), startTime, endTime)
}
