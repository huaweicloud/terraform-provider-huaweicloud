package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcConnectionTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_connection_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcConnectionTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.foo"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.key"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcConnectionTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_connection_tags" "test" {
  depends_on = [
    huaweicloud_cc_connection.test1,
    huaweicloud_cc_connection.test2,
    huaweicloud_cc_connection.test3,
  ]
}
`, testAccDatasourceCreateCloudConnections(name))
}
