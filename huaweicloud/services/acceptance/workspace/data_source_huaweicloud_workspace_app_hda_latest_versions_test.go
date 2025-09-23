package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppHdaLatestVersions_basic(t *testing.T) {
	rName := "data.huaweicloud_workspace_app_hda_latest_versions.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppHdaLatestVersions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_latest_version_set", "true"),
					resource.TestCheckOutput("is_hda_type_set", "true"),
				),
			},
		},
	})
}

const testAccDataAppHdaLatestVersions_basic string = `
data "huaweicloud_workspace_app_hda_latest_versions" "test" {}

locals {
  hda_latest_versions      = data.huaweicloud_workspace_app_hda_latest_versions.test.hda_latest_versions
  first_hda_latest_version = try(local.hda_latest_versions[0], {})
}

output "is_latest_version_set" {
  value = length(local.hda_latest_versions) != 0 ? try(local.first_hda_latest_version.latest_version != "", false) : true 
}

output "is_hda_type_set" {
  value = length(local.hda_latest_versions) != 0 ? try(local.first_hda_latest_version.hda_type != "", false) : true 
}
`
