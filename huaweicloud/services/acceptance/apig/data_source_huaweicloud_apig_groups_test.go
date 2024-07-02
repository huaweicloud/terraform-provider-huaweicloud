package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGroups_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_apig_groups.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_apig_groups.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_groups.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`[1-9]\d*`)),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "groups.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "groups.0.id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(byId, "groups.0.name", "huaweicloud_apig_group.test", "name"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.status"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.sl_domain"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.updated_at"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.on_sell_status"),
					resource.TestCheckResourceAttr(byId, "groups.0.url_domains.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.url_domains.0.id"),
					resource.TestCheckResourceAttr(byId, "groups.0.url_domains.0.name", "www.terraform.test3.com"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.url_domains.0.cname_status"),
					resource.TestCheckResourceAttr(byId, "groups.0.url_domains.0.ssl_id", ""),
					resource.TestCheckResourceAttr(byId, "groups.0.url_domains.0.ssl_name", ""),
					resource.TestCheckResourceAttr(byId, "groups.0.url_domains.0.min_ssl_version", "TLSv1.1"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.url_domains.0.verified_client_certificate_enabled"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.url_domains.0.is_has_trusted_root_ca"),
					resource.TestCheckResourceAttr(byId, "groups.0.sl_domains.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.description"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.is_default"),
					resource.TestCheckResourceAttr(byId, "groups.0.environment.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "groups.0.environment.0.environment_id", "huaweicloud_apig_environment.test", "id"),
					resource.TestCheckResourceAttr(byId, "groups.0.environment.0.variable.#", "2"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.environment.0.variable.0.name"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.environment.0.variable.0.value"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.environment.0.variable.0.id"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGroups_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

variable "variables_configuration" {
  type = list(object({
    name  = string
    value = string
  }))
  default = [
    {name="TEST_VAR_1", value="TEST_VALUE_1"},
    {name="TEST_VAR_2", value="TEST_VALUE_2"},
  ]
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "huaweicloud_apig_group" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
  description = "Created by script"

  environment {
    environment_id = huaweicloud_apig_environment.test.id

    dynamic "variable" {
      for_each = var.variables_configuration

      content {
        name  = variable.value.name
        value = variable.value.value
      }
    }
  }

  url_domains {
    name                      = "www.terraform.test3.com"
    min_ssl_version           = "TLSv1.1"
    is_http_redirect_to_https = true
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceGroups_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_groups" "test" {
  depends_on = [huaweicloud_apig_group.test]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  group_id = huaweicloud_apig_group.test.id
}

data "huaweicloud_apig_groups" "filter_by_id" {
  instance_id = local.instance_id
  group_id    = local.group_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_groups.filter_by_id.groups[*].id : v == local.group_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  group_name = huaweicloud_apig_group.test.name
}

data "huaweicloud_apig_groups" "filter_by_name" {
  // The behavior of parameter 'name' is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_apig_group.test]

  instance_id = local.instance_id
  name        = local.group_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_groups.filter_by_name.groups[*].name : v == local.group_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name and the name is not exist
data "huaweicloud_apig_groups" "filter_by_not_found_name" {
  // Since a specified name is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [huaweicloud_apig_group.test]  

  instance_id = local.instance_id
  name        = "not_found_name"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_apig_groups.filter_by_not_found_name.groups) == 0
}
`, testAccDataSourceGroups_base())
}
