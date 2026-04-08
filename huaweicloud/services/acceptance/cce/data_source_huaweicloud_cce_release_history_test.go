package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEReleaseHistory_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_release_history.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceReleaseHistory_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "releases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.chart_name"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.chart_public"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.chart_version"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.status_description"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.update_at"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.values"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.version"),
				),
			},
		},
	})
}

func testAccDataSourceReleaseHistory_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_releases" "namespace_filter" {
  cluster_id = "%[1]s"
  namespace  = "default"
}

locals {
  release_name = data.huaweicloud_cce_releases.namespace_filter.releases[0].name
}

data "huaweicloud_cce_release_history" "test" {
  cluster_id = "%[1]s"
  namespace  = "default"
  name       = local.release_name
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
