package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataEventRouterClusters_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_eg_eventrouter_clusters.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_eg_eventrouter_clusters.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEventRouterClusters_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "clusters.#", regexp.MustCompile(`^[0-9]+$`)),

					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataEventRouterClusters_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_eg_eventrouter_cluster" "test" {
  count = 2

  name               = format("%[2]s-%%d", count.index)
  source_type        = "KAFKA"
  sink_type          = "KAFKA"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  description        = "Created by terraform script"
  availability_zones = join(",", slice(data.huaweicloud_availability_zones.test.names, 0, 2))
  flavor             = "small"
}
`, common.TestVpc(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name)
}

func testAccDataEventRouterClusters_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_eg_eventrouter_clusters" "all" {
  depends_on = [huaweicloud_eg_eventrouter_cluster.test]
}

# Filter by 'name' parameter.
locals {
  cluster_name_keyword = "%[2]s"
}

data "huaweicloud_eg_eventrouter_clusters" "filter_by_name" {
  depends_on = [huaweicloud_eg_eventrouter_cluster.test]

  name = local.cluster_name_keyword
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_eg_eventrouter_clusters.filter_by_name.clusters[*].name :
      strcontains(v, local.cluster_name_keyword)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, testAccDataEventRouterClusters_basic_base(name), name)
}
