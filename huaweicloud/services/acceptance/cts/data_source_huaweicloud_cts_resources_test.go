package cts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCtsResources_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cts_resources.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCtsResources_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource.#"),
				),
			},
		},
	})
}

const testDataSourceCtsResources_basic = `data "huaweicloud_cts_resources" "test" {}`
