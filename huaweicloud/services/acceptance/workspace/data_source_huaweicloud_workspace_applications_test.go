package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplications_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_workspace_applications.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByName   = "data.huaweicloud_workspace_applications.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.name"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.version"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.description"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.authorization_type"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.install_type"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.support_os"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.catalog_id"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.application_icon_url"),
					resource.TestCheckOutput("is_install_command_set_and_valid", "true"),
					resource.TestCheckOutput("is_uninstall_command_set_and_valid", "true"),
					resource.TestCheckOutput("is_install_info_set_and_valid", "true"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.status"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.application_source"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.create_time"),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.catalog"),
					resource.TestMatchResourceAttr(dcName, "applications.0.application_file_store.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "applications.0.application_file_store.0.store_type"),
					dcFilterByName.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByName, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataApplications_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_application_catalogs" "test" {}

resource "huaweicloud_workspace_application" "with_file_store" {
  name               = "%[1]s_with_file_store"
  version            = "1.0.0"
  description        = "Created by terraform script"
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Windows"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[0].id, "NOT_FOUND")
  install_command    = "terraform test install"
  uninstall_command  = "terraform test uninstall"
  install_info       = "{\"user\":\"Terraform\"}"

  application_file_store {
    store_type = "LINK"
    file_link  = "https://www.huaweicloud.com/TerraformTest.msi"
  }
}

resource "huaweicloud_workspace_application" "with_obs_store" {
  name               = "%[1]s_with_obs_store"
  version            = "1.0.0"
  description        = "created by terraform script."
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Linux"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[1].id, "NOT_FOUND")
  install_command    = "terraform test install"
  uninstall_command  = "terraform test uninstall"
  install_info       = "{\"user\":\"Terraform\"}"

  application_file_store {
    store_type = "OBS"

    bucket_store {
      bucket_name      = "app-center-%[2]s"
      bucket_file_path = "dir1/TerraformTest.apk"
    }
  }
}`, name, acceptance.HW_PROJECT_ID)
}

func testAccDataApplications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_applications" "all" {
  depends_on = [
    huaweicloud_workspace_application.with_file_store,
    huaweicloud_workspace_application.with_obs_store,
  ]
}

locals {
  component_id = huaweicloud_workspace_application.with_file_store.id
  manully_filter_component_result = try([
    for v in data.huaweicloud_workspace_applications.all.applications : v if v.id == local.component_id
  ][0], null)
}

output "is_install_command_set_and_valid" {
  value = try(local.manully_filter_component_result.install_command == huaweicloud_workspace_application.with_file_store.install_command
  , false)
}

output "is_uninstall_command_set_and_valid" {
  value = try(local.manully_filter_component_result.uninstall_command == huaweicloud_workspace_application.with_file_store.uninstall_command
  , false)
}

output "is_install_info_set_and_valid" {
  value = try(local.manully_filter_component_result.install_info == huaweicloud_workspace_application.with_file_store.install_info
  , false)
}

# Filter by name
data "huaweicloud_workspace_applications" "filter_by_name" {
  name = "%[2]s"

  depends_on = [
    huaweicloud_workspace_application.with_file_store,
    huaweicloud_workspace_application.with_obs_store,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_workspace_applications.filter_by_name.applications) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_applications.filter_by_name.applications[*].name : strcontains(v, "%[2]s")]
  )
}`, testAccDataApplications_base(name), name)
}
