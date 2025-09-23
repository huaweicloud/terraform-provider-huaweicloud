package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceLoginUrl_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cbh_instance_login_url.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a CBH instance ID and config it to the environment variable.
			acceptance.TestAccPreCheckCbhInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstanceLoginUrl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "login_url"),
				),
			},
		},
	})
}

func testDataSourceInstanceLoginUrl_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbh_instance_login_url" "test" {
  server_id = "%s"
}
`, acceptance.HW_CBH_INSTANCE_ID)
}
