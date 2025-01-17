package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.enterprise_project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.instance_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.vcpus_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.ram_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.volume_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.instance_used"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.vcpus_used"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.ram_used"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_quotas.0.volume_used"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussQuotas_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_quotas" "test" {}

locals {
  enterprise_project_id = "%[1]s"
}
data "huaweicloud_gaussdb_opengauss_quotas" "enterprise_project_id_filter" {
  enterprise_project_id = "%[1]s"
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_quotas.enterprise_project_id_filter.eps_quotas) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_quotas.enterprise_project_id_filter.eps_quotas[*].enterprise_project_id :
  v == local.enterprise_project_id]
  )
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
