package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceOmUrl_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cbh_instance_om_url.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a CBH instance ID and config it to the environment variable.
			acceptance.TestAccPreCheckCbhServerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstanceOmUrl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "state", "HOST_NOT_MANAGE"),
					resource.TestCheckResourceAttr(dataSource, "description", "主机未被CBH纳管。"),
				),
			},
		},
	})
}

// The IP address and account name of the managed host are fixed values.
func testDataSourceInstanceOmUrl_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbh_instance_om_url" "test" {
  server_id         = "%s"
  ip_address        = "192.168.0.7"
  host_account_name = "root"
}
`, acceptance.HW_CBH_SERVER_ID)
}
