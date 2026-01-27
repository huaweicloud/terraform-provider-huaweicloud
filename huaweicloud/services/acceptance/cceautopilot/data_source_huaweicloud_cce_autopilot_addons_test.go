package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAddons_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_addons.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAddons_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.kind"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.api_version"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.0.alias"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.metadata.0.uid"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.addon_template_labels.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.addon_template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.addon_template_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.values.%"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.spec.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.reason"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.input.%"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.stable"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.support_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.translate.%"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.update_timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.current_version.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.is_rollbackable"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.reason"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status.0.target_versions.#"),
					resource.TestCheckOutput("addon_template_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAddons_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_autopilot_addons" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_cce_autopilot_addons" "addon_template_name_filter" {
  cluster_id          = "%[1]s"
  addon_template_name = data.huaweicloud_cce_autopilot_addons.test.items[0].metadata[0].name
}
locals {
  name = data.huaweicloud_cce_autopilot_addons.test.items[0].metadata[0].name
}
output "addon_template_name_filter_is_useful" {
  value = length(data.huaweicloud_cce_autopilot_addons.addon_template_name_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_cce_autopilot_addons.addon_template_name_filter.items[*].metadata : alltrue(
      [for vv in v[*].name : vv == local.name]
    )]
  )
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
