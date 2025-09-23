package aom

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAomAccessCodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_access_codes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAomAccessCodes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "access_codes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_codes.0.access_code_id"),
					resource.TestCheckResourceAttrSet(dataSource, "access_codes.0.access_code"),
					resource.TestCheckResourceAttrSet(dataSource, "access_codes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "access_codes.0.create_at"),
				),
			},
		},
	})
}

const testDataSourceDataSourceAomAccessCodes_basic = `data "huaweicloud_aom_access_codes" "test" {}`
