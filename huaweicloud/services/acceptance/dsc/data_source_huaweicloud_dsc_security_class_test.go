package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityClass_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_security_class.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDSCScanTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecurityClass_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "template_category"),
					resource.TestCheckResourceAttrSet(dataSource, "template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.#"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children_nums"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.depth"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children.#"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children.0.children_nums"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children.0.depth"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_trees.0.children.0.name"),
				),
			},
		},
	})
}

func testDataSourceSecurityClass_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_security_class" "test" {
  template_id = "%[1]s"
}
`, acceptance.HW_DSC_SCAN_TEMPLATE_ID)
}
