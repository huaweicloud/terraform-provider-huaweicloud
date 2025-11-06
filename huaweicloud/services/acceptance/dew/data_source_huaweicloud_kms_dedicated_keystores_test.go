package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDedicatedKeystores_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_kms_dedicated_keystores.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedKeystores_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "keystores.#"),
				),
			},
		},
	})
}

const testAccDataSourceDedicatedKeystores_basic string = `data "huaweicloud_kms_dedicated_keystores" "test" {}`
