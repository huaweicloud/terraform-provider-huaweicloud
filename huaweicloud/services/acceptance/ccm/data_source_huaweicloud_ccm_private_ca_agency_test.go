package ccm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePrivateCaAgency_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_ccm_private_ca_agency.test"
		dc         = acceptance.InitDataSourceCheck(datasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateCaAgency_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "agency_granted"),
				),
			},
		},
	})
}

const testAccDatasourcePrivateCaAgency_basic = `
data "huaweicloud_ccm_private_ca_agency" "test" {}
`
