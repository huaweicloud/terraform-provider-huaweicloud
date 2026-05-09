package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNotebooks_basic(t *testing.T) {
	var (
		rName = "huaweicloud_modelarts_notebook.test"

		all   = "data.huaweicloud_modelarts_notebooks.all"
		dcAll = acceptance.InitDataSourceCheck(all)

		byFeature   = "data.huaweicloud_modelarts_notebooks.filter_by_feature"
		dcByFeature = acceptance.InitDataSourceCheck(byFeature)

		byNotebookId   = "data.huaweicloud_modelarts_notebooks.filter_by_notebook_id"
		dcByNotebookId = acceptance.InitDataSourceCheck(byNotebookId)

		byName   = "data.huaweicloud_modelarts_notebooks.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_modelarts_notebooks.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byWorkspaceId   = "data.huaweicloud_modelarts_notebooks.filter_by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)

		byFlavor   = "data.huaweicloud_modelarts_notebooks.filter_by_flavor"
		dcByFlavor = acceptance.InitDataSourceCheck(byFlavor)

		byImageId   = "data.huaweicloud_modelarts_notebooks.filter_by_image_id"
		dcByImageId = acceptance.InitDataSourceCheck(byImageId)

		byBilling   = "data.huaweicloud_modelarts_notebooks.filter_by_billing"
		dcByBilling = acceptance.InitDataSourceCheck(byBilling)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRunnerPublicIPs(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNotebooks_nonExistentNotebook(),
				Check: resource.ComposeTestCheckFunc(
					dcAll.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "notebooks.#", "0"),
				),
			},
			{
				Config: testAccDataNotebooks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "notebooks.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByFeature.CheckResourceExists(),
					resource.TestCheckOutput("is_feature_filter_useful", "true"),
					dcByNotebookId.CheckResourceExists(),
					resource.TestCheckOutput("is_notebook_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byNotebookId, "notebooks.0.status", "RUNNING"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.id", rName, "id"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.name", rName, "name"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.flavor_id", rName, "flavor_id"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.image_id", rName, "image_id"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.image_type", rName, "image_type"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.image_swr_path", rName, "image_swr_path"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.image_name", rName, "image_name"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.description", rName, "description"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.key_pair", rName, "key_pair"),
					resource.TestCheckResourceAttrPair(byNotebookId, "notebooks.0.workspace_id", rName, "workspace_id"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
					dcByFlavor.CheckResourceExists(),
					resource.TestCheckOutput("is_flavor_filter_useful", "true"),
					dcByImageId.CheckResourceExists(),
					resource.TestCheckOutput("is_image_id_filter_useful", "true"),
					dcByBilling.CheckResourceExists(),
					resource.TestCheckOutput("is_billing_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataNotebooks_nonExistentNotebook() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
// Create a new workspace to make sure no notebook exists in the workspace.
resource "huaweicloud_modelarts_workspace" "test" {
  name = "%[1]s"
}

data "huaweicloud_modelarts_notebooks" "all" {
  workspace_id = huaweicloud_modelarts_workspace.test.id
}
`, name)
}

func testAccDataNotebooks_basic_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
data "huaweicloud_modelarts_notebook_flavors" "test" {
  type     = "MANAGED"
  category = "CPU"
}

data "huaweicloud_modelarts_notebook_images" "test" {
  type     = "BUILD_IN"
  cpu_arch = try(data.huaweicloud_modelarts_notebook_flavors.test.flavors[0].arch, "x86_64")
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_modelarts_notebook" "test" {
  depends_on = [huaweicloud_kps_keypair.test]

  name               = "%[1]s"
  flavor_id          = data.huaweicloud_modelarts_notebook_flavors.test.flavors[0].id
  image_id           = data.huaweicloud_modelarts_notebook_images.test.images[0].id
  description        = "Created by Terraform"
  key_pair           = huaweicloud_kps_keypair.test.name
  allowed_access_ips = split(",", "%[2]s")
  workspace_id       = "0"

  volume {
    type = "EVS"
    size = 100
  }
}
`, name, acceptance.HW_RUNNER_PUBLIC_IPS)
}

func testAccDataNotebooks_basic() string {
	return fmt.Sprintf(`
%[1]s

# Query without any filter.
data "huaweicloud_modelarts_notebooks" "all" {
  depends_on = [huaweicloud_modelarts_notebook.test]
}

# Filter by 'feature' parameter.
locals {
  notebook_feature = "NOTEBOOK"
}

data "huaweicloud_modelarts_notebooks" "filter_by_feature" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  feature = local.notebook_feature
}

locals {
  feature_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_feature.notebooks :
    v.feature == local.notebook_feature]
}

output "is_feature_filter_useful" {
  value = length(local.feature_filter_result) > 0 && alltrue(local.feature_filter_result)
}

# Filter by 'name' parameter.
data "huaweicloud_modelarts_notebooks" "filter_by_name" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  name = huaweicloud_modelarts_notebook.test.name
}

locals {
  name_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_name.notebooks :
    v.name == huaweicloud_modelarts_notebook.test.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'notebook_id' parameter.
locals {
  notebook_id = huaweicloud_modelarts_notebook.test.id
}

data "huaweicloud_modelarts_notebooks" "filter_by_notebook_id" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  notebook_id = local.notebook_id
}

locals {
  notebook_id_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_notebook_id.notebooks :
    v.id == local.notebook_id]
}

output "is_notebook_id_filter_useful" {
  value = length(local.notebook_id_filter_result) > 0 && alltrue(local.notebook_id_filter_result)
}

# Filter by 'status' parameter.
locals {
  notebook_status = huaweicloud_modelarts_notebook.test.status
}

data "huaweicloud_modelarts_notebooks" "filter_by_status" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  status = local.notebook_status
}

locals {
  status_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_status.notebooks :
    v.status == local.notebook_status]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by 'workspace_id' parameter.
locals {
  workspace_id = huaweicloud_modelarts_notebook.test.workspace_id
}

data "huaweicloud_modelarts_notebooks" "filter_by_workspace_id" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  workspace_id = local.workspace_id
}

locals {
  workspace_id_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_workspace_id.notebooks :
    v.workspace_id == local.workspace_id]
}

output "is_workspace_id_filter_useful" {
  value = length(local.workspace_id_filter_result) > 0 && alltrue(local.workspace_id_filter_result)
}

# Filter by 'flavor_id' parameter.
locals {
  notebook_flavor_id = huaweicloud_modelarts_notebook.test.flavor_id
}

data "huaweicloud_modelarts_notebooks" "filter_by_flavor" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  flavor_id = local.notebook_flavor_id
}

locals {
  flavor_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_flavor.notebooks :
    v.flavor_id == local.notebook_flavor_id]
}

output "is_flavor_filter_useful" {
  value = length(local.flavor_filter_result) > 0 && alltrue(local.flavor_filter_result)
}

# Filter by 'image_id' parameter.
locals {
  notebook_image_id = huaweicloud_modelarts_notebook.test.image_id
}

data "huaweicloud_modelarts_notebooks" "filter_by_image_id" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  image_id = local.notebook_image_id
}

locals {
  image_id_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_image_id.notebooks :
    v.image_id == local.notebook_image_id]
}

output "is_image_id_filter_useful" {
  value = length(local.image_id_filter_result) > 0 && alltrue(local.image_id_filter_result)
}

# Filter by 'billing' parameter.
locals {
  notebook_billing = "STORAGE"
}

data "huaweicloud_modelarts_notebooks" "filter_by_billing" {
  depends_on = [huaweicloud_modelarts_notebook.test]

  billing = local.notebook_billing
}

locals {
  billing_filter_result = [for v in data.huaweicloud_modelarts_notebooks.filter_by_billing.notebooks :
    contains(v.billing_items, local.notebook_billing)]
}

output "is_billing_filter_useful" {
  value = length(local.billing_filter_result) > 0 && anytrue(local.billing_filter_result)
}
`, testAccDataNotebooks_basic_base())
}
