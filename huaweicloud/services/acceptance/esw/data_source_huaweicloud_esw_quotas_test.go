package esw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswQuotas_basic(t *testing.T) {
	rName := "data.huaweicloud_esw_quotas.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "quotas.#"),
					resource.TestCheckResourceAttrSet(rName, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(rName, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(rName, "quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(rName, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(rName, "quotas.0.resources.0.min"),
					resource.TestCheckResourceAttrSet(rName, "quotas.0.resources.0.max"),
				),
			},
		},
	})
}

func testAccDatasourceEswQuotas_basic() string {
	return `
data "huaweicloud_esw_quotas" "test" {}
`
}
