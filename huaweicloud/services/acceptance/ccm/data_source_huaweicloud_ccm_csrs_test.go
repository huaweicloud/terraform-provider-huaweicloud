package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcmCsrs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ccm_csrs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcmCsrs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.private_key_algo"),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.usage"),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "csr_list.0.update_time"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("private_key_algo_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcmCsrs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_csrs" "test" {
  depends_on = [huaweicloud_ccm_csr.test]
}

# Filter by name
locals {
  name = data.huaweicloud_ccm_csrs.test.csr_list[0].name
}

data "huaweicloud_ccm_csrs" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ccm_csrs.filter_by_name.csr_list[*].name : strcontains(v, local.name)
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by private_key_algo
locals {
  private_key_algo = data.huaweicloud_ccm_csrs.test.csr_list[0].private_key_algo
}

data "huaweicloud_ccm_csrs" "filter_by_private_key_algo" {
  private_key_algo = local.private_key_algo
}

locals {
  private_key_algo_filter_result = [
    for v in data.huaweicloud_ccm_csrs.filter_by_private_key_algo.csr_list[*].private_key_algo : v == local.private_key_algo
  ]
}

output "private_key_algo_filter_is_useful" {
  value = alltrue(local.private_key_algo_filter_result) && length(local.private_key_algo_filter_result) > 0
}
`, testCCMCsr_basic(name))
}
