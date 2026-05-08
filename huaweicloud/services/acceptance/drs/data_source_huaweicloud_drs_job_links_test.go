package drs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceJobLinks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_job_links.test"
		rName          = acceptance.RandomAccResourceName()
		dbName         = acceptance.RandomAccResourceName()
		pwd            = acceptance.RandomPassword()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceJobLinks_basic(rName, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "job_links.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("job_type_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.job_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.engine_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.net_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.task_types.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.job_direction"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.cluster_modes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.source_endpoint_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_links.0.target_endpoint_type"),
				),
			},
		},
	})
}

const testSecgroupRule string = `
resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "3306,9092"
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.0.0/16"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.0.0/16"
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`

func testAccDataSourceJobLinks_base(name, dbName, pwd, action, startTime string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

%[3]s

%[4]s

resource "huaweicloud_drs_job" "test" {
  name           = "%[5]s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "vpc"
  migration_type = "FULL_INCR_TRANS"
  description    = "%[5]s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%[6]s"
    vpc_id      = huaweicloud_vpc.test.id
    subnet_id   = huaweicloud_vpc_subnet.test.id
  }

  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[6]s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  action     = "%[7]s"
  start_time = "%[8]s"

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, testSecgroupRule, sourceDb, destDb, name, pwd, action, startTime)
}

func testAccDataSourceJobLinks_basic(rName, dbName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_drs_job_links" "test" {
  job_type = "migration"

  depends_on = [
    huaweicloud_drs_job.test
  ]
}

output "job_type_filter_is_useful" {
  value = length(data.huaweicloud_drs_job_links.test.job_links) > 0
}
`, testAccDataSourceJobLinks_base(rName, dbName, pwd, "start", ""))
}
