package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceInfos_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_resource_infos.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceInfos_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_service.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_service.0.service_num"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_service.0.resource_num"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_service.0.resource_distribution.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ccsp_service.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ccsp_service.0.instance_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "ccsp_service.0.instance_distribution.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vsm.#"),
					resource.TestCheckResourceAttrSet(dataSource, "app.#"),
					resource.TestCheckResourceAttrSet(dataSource, "app.0.app_num"),
					resource.TestCheckResourceAttrSet(dataSource, "kms.#"),
					resource.TestCheckResourceAttrSet(dataSource, "kms.0.total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "kms.0.result.#"),
				),
			},
		},
	})
}

const testDataSourceResourceInfos_basic = `data "huaweicloud_cpcs_resource_infos" "test" {}`
