package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssClusterResourceTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_cluster_resource_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssClusterResourceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "tags.key", "value"),
					resource.TestCheckResourceAttr(dataSource, "tags.foo", "bar"),
				),
			},
		},
	})
}

func testDataSourceCssClusterResourceTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_cluster_resource_tags" "test" {
  resource_type = "css-cluster"
  cluster_id    = "%s"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
