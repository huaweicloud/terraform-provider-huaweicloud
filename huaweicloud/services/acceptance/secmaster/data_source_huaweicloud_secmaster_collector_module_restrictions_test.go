package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCollectorModuleRestrictions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_collector_module_restrictions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterTemplateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSocCollectorModuleRestrictions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.#"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.default_value"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.example"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.required"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.restrictions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.restrictions.0.logic"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.restrictions.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.restrictions.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.restrictions.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.template_field_id"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "module_restrictions.0.fields.0.type"),
				),
			},
		},
	})
}

func testAccDataSourceSocCollectorModuleRestrictions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collector_module_restrictions" "test" {
  workspace_id = "%[1]s"
  template_ids = ["%[2]s"]
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_TEMPLATE_ID)
}
