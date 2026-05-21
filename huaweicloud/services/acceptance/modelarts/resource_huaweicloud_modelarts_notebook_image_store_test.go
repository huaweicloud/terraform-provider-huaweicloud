package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNotebookImageStore_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_modelarts_notebook_image_store.test"
		name         = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccNotebookImageStore_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tag", "v1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "notebook_id",
						"huaweicloud_modelarts_notebook.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "namespace", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttrSet(resourceName, "swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "arch"),
					resource.TestCheckResourceAttrSet(resourceName, "origin"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccNotebookImageStore_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_modelarts_notebook_flavors" "test" {
  type     = "MANAGED"
  category = "CPU"
}

locals {
  available_notebook_flavors_with_onsale = [
    for o in data.huaweicloud_modelarts_notebook_flavors.test.flavors : o if
    !o.sold_out && !strcontains(o.id, "free")
  ]
}

data "huaweicloud_modelarts_notebook_images" "test" {
  type     = "BUILD_IN"
  cpu_arch = try(data.huaweicloud_modelarts_notebook_flavors.test.flavors[0].arch, "x86_64")
}

locals {
  available_notebook_images_with_onsale = [
    for o in data.huaweicloud_modelarts_notebook_images.test.images : o if
    contains(o.resource_categories, "CPU") && o.status == "ACTIVE" && contains(o.dev_services, "NOTEBOOK")
  ]
}

resource "huaweicloud_modelarts_notebook" "test" {
  name      = "%[1]s"
  flavor_id = try(local.available_notebook_flavors_with_onsale[0].id, null)
  image_id  = try(local.available_notebook_images_with_onsale[0].id, null)

  volume {
    type      = "EVS"
    ownership = "MANAGED"
    size      = 40
  }

  lifecycle {
    ignore_changes = [image_id]
  }
}
`, name)
}

func testAccNotebookImageStore_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_notebook_image_store" "test" {
  notebook_id = huaweicloud_modelarts_notebook.test.id
  name        = "%[2]s"
  namespace   = "%[3]s"
  tag         = "v1.0.0"
  description = "Created by terraform script"

  enable_force_new = "true"
}
`, testAccNotebookImageStore_base(name), name, acceptance.HW_DOMAIN_NAME)
}
