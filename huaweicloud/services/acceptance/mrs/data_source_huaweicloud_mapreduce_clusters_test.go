package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceMrsClusters_basic(t *testing.T) {
	rName := "data.huaweicloud_mapreduce_clusters.name_filter"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	pwd := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMrsClusters_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.name", "huaweicloud_mapreduce_cluster.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.enterprise_project_id",
						"huaweicloud_mapreduce_cluster.test", "enterprise_project_id"),
					resource.TestCheckResourceAttr(rName, "clusters.0.type", "1"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.version",
						"huaweicloud_mapreduce_cluster.test", "version"),
					resource.TestCheckResourceAttr(rName, "clusters.0.safe_mode", "1"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vnc"),
					resource.TestCheckResourceAttr(rName, "clusters.0.component_list.0.component_name", "Storm"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceMrsClusters_basic(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_mapreduce_clusters" "status_filter" {
  status = "running"

  depends_on = [
    huaweicloud_mapreduce_cluster.test
  ]
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_mapreduce_clusters.status_filter.clusters) > 0 && alltrue(
    [for v in data.huaweicloud_mapreduce_clusters.status_filter.clusters[*].status : v == "running"]
  )  
}

data "huaweicloud_mapreduce_clusters" "name_filter" {
  name   = "%[2]s"
  status = "existing"

  depends_on = [
    huaweicloud_mapreduce_cluster.test
  ]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_mapreduce_clusters.name_filter.clusters) > 0 && alltrue(
    [for v in data.huaweicloud_mapreduce_clusters.name_filter.clusters[*].name : v == "%[2]s"]
  )  
}

data "huaweicloud_mapreduce_clusters" "tags_filter" {
  tags = "foo*bar,key*value"

  depends_on = [
    huaweicloud_mapreduce_cluster.test
  ]
}
output "tags_filter_is_useful" {
  value = length(data.huaweicloud_mapreduce_clusters.tags_filter.clusters) > 0 && alltrue(
    [for v in data.huaweicloud_mapreduce_clusters.tags_filter.clusters[*].tags.foo : v == "bar"]
  )  
}
`, testAccMrsMapReduceClusterConfig_basic(name, pwd), name)
}
