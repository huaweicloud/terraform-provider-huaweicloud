package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case cannot be executed successfully.
func TestAccDataSourceConnectionNumbers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_connection_numbers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConnectionNumbers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

// Parameter `ip` are mock data.
func testDataSourceConnectionNumbers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_connection_numbers" "test" {
  instance_id = "%[1]s"
  start_time  = "1755734400"
  end_time    = "1755820800"
  ip          = "192.168.1.1"
}
`, acceptance.HW_AAD_INSTANCE_ID)
}
