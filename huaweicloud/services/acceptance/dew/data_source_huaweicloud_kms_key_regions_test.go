package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKMSKeyRegionsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_kms_key_regions.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testKMSKeyRegions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "regions.#"),
				),
			},
		},
	})
}

const testKMSKeyRegions_basic = `data "huaweicloud_kms_key_regions" "test" {}`
