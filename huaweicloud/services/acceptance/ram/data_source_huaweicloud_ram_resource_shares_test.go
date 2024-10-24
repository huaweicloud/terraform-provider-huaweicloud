package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRamShares_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_ram_resource_shares.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_ram_resource_shares.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byPermissionID   = "data.huaweicloud_ram_resource_shares.filter_by_permission_id"
		dcByPermissionID = acceptance.InitDataSourceCheck(byPermissionID)

		byResourceShareIDs   = "data.huaweicloud_ram_resource_shares.filter_by_resource_share_ids"
		dcByResourceShareIDs = acceptance.InitDataSourceCheck(byResourceShareIDs)

		byStatus   = "data.huaweicloud_ram_resource_shares.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byTags   = "data.huaweicloud_ram_resource_shares.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMResourceShare(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRamShares_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.allow_external_principals"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.owning_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_shares.0.tags.%"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcByPermissionID.CheckResourceExists(),
					resource.TestCheckOutput("is_permission_id_filter_useful", "true"),

					dcByResourceShareIDs.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_share_ids_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRamShares_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ram_resource_shares" "test" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_owner = "self"
}

# Filter by name
locals {
  name = data.huaweicloud_ram_resource_shares.test.resource_shares[0].name
}

data "huaweicloud_ram_resource_shares" "filter_by_name" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_owner = "self"
  name           = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ram_resource_shares.filter_by_name.resource_shares[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by permission_id
data "huaweicloud_ram_resource_share_associated_permissions" "test" {
  resource_share_id = huaweicloud_ram_resource_share.test.id
}

locals {
  permission_id = data.huaweicloud_ram_resource_share_associated_permissions.test.associated_permissions[0].permission_id
}

data "huaweicloud_ram_resource_shares" "filter_by_permission_id" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_owner = "self"
  permission_id  = local.permission_id
}

output "is_permission_id_filter_useful" {
  value = length(data.huaweicloud_ram_resource_shares.filter_by_permission_id.resource_shares) > 0
}

# Filter by resource_share_ids
locals {
  resource_share_ids = split(",", data.huaweicloud_ram_resource_shares.test.resource_shares[0].id)
}

data "huaweicloud_ram_resource_shares" "filter_by_resource_share_ids" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_owner     = "self"
  resource_share_ids = local.resource_share_ids
}

locals {
  resource_share_ids_filter_result = [
    for v in data.huaweicloud_ram_resource_shares.filter_by_resource_share_ids.resource_shares[*].id : contains(local.resource_share_ids, v)
  ]
}

output "is_resource_share_ids_filter_useful" {
  value = length(local.resource_share_ids_filter_result) > 0 && alltrue(local.resource_share_ids_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_ram_resource_shares.test.resource_shares[0].status
}

data "huaweicloud_ram_resource_shares" "filter_by_status" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_owner = "self"
  status         = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ram_resource_shares.filter_by_status.resource_shares[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by tag_filters
locals {
  tag_value = data.huaweicloud_ram_resource_shares.test.resource_shares[0].tags["foo"]
}

data "huaweicloud_ram_resource_shares" "filter_by_tags" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_owner = "self"

  tag_filters {
    key    = "foo"
    values = [local.tag_value]
  }
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_ram_resource_shares.filter_by_tags.resource_shares[*].tags["foo"] : v == local.tag_value
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testRAMShare_basic(name))
}
