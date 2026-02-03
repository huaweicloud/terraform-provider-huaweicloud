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

func getV3GroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return iam.GetV3GroupById(client, state.Primary.ID)
}

func TestAccV3Group_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_group.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3GroupResourceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
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
				Config: testAccV3Group_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
				),
			},
			{
				Config: testAccV3Group_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV3Group_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"
}
`, name)
}

func testAccV3Group_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name = "%[1]s"
}
`, name)
}
