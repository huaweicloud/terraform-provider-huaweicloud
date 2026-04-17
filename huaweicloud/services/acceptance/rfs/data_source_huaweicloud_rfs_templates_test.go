package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsTemplates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_templates.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRfsTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.latest_version_id"),
				),
			},
		},
	})
}

func testAccDataSourceRfsTemplate_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_templates" "test" {
	depends_on = [huaweicloud_rfs_template.test]
}
`, testRfsTemplate_basic(name))
}
