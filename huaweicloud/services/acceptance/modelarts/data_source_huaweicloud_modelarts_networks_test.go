package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNetworks_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_networks.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byLabelSelector   = "data.huaweicloud_modelarts_networks.filter_by_label_selector"
		dcByLabelSelector = acceptance.InitDataSourceCheck(byLabelSelector)

		rName = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNetworks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "networks.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "networks.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(all, "networks.0.metadata.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "networks.0.metadata.0.labels.0.os_modelarts_name"),
					resource.TestCheckResourceAttrSet(all, "networks.0.metadata.0.labels.0.os_modelarts_workspace_id"),
					resource.TestCheckResourceAttrSet(all, "networks.0.spec.0.cidr"),
					resource.TestCheckResourceAttrSet(all, "networks.0.status.0.phase"),

					dcByLabelSelector.CheckResourceExists(),
					resource.TestCheckOutput("is_label_selector_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataNetworks_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_network" "test" {
  count = 2
 
  name = "%[1]s-${count.index}"
  cidr = "192.168.20.0/24"

}

# Without any filter parameters.
data "huaweicloud_modelarts_networks" "all" {
  depends_on = [
    huaweicloud_modelarts_network.test
  ]
}

# Filter by label selector.
locals {
  resource_pool_name = try(data.huaweicloud_modelarts_networks.all.networks[0].metadata[0].labels[0].os_modelarts_name, "")
  label_selector     = "os.modelarts/name=${local.resource_pool_name}"
}

data "huaweicloud_modelarts_networks" "filter_by_label_selector" {
  label_selector = local.label_selector

  depends_on = [
    huaweicloud_modelarts_network.test
  ]
}

locals {
  label_selector_filter_result = [
    for v in data.huaweicloud_modelarts_networks.filter_by_label_selector.networks :
    try(v.metadata[0].labels[0].os_modelarts_name, "") == local.resource_pool_name
  ]
}

output "is_label_selector_filter_useful" {
  value = length(local.label_selector_filter_result) > 0 && alltrue(local.label_selector_filter_result)
}`, name)
}
