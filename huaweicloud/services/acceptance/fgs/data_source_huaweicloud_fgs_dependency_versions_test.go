package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDependencyVersions_basic(t *testing.T) {
	var (
		rName          = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_fgs_dependency_versions.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byVersionId   = "data.huaweicloud_fgs_dependency_versions.filter_by_version_id"
		dcByVersionId = acceptance.InitDataSourceCheck(byVersionId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsDependencyLink(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDependencyVersions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.#"),
					dcByVersionId.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.id", "huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.version", "huaweicloud_fgs_dependency_version.test", "version"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.dependency_id", "huaweicloud_fgs_dependency_version.test",
						"dependency_id"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.dependency_name", "huaweicloud_fgs_dependency_version.test", "name"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.runtime", "huaweicloud_fgs_dependency_version.test", "runtime"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.link", "huaweicloud_fgs_dependency_version.test", "link"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.size", "huaweicloud_fgs_dependency_version.test", "size"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.etag", "huaweicloud_fgs_dependency_version.test", "etag"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.owner", "huaweicloud_fgs_dependency_version.test", "owner"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.description", "huaweicloud_fgs_dependency_version.test",
						"description"),
					resource.TestCheckOutput("is_version_id_filter_useful", "true"),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
				),
			},
		},
	})
}
func testAccDependencyVersions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_dependency_versions" "test" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
}

// Filter by version ID
locals {
  version_id = huaweicloud_fgs_dependency_version.test.version_id
}

data "huaweicloud_fgs_dependency_versions" "filter_by_version_id" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  version_id    = local.version_id
}

output "is_version_id_filter_useful" {
  value = length(data.huaweicloud_fgs_dependency_versions.filter_by_version_id.versions) >= 1 && alltrue(
    [for v in data.huaweicloud_fgs_dependency_versions.filter_by_version_id.versions[*].id : v == local.version_id]
  )
}

// Filter by version
locals {
  version = huaweicloud_fgs_dependency_version.test.version
}

data "huaweicloud_fgs_dependency_versions" "filter_by_version" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  version       = local.version
}

output "is_version_filter_useful" {
  value = length(data.huaweicloud_fgs_dependency_versions.filter_by_version.versions) > 0 && alltrue(
    [for v in data.huaweicloud_fgs_dependency_versions.filter_by_version.versions[*].version : v == local.version]
  )
}

// Filter by runtime
locals {
  runtime = data.huaweicloud_fgs_dependency_versions.test.versions[0].runtime
}

data "huaweicloud_fgs_dependency_versions" "filter_by_runtime" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  runtime       = local.runtime
}

output "is_runtime_filter_useful" {
  value = length(data.huaweicloud_fgs_dependency_versions.filter_by_runtime.versions) > 0 && alltrue(
    [for v in data.huaweicloud_fgs_dependency_versions.filter_by_runtime.versions[*].runtime : v == local.runtime]
  )
}
`, testAccDependencyVersion_basic(name, acceptance.HW_FGS_DEPENDENCY_OBS_LINK))
}
