package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDependencyVersions_basic(t *testing.T) {
	var (
		base = "huaweicloud_fgs_dependency_version.test"

		all                        = "data.huaweicloud_fgs_dependency_versions.all"
		dcForAllDependencyVersions = acceptance.InitDataSourceCheck(all)

		byVersionId           = "data.huaweicloud_fgs_dependency_versions.filter_by_version_id"
		dcByVersionId         = acceptance.InitDataSourceCheck(byVersionId)
		byNotFoundVersionId   = "data.huaweicloud_fgs_dependency_versions.filter_by_not_found_version_id"
		dcByNotFoundVersionId = acceptance.InitDataSourceCheck(byNotFoundVersionId)

		byVersion           = "data.huaweicloud_fgs_dependency_versions.filter_by_version"
		dcByVersion         = acceptance.InitDataSourceCheck(byVersion)
		byNotFoundVersion   = "data.huaweicloud_fgs_dependency_versions.filter_by_not_found_version"
		dcByNotFoundVersion = acceptance.InitDataSourceCheck(byNotFoundVersion)

		byRuntime           = "data.huaweicloud_fgs_dependency_versions.filter_by_runtime"
		dcByRuntime         = acceptance.InitDataSourceCheck(byRuntime)
		byNotFoundRuntime   = "data.huaweicloud_fgs_dependency_versions.filter_by_not_found_runtime"
		dcByNotFoundRuntime = acceptance.InitDataSourceCheck(byNotFoundRuntime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsDependencyLink(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDependencyVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllDependencyVersions.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "versions.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by version ID.
					dcByVersionId.CheckResourceExists(),
					resource.TestCheckOutput("is_version_id_filter_useful", "true"),
					dcByNotFoundVersionId.CheckResourceExists(),
					resource.TestCheckOutput("version_id_not_found_validation_pass", "true"),
					// Filter by version number.
					dcByVersion.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
					dcByNotFoundVersion.CheckResourceExists(),
					resource.TestCheckOutput("version_not_found_validation_pass", "true"),
					// Filter by runtime of the dependency version.
					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
					dcByNotFoundRuntime.CheckResourceExists(),
					resource.TestCheckOutput("runtime_not_found_validation_pass", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.id", base, "version_id"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.version", base, "version"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.dependency_id", base, "dependency_id"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.dependency_name", base, "name"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.runtime", base, "runtime"),
					// The link will be replaced with a new link which belongs to the FunctionGraph service.
					resource.TestCheckResourceAttrSet(byVersionId, "versions.0.link"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.size", base, "size"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.etag", base, "etag"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.owner", base, "owner"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.description", base, "description"),
				),
			},
		},
	})
}

func testAccDataDependencyVersions_basic() string {
	var (
		name             = acceptance.RandomAccResourceName()
		randVersionId, _ = uuid.GenerateUUID()
		randVersion      = acctest.RandIntRange(10000, 99999)
	)

	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_fgs_dependency_versions" "all" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
}

# Filter by version ID.
locals {
  version_id = huaweicloud_fgs_dependency_version.test.version_id
}

data "huaweicloud_fgs_dependency_versions" "filter_by_version_id" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  version_id    = local.version_id
}

data "huaweicloud_fgs_dependency_versions" "filter_by_not_found_version_id" {
  # Query dependency versions using a not exist version ID after dependency version resource create.
  depends_on = [
    huaweicloud_fgs_dependency_version.test,
  ]

  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  version_id    = "%[2]s"
}

locals {
  version_id_filter_result = [for v in data.huaweicloud_fgs_dependency_versions.filter_by_version_id.versions[*].id :
    v == local.version_id]
}

output "is_version_id_filter_useful" {
  value = length(local.version_id_filter_result) > 0 && alltrue(local.version_id_filter_result)
}

output "version_id_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_dependency_versions.filter_by_not_found_version_id.versions) == 0
}

# Filter by version.
locals {
  version = huaweicloud_fgs_dependency_version.test.version
}

data "huaweicloud_fgs_dependency_versions" "filter_by_version" {
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  version       = local.version
}

data "huaweicloud_fgs_dependency_versions" "filter_by_not_found_version" {
  # Query dependency versions using a not exist version number after dependency version resource create.
  depends_on = [
    huaweicloud_fgs_dependency_version.test,
  ]

  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  version       = %[3]d
}

locals {
  version_filter_result = [for v in data.huaweicloud_fgs_dependency_versions.filter_by_version.versions[*].version :
    v == local.version]
}

output "is_version_filter_useful" {
  value = length(local.version_filter_result) > 0 && alltrue(local.version_filter_result)
}

output "version_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_dependency_versions.filter_by_not_found_version.versions) == 0
}

# Filter by runtime.
locals {
  version_runtime = huaweicloud_fgs_dependency_version.test.runtime
}

data "huaweicloud_fgs_dependency_versions" "filter_by_runtime" {
  # The behavior of parameter 'runtime' of the resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  runtime       = local.version_runtime
}

data "huaweicloud_fgs_dependency_versions" "filter_by_not_found_runtime" {
  # Query dependency versions using a not exist version runtime after dependency version resource create.
  depends_on = [
    huaweicloud_fgs_dependency_version.test,
  ]

  dependency_id = huaweicloud_fgs_dependency_version.test.dependency_id
  runtime       = "runtime_not_found"
}

locals {
  runtime_filter_result = [for v in data.huaweicloud_fgs_dependency_versions.filter_by_runtime.versions[*].runtime :
    v == local.version_runtime]
}

output "is_runtime_filter_useful" {
  value = length(local.runtime_filter_result) > 0 && alltrue(local.runtime_filter_result)
}

output "runtime_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_dependency_versions.filter_by_not_found_runtime.versions) == 0
}
`, testAccDependencyVersion_basic(name), randVersionId, randVersion)
}
