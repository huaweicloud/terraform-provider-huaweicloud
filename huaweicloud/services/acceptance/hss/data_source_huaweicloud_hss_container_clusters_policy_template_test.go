package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerClustersPolicyTemplate_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_clusters_policy_template.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerClustersPolicyTemplate_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "target_kind"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "constraint_template"),
				),
			},
		},
	})
}

func testAccDataSourceContainerClustersPolicyTemplate_base() string {
	return `
data "huaweicloud_hss_container_clusters_policy_templates" "test" {
  enterprise_project_id = "0"
}
`
}

func testAccDataSourceContainerClustersPolicyTemplate_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_container_clusters_policy_template" "test" {
  policy_template_id    = data.huaweicloud_hss_container_clusters_policy_templates.test.data_list[0].id
  enterprise_project_id = "0"
}
`, testAccDataSourceContainerClustersPolicyTemplate_base())
}
