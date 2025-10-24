package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to limitations in the testing environment, the test case was not actually executed successfully.
func TestAccDataSourceContainerIacFileRisks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_iac_file_risks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSIACFileId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerIacFileRisks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceContainerIacFileRisks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_container_iac_file_risks" "test" {
  file_id = "%s"
}
`, acceptance.HW_HSS_IAC_FILE_ID)
}
