package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInstanceGroupAssign_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_instance_group_assign.test"

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupAssign_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "group_id"),
					resource.TestMatchResourceAttr(rName, "instance_ids.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
				),
			},
			{
				Config: testAccInstanceGroupAssign_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "group_id"),
					resource.TestMatchResourceAttr(rName, "instance_ids.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
				),
			},
		},
	})
}

func testAccInstanceGroupAssign_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_das_instance_group" "test" {
  datastore_type = "MySQL"
  group_name     = "%[1]s"
  description    = "Created by terraform script"
}

locals {
  instance_ids = split(",", "%[2]s")
}

`, name, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccInstanceGroupAssign_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_instance_group_assign" "test" {
  group_id     = huaweicloud_das_instance_group.test.id
  instance_ids = slice(local.instance_ids, 0, 1) 
}
`, testAccInstanceGroupAssign_base(name))
}

func testAccInstanceGroupAssign_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

# Assigning an instance group to the 0th instance repeatedly will not result in an error
resource "huaweicloud_das_instance_group_assign" "test" {
  group_id     = huaweicloud_das_instance_group.test.id
  instance_ids = local.instance_ids

  enable_force_new = "true"
}
`, testAccInstanceGroupAssign_base(name))
}
