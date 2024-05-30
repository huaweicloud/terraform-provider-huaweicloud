package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssClusterTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_cluster_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssClusterTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "tags.key", "value"),
					resource.TestCheckResourceAttr(dataSource, "tags.foo", "bar"),
				),
			},
		},
	})
}

func testDataSourceCssClusterTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_css_cluster_tags" "test" {
  depends_on = [huaweicloud_css_cluster.test]

  resource_type = "css-cluster"
}
`, testAccCssCluster_basic(name, "Test@passw0rd", 7, "bar"))
}
