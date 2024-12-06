package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppNasStorages_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_app_nas_storages.test"
		dc  = acceptance.InitDataSourceCheck(all)

		filterById   = "data.huaweicloud_workspace_app_nas_storages.filter_by_id"
		dcFilterById = acceptance.InitDataSourceCheck(filterById)

		filterByName   = "data.huaweicloud_workspace_app_nas_storages.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSfsFileSystemNames(t, 3)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppNasStorages_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "storages.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcFilterById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppNasStorages_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_nas_storage" "test" {
  count = length(slice(split(",", "%[1]s"), 0, 2))

  name = format("%[2]s_%%d", count.index)

  storage_metadata {
    storage_handle = element(slice(split(",", "%[1]s"), 0, 2), count.index)
    storage_class  = "sfs"
  }
}

resource "huaweicloud_workspace_app_nas_storage" "comparison" {
  count = length(slice(split(",", "%[1]s"), 2, 3))

  name = "storage_comparison"

  storage_metadata {
    storage_handle = element(split(",", "%[1]s"), 2)
    storage_class  = "sfs"
  }
}

data "huaweicloud_workspace_app_nas_storages" "test" {
  depends_on = [
    huaweicloud_workspace_app_nas_storage.test,
    huaweicloud_workspace_app_nas_storage.comparison,
  ]
}

// Filter by ID
locals {
  storage_id = element(huaweicloud_workspace_app_nas_storage.test[*].id, 0)
}

data "huaweicloud_workspace_app_nas_storages" "filter_by_id" {
  depends_on = [
    huaweicloud_workspace_app_nas_storage.test,
    huaweicloud_workspace_app_nas_storage.comparison,
  ]

  storage_id = local.storage_id
}

locals {
  id_filter_result = [for o in data.huaweicloud_workspace_app_nas_storages.filter_by_id.storages: o.id == local.storage_id]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) == 1 && alltrue(local.id_filter_result)
}

// Filter by name
locals {
  storage_name_prefix = "%[2]s"
}

data "huaweicloud_workspace_app_nas_storages" "filter_by_name" {
  depends_on = [
    huaweicloud_workspace_app_nas_storage.test,
    huaweicloud_workspace_app_nas_storage.comparison,
  ]

  name = local.storage_name_prefix
}

locals {
  name_filter_result = [for o in data.huaweicloud_workspace_app_nas_storages.filter_by_name.storages: strcontains(o.name, local.storage_name_prefix)]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) == 2 && alltrue(local.name_filter_result)
}
`, acceptance.HW_SFS_FILE_SYSTEM_NAMES, name)
}
