package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssociations_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSourcePrincipal = "data.huaweicloud_ram_resource_share_associations.test_principal"
		dcTestPrincipal     = acceptance.InitDataSourceCheck(dataSourcePrincipal)

		dataSourceResource = "data.huaweicloud_ram_resource_share_associations.test_resource"
		dcTestResource     = acceptance.InitDataSourceCheck(dataSourceResource)

		byPrincipal   = "data.huaweicloud_ram_resource_share_associations.filter_by_principal"
		dcByPrincipal = acceptance.InitDataSourceCheck(byPrincipal)

		byResourceShareIds   = "data.huaweicloud_ram_resource_share_associations.filter_by_resource_share_ids"
		dcByResourceShareIds = acceptance.InitDataSourceCheck(byResourceShareIds)

		byResourceUrn   = "data.huaweicloud_ram_resource_share_associations.filter_by_resource_urn"
		dcByResourceUrn = acceptance.InitDataSourceCheck(byResourceUrn)

		byStatus   = "data.huaweicloud_ram_resource_share_associations.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMResourceShare(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssociations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcTestPrincipal.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.#"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.associated_entity"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.association_type"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.external"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.resource_share_id"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.resource_share_name"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.status"),
					resource.TestCheckResourceAttrSet(dataSourcePrincipal, "associations.0.updated_at"),

					dcTestResource.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.#"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.associated_entity"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.association_type"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.external"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.resource_share_id"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.resource_share_name"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceResource, "associations.0.updated_at"),

					dcByPrincipal.CheckResourceExists(),
					resource.TestCheckOutput("is_principal_filter_useful", "true"),

					dcByResourceShareIds.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_share_ids_filter_useful", "true"),

					dcByResourceUrn.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_urn_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAssociations_basic(name string) string {
	return fmt.Sprintf(`
%s

# Filter by principal association type. Field association_type is required.
data "huaweicloud_ram_resource_share_associations" "test_principal" {
  depends_on = [huaweicloud_ram_resource_share.test]

  association_type = "principal"
}

# Filter by resource association type. Field association_type is required.
data "huaweicloud_ram_resource_share_associations" "test_resource" {
  depends_on = [huaweicloud_ram_resource_share.test]

  association_type = "resource"
}

# Filter by principal
locals {
  principal = data.huaweicloud_ram_resource_share_associations.test_principal.associations[0].associated_entity
}

data "huaweicloud_ram_resource_share_associations" "filter_by_principal" {
  depends_on = [huaweicloud_ram_resource_share.test]

  association_type = "principal"
  principal        = local.principal
}

locals {
  principal_filter_result = [
    for v in data.huaweicloud_ram_resource_share_associations.filter_by_principal.associations[*].associated_entity :
    v == local.principal
  ]
}

output "is_principal_filter_useful" {
  value = length(local.principal_filter_result) > 0 && alltrue(local.principal_filter_result)
}

# Filter by resource_share_ids
locals {
  resource_share_ids = split(",",
  data.huaweicloud_ram_resource_share_associations.test_principal.associations[0].resource_share_id)
}

data "huaweicloud_ram_resource_share_associations" "filter_by_resource_share_ids" {
  depends_on = [huaweicloud_ram_resource_share.test]

  association_type   = "principal"
  resource_share_ids = local.resource_share_ids
}

locals {
  resource_share_ids_filter_result = [
    for v in data.huaweicloud_ram_resource_share_associations.filter_by_resource_share_ids.associations[*].
    resource_share_id : contains(local.resource_share_ids, v)
  ]
}

output "is_resource_share_ids_filter_useful" {
  value = length(local.resource_share_ids_filter_result) > 0 && alltrue(local.resource_share_ids_filter_result)
}

# Filter by resource_urn
locals {
  resource_urn = data.huaweicloud_ram_resource_share_associations.test_resource.associations[0].associated_entity
}

data "huaweicloud_ram_resource_share_associations" "filter_by_resource_urn" {
  depends_on = [huaweicloud_ram_resource_share.test]

  association_type = "resource"
  resource_urn     = local.resource_urn
}

locals {
  resource_urn_filter_result = [
    for v in data.huaweicloud_ram_resource_share_associations.filter_by_resource_urn.associations[*].associated_entity :
    v == local.resource_urn
  ]
}

output "is_resource_urn_filter_useful" {
  value = length(local.resource_urn_filter_result) > 0 && alltrue(local.resource_urn_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_ram_resource_share_associations.test_resource.associations[0].status
}

data "huaweicloud_ram_resource_share_associations" "filter_by_status" {
  depends_on = [huaweicloud_ram_resource_share.test]

  association_type = "resource"
  status           = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ram_resource_share_associations.filter_by_status.associations[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testRAMShare_basic(name))
}
