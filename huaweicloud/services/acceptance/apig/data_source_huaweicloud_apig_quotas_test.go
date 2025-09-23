package apig

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_apig_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQuotas,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "quotas.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.created_at"),
				),
			},
		},
	})
}

const testAccDataSourceQuotas = `data "huaweicloud_apig_quotas" "test" {}`
