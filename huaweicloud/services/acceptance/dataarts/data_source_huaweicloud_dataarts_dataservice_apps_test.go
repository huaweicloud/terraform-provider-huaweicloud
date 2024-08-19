package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataServiceApps_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_dataservice_apps.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_dataservice_apps.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_dataarts_dataservice_apps.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataServiceApps_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "apps.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDataServiceApps_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_app" "test" {
  count = 2

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"

  name         = format("%[2]s_%%d", count.index)
  app_type     = "APP"
  description  = "Created by terraform script"
}

data "huaweicloud_dataarts_dataservice_apps" "test" {
  depends_on = [
    huaweicloud_dataarts_dataservice_app.test
  ]

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
}

# Filter by name
data "huaweicloud_dataarts_dataservice_apps" "filter_by_name" {
  depends_on = [huaweicloud_dataarts_dataservice_app.test]

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"

  name = "%[2]s" # Fuzzy search
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_apps.filter_by_name.apps[*].name : strcontains(v, "%[2]s")
  ]
}

output "is_name_filter_useful" {
  // At least two applications have the same name prefix.
  value = length(local.name_filter_result) >= 2 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_dataarts_dataservice_apps" "filter_by_not_found_name" {
  depends_on = [huaweicloud_dataarts_dataservice_app.test]

  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"

  name = local.not_found_name # This name is not exist
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_dataarts_dataservice_apps.filter_by_not_found_name.apps) == 0
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
