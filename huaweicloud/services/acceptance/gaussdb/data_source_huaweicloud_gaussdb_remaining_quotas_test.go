package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbRemainingQuotas_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_remaining_quotas.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbRemainingQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.eps_tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.instance_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.cpu_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.mem_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.volume_eps_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.instance_eps_remaining_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.cpu_eps_remaining_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.mem_eps_remaining_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eps_remaining_quotas.0.volume_eps_remaining_quota"),
				),
			},
		},
	})
}

func testAccGaussDbRemainingQuotas_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_remaining_quotas" "test" {
  eps_tags = ["%s"]
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID)
}
