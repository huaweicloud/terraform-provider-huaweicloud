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

func getV5ResourceTagResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5ResourceTagsById(client, state.Primary.Attributes["resource_type"], state.Primary.Attributes["resource_id"])
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5ResourceTag_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_resource_tag.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5ResourceTagResourceFunc)
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
				Config: testAccV5ResourceTag_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resource_type", "user"),
					resource.TestCheckResourceAttrSet(rName, "resource_id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config: testAccV5ResourceTag_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resource_type", "user"),
					resource.TestCheckResourceAttrSet(rName, "resource_id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value"),
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

func testAccV5ResourceTag_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}
`, name)
}

func testAccV5ResourceTag_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccV5ResourceTag_base(name))
}

func testAccV5ResourceTag_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.test.id

  tags = {
    foo  = "bar1"
    key1 = "value"
  }
}
`, testAccV5ResourceTag_base(name))
}
