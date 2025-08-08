package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplicationNextModels_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_next_models.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocApplicationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplicationNextModels_applicationID_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.parent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_applications.0.update_time"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceCocApplicationNextModels_componentID_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_next_models.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocComponentID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplicationNextModels_componentID_basic(),
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
					resource.TestCheckOutput("component_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplicationNextModels_applicationID_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_application_next_models" "test" {
  application_id = "%s"
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_next_models.test.sub_applications) > 0
}
`, acceptance.HW_COC_APPLICATION_ID)
}

func testDataSourceDataSourceCocApplicationNextModels_componentID_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_application_next_models" "test" {
  component_id = "%s"
}

output "component_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_next_models.test.groups) > 0
}
`, acceptance.HW_COC_COMPONENT_ID)
}
