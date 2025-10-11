package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineSecurityChecksDirectories_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_security_checks_directories.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBaselineSecurityChecksDirectories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "task_condition.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baseline_directory_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baseline_directory_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baseline_directory_list.0.standard"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baseline_directory_list.0.data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pwd_directory_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pwd_directory_list.0.tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pwd_directory_list.0.sub_tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pwd_directory_list.0.checked"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pwd_directory_list.0.key"),
				),
			},
		},
	})
}

func testAccDataSourceBaselineSecurityChecksDirectories_basic() string {
	return `
data "huaweicloud_hss_baseline_security_checks_directories" "test" {
  support_os  = "Linux"
  select_type = "check_type"
}
`
}
