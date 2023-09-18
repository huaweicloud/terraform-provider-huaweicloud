package cdm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCdmClusters_basic(t *testing.T) {
	rName := "data.huaweicloud_cdm_clusters.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCdmClusters_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "clusters.0.status", "200"),
					resource.TestCheckResourceAttr(rName, "clusters.0.recent_event", "1"),
					resource.TestCheckResourceAttr(rName, "clusters.0.is_frozen", "0"),
					resource.TestCheckResourceAttr(rName, "clusters.0.is_auto_off", "false"),
					resource.TestCheckResourceAttr(rName, "clusters.0.is_failure_remind", "false"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.name"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.availability_zone"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.version"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.instances.#"),

					resource.TestCheckOutput("default_is_useful", "true"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("az_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("error_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCdmClusters_basic(name string) string {
	clustersConfig := testAccCdmCluster_basic(name)
	name2 := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cdm_cluster" "test2" {
  availability_zone = data.huaweicloud_availability_zones.test.names[1]
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name              = "%[2]s"
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}

data "huaweicloud_cdm_clusters" "test" {
  depends_on = [
    huaweicloud_cdm_cluster.test,
    huaweicloud_cdm_cluster.test2
  ]
}

output "default_is_useful" {
  value = length(data.huaweicloud_cdm_clusters.test.clusters) >= 2
}

data "huaweicloud_cdm_clusters" "name_filter" {
  name = "%[3]s"

  depends_on = [
    huaweicloud_cdm_cluster.test,
    huaweicloud_cdm_cluster.test2
  ]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_cdm_clusters.name_filter.clusters) > 0 && alltrue(
    [for v in data.huaweicloud_cdm_clusters.name_filter.clusters[*].name : v == "%[3]s"]
  )  
}

data "huaweicloud_cdm_clusters" "az_filter" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  depends_on = [
    huaweicloud_cdm_cluster.test,
    huaweicloud_cdm_cluster.test2
  ]  
}
output "az_filter_is_useful" {
  value = length(data.huaweicloud_cdm_clusters.az_filter.clusters) > 0 && alltrue(
    [for v in data.huaweicloud_cdm_clusters.az_filter.clusters[*].availability_zone : v == data.huaweicloud_availability_zones.test.names[0]]
  )  
}

data "huaweicloud_cdm_clusters" "status_filter" {
  status = "200"

  depends_on = [
    huaweicloud_cdm_cluster.test,
    huaweicloud_cdm_cluster.test2
  ]  
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_cdm_clusters.status_filter.clusters) > 0 && alltrue(
    [for v in data.huaweicloud_cdm_clusters.status_filter.clusters[*].status : v == "200"]
  )  
}

data "huaweicloud_cdm_clusters" "error_status_filter" {
  status = "300"

  depends_on = [
    huaweicloud_cdm_cluster.test,
    huaweicloud_cdm_cluster.test2
  ]  
}
output "error_status_filter_is_useful" {
  value = length(data.huaweicloud_cdm_clusters.error_status_filter.clusters) == 0
}
`, clustersConfig, name2, name)
}
