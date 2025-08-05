package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplicationNextModels_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_next_models.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplicationNextModels_applicationID_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "components.#"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.application_id"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.domain_id"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceCocApplicationNextModels_componentID_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_next_models.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplicationNextModels_componentID_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.application_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.sync_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.vendor"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.sync_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.sync_rules.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.sync_rules.0.rule_tags"),
					resource.TestCheckOutput("component_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplicationNextModels_applicationID_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_application_next_models" "test" {
  application_id = huaweicloud_coc_application.test.id

  depends_on = [huaweicloud_coc_component.test]
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_next_models.test.components) > 0 && alltrue(
    [for v in data.huaweicloud_coc_application_next_models.test.components[*].application_id :
      v == huaweicloud_coc_application.test.id]
  )
}
`, testAccComponent_basic(name))
}

func testDataSourceDataSourceCocApplicationNextModels_componentID_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_application_next_models" "test" {
  component_id = huaweicloud_coc_component.test.id

  depends_on = [huaweicloud_coc_group.test]
}

output "component_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_next_models.test.groups) > 0 && alltrue(
    [for v in data.huaweicloud_coc_application_next_models.test.groups[*].component_id :
      v == huaweicloud_coc_component.test.id]
  )
}
`, testAccGroup_basic(name))
}
