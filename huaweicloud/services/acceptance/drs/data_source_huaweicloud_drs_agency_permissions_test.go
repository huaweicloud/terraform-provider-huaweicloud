package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsAgencyPermissions_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_agency_permissions.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgencyPermissions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "common_permissions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "engine_permissions.#"),
				),
			},
		},
	})
}

const testAccAgencyPermissions_basic = `
data "huaweicloud_drs_agency_permissions" "test" {
  is_non_dbs = false
}
`
