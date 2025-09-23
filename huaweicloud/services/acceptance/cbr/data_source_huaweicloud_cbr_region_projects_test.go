package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCbrRegionProjects_basic(t *testing.T) {
	var (
		dc = acceptance.InitDataSourceCheck("data.huaweicloud_cbr_region_projects.test")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCbrRegionProjects_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "region"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "links.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.domain_id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.is_domain"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.parent_id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.name"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.enabled"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.links.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cbr_region_projects.test", "projects.0.links.0.self"),
				),
			},
		},
	})
}

const testAccDataSourceCbrRegionProjects_basic = `data "huaweicloud_cbr_region_projects" "test" {}`
