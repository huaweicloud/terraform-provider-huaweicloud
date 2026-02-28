package rds

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPublicationCandidates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_publication_candidates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsPublicationCandidates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_publications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_publications.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_publications.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_publications.0.publication_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_publications.0.publication_name"),
					resource.TestCheckOutput("publication_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("publication_instance_name_filter_is_useful", "true"),
					resource.TestCheckOutput("publication_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsPublicationCandidates_base(name string) string {
	currentDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2022_SE"
  instance_mode = "single"
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  tde_enabled       = true

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name, currentDate)
}

func testDataSourceRdsPublicationCandidates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_publications" "test" {
  instance_id = "%[2]s"
}

data "huaweicloud_rds_publication_candidates" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_publication_candidates" "publication_instance_id_filter" {
  instance_id             = huaweicloud_rds_instance.test.id
  publication_instance_id = "%[2]s"
}

locals {
  publication_instance_id =  "%[2]s"
  publication_instance_id_filter_data = data.huaweicloud_rds_publication_candidates.publication_instance_id_filter
}

output "publication_instance_id_filter_is_useful" {
  value = length(local.publication_instance_id_filter_data.instance_publications) > 0 && alltrue(
    [for v in local.publication_instance_id_filter_data.instance_publications[*].instance_id :
      v == local.publication_instance_id]
  )
}

data "huaweicloud_rds_publication_candidates" "publication_instance_name_filter" {
  instance_id               = huaweicloud_rds_instance.test.id
  publication_instance_name = "%[3]s"
}

locals {
  publication_instance_name = "%[3]s"
  publication_instance_name_filter_data = data.huaweicloud_rds_publication_candidates.publication_instance_name_filter
}

output "publication_instance_name_filter_is_useful" {
  value = length(local.publication_instance_name_filter_data.instance_publications) > 0 && alltrue(
    [for v in local.publication_instance_name_filter_data.instance_publications[*].instance_name :
      v == local.publication_instance_name]
  )
}

data "huaweicloud_rds_publication_candidates" "publication_name_filter" {
  instance_id      = huaweicloud_rds_instance.test.id
  publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
}

locals {
  publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
  publication_name_filter_data = data.huaweicloud_rds_publication_candidates.publication_name_filter
}

output "publication_name_filter_is_useful" {
  value = length(local.publication_name_filter_data.instance_publications) > 0 && alltrue(
    [for v in local.publication_name_filter_data.instance_publications[*].publication_name : v == local.publication_name]
  )
}
`, testDataSourceRdsPublicationCandidates_base(name), acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_INSTANCE_NAME)
}
