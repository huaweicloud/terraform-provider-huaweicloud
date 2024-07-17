package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCustomAuthorizers_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_custom_authorizers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		rName      = acceptance.RandomAccResourceName()

		byId   = "data.huaweicloud_apig_custom_authorizers.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_custom_authorizers.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_apig_custom_authorizers.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomAuthorizers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "authorizers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "authorizers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizers.0.type"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("authorizer_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCustomAuthorizers_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "FRONTEND"
  is_body_send     = true
  cache_age        = 60
  
  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
`, testAccCustomAuthorizer_base(name), name)
}

func testAccDataSourceCustomAuthorizers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_custom_authorizers" "test" {
  depends_on = [
    huaweicloud_apig_custom_authorizer.test
  ]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  authorizer_id = huaweicloud_apig_custom_authorizer.test.id
}

data "huaweicloud_apig_custom_authorizers" "filter_by_id" {
  instance_id  = local.instance_id
  authorizer_id = local.authorizer_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_custom_authorizers.filter_by_id.authorizers[*].id : v == local.authorizer_id
  ]
}

output "authorizer_id_filter_is_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  name = huaweicloud_apig_custom_authorizer.test.name
}

data "huaweicloud_apig_custom_authorizers" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_custom_authorizer.test
  ]

  instance_id = local.instance_id
  name        = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_custom_authorizers.filter_by_name.authorizers[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by type
locals {
  type = huaweicloud_apig_custom_authorizer.test.type
}

data "huaweicloud_apig_custom_authorizers" "filter_by_type" {
  depends_on = [
    huaweicloud_apig_custom_authorizer.test
  ]

  instance_id = local.instance_id
  type        = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_custom_authorizers.filter_by_type.authorizers[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, testAccDataSourceCustomAuthorizers_base(name))
}
