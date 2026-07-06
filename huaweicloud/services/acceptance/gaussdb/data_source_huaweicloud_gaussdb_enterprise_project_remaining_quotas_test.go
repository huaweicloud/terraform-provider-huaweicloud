package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBEnterpriseProjectRemainingQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_enterprise_project_remaining_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBEnterpriseProjectRemainingQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "eps_tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.eps_tag"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.instance_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.cpu_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.mem_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.volume_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.instance_eps_remaining_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.cpu_eps_remaining_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.mem_eps_remaining_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eps_remaining_quotas.0.volume_eps_remaining_quota"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBEnterpriseProjectRemainingQuotas_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_enterprise_project_remaining_quotas" "test" {
  eps_tags = ["%s"]
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
