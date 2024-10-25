package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRAMSharedResources_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_ram_shared_resources.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byPrincipal   = "data.huaweicloud_ram_shared_resources.principal_filter"
		dcByPrincipal = acceptance.InitDataSourceCheck(byPrincipal)

		byResourceUrns   = "data.huaweicloud_ram_shared_resources.resource_urns_filter"
		dcByResourceUrns = acceptance.InitDataSourceCheck(byResourceUrns)

		byResourceShareIDs   = "data.huaweicloud_ram_shared_resources.resource_share_ids_filter"
		dcByResourceShareIDs = acceptance.InitDataSourceCheck(byResourceShareIDs)

		byResourceRegion   = "data.huaweicloud_ram_shared_resources.resource_region_filter"
		dcByResourceRegion = acceptance.InitDataSourceCheck(byResourceRegion)

		byResourceType   = "data.huaweicloud_ram_shared_resources.resource_type_filter"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMShareId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceRAMSharedReources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.resource_urn"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.resource_type"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.resource_share_id"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.status"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.updated_at"),

					dcByPrincipal.CheckResourceExists(),
					resource.TestCheckOutput("is_principal_filter_useful", "true"),

					dcByResourceUrns.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_urns_filter_useful", "true"),

					dcByResourceShareIDs.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_share_ids_filter_useful", "true"),

					dcByResourceRegion.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_region_filter_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceRAMSharedReources_base() string {
	return fmt.Sprintf(`
data "huaweicloud_ram_resource_share_associations" "test_principal" {
  resource_share_ids = ["%[1]s"]
  association_type   = "principal"
}

data "huaweicloud_ram_resource_share_associations" "test_resource" {
  resource_share_ids = ["%[1]s"]
  association_type   = "resource"
}

locals {
  principal    = data.huaweicloud_ram_resource_share_associations.test_principal.associations[0].associated_entity
  resource_urn = data.huaweicloud_ram_resource_share_associations.test_resource.associations[0].associated_entity
}
`, acceptance.HW_RAM_SHARE_ID)
}

func testAccDatasourceRAMSharedReources_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ram_shared_resources" "test" {
  resource_owner = "self"
}

// Filter by principal
data "huaweicloud_ram_shared_resources" "principal_filter" {
  resource_owner = "self"
  principal      = local.principal
}

output "is_principal_filter_useful" {
  value = length(data.huaweicloud_ram_shared_resources.principal_filter.shared_resources) > 0
}

// Filter by resource_urns
data "huaweicloud_ram_shared_resources" "resource_urns_filter" {
  resource_owner = "self"
  resource_urns  = [local.resource_urn]
}

locals {
  resource_urns_filter_result = [
    for v in data.huaweicloud_ram_shared_resources.resource_urns_filter.shared_resources[*].resource_urn :
    v == local.resource_urn
  ]
}

output "is_resource_urns_filter_useful" {
  value = alltrue(local.resource_urns_filter_result) && length(local.resource_urns_filter_result) > 0
}

# Filter by resource_share_ids
locals {
  resource_share_ids = split(",", data.huaweicloud_ram_shared_resources.test.shared_resources[0].resource_share_id)
}

data "huaweicloud_ram_shared_resources" "resource_share_ids_filter" {
  resource_owner     = "self"
  resource_share_ids = local.resource_share_ids
}

locals {
  resource_share_ids_filter_result = [
    for v in data.huaweicloud_ram_shared_resources.resource_share_ids_filter.shared_resources[*].resource_share_id :
    contains(local.resource_share_ids, v)
  ]
}

output "is_resource_share_ids_filter_useful" {
  value = alltrue(local.resource_share_ids_filter_result) && length(local.resource_share_ids_filter_result) > 0
}

# Filter by resource_region
data "huaweicloud_ram_shared_resources" "resource_region_filter" {
  resource_owner  = "self"
  resource_region = "%s"
}

output "is_resource_region_filter_useful" {
  value = length(data.huaweicloud_ram_shared_resources.resource_region_filter.shared_resources) > 0
}

# Filter by resource_type
locals {
  resource_type = data.huaweicloud_ram_shared_resources.test.shared_resources[0].resource_type
}

data "huaweicloud_ram_shared_resources" "resource_type_filter" {
  resource_owner = "self"
  resource_type  = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.huaweicloud_ram_shared_resources.resource_type_filter.shared_resources[*].resource_type :
    v == local.resource_type
  ]
}

output "is_resource_type_filter_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}
`, testAccDatasourceRAMSharedReources_base(), acceptance.HW_REGION_NAME)
}
