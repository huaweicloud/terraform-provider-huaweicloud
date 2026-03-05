package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcepserviceConnections_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_vpcep_connections.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcepserviceConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "connections.#"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.max_session"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.specification_name"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.vpcep_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.permission"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.permission_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.created_at"),
				),
			},
		},
	})
}

func testDataSourceVpcepserviceConnections_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_vpcep_connections" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
