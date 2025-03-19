package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAddons_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_addons.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCceAddons_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCceAddons_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_addons" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_addons.test.items) > 0
}
`, testAccAddon_basic(name))
}
