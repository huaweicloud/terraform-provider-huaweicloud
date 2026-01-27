package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5GroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5GroupById(client, state.Primary.ID)
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5Group_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_group.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5GroupFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5Group_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_name", name),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
				),
			},
			{
				Config: testAccV5Group_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_name", updateName),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
				),
			},
			{
				Config: testAccV5Group_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_name", updateName),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV5Group_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "test" {
  group_name  = "%[1]s"
  description = "Created by terraform script"
}
`, name)
}

func testAccV5Group_basic_step2(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "test" {
  group_name  = "%[1]s"
  description = "Updated by terraform script"
}
`, updateName)
}

func testAccV5Group_basic_step3(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "test" {
  group_name = "%[1]s"
}
`, updateName)
}
