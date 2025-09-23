package aom

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCloudServiceAuthorizations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_cloud_service_authorizations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCloudServiceAuthorizations_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.service"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.role_name.#"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.status"),
				),
			},
		},
	})
}

const testDataSourceCloudServiceAuthorizations_basic = `data "huaweicloud_aom_cloud_service_authorizations" "test" {}`
