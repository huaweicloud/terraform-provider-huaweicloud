package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssImageRegistryStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_image_registry_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssImageRegistryStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "fail_num"),
					resource.TestCheckResourceAttrSet(dataSource, "success_num"),
				),
			},
		},
	})
}

const testDataSourceDataSourceHssImageRegistryStatistics_basic = `
data "huaweicloud_hss_image_registry_statistics" "test" {
  enterprise_project_id = "all_granted_eps"
}`
