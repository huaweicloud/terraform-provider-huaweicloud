package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssClusters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_clusters.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCssClusters_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.datastore.0.version"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_type_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_version_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCssClusters_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_clusters" "test" {
  depends_on = [
    huaweicloud_css_cluster.test,
    huaweicloud_css_logstash_cluster.test,
  ]
}

locals {
  cluster_id     = data.huaweicloud_css_clusters.test.clusters[0].id
  name           = data.huaweicloud_css_clusters.test.clusters[0].name
  engine_type    = data.huaweicloud_css_clusters.test.clusters[0].datastore[0].type
  engine_version = data.huaweicloud_css_clusters.test.clusters[0].datastore[0].version
}

data "huaweicloud_css_clusters" "filter_by_id" {
  cluster_id = local.cluster_id
}

data "huaweicloud_css_clusters" "filter_by_name" {
  name = local.name
}

data "huaweicloud_css_clusters" "filter_by_engine_type" {
  engine_type = local.engine_type
}

data "huaweicloud_css_clusters" "filter_by_engine_version" {
  engine_version = local.engine_version
}

locals {
  list_by_id             = data.huaweicloud_css_clusters.filter_by_id.clusters
  list_by_name           = data.huaweicloud_css_clusters.filter_by_name.clusters
  list_by_engine_type    = data.huaweicloud_css_clusters.filter_by_engine_type.clusters
  list_by_engine_version = data.huaweicloud_css_clusters.filter_by_engine_version.clusters
}

output "id_filter_is_useful" {
  value = length(local.list_by_id) > 0 && alltrue(
    [for v in local.list_by_id[*].id : v == local.cluster_id]
  )
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}

output "engine_type_filter_is_useful" {
  value = length(local.list_by_engine_type) > 0 && alltrue(
    [for v in local.list_by_engine_type[*].datastore[0].type : v == local.engine_type]
  )
}

output "engine_version_filter_is_useful" {
  value = length(local.list_by_engine_version) > 0 && alltrue(
    [for v in local.list_by_engine_version[*].datastore[0].version : v == local.engine_version]
  )
}
`, testDataSourceCssClusters_data_basic(name))
}

func testDataSourceCssClusters_data_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"
  security_mode  = true
  https_enabled  = true
  password       = "Test@passw0rd"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  backup_strategy {
    keep_days   = 7
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  public_access {
    bandwidth         = 5
    whitelist_enabled = true
    whitelist         = "116.204.111.47"
  }

  kibana_public_access {
    bandwidth         = 5
    whitelist_enabled = true
    whitelist         = "116.204.111.47"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s_1"
  engine_version = "7.10.0"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}
`, testAccCssBase(name), name)
}
