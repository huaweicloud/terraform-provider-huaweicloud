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

func getV5AccessKeyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5AccessKeyById(client, state.Primary.Attributes["user_id"], state.Primary.ID)
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5AccessKey_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_access_key.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5AccessKeyResourceFunc)
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
				Config: testAccAccessKeyV5_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "user_id", "huaweicloud_identityv5_user.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "inactive"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "access_key_id"),
					resource.TestCheckResourceAttrSet(rName, "secret_access_key"),
					// If the access key has never been used, 'last_used_at' is empty, so we don't check it.
				),
			},
			{
				Config: testAccAccessKeyV5_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "status", "active"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccV5AccessKeyImportState(rName),
				ImportStateVerifyIgnore: []string{"secret_access_key"},
			},
		},
	})
}

func testAccAccessKeyV5_basic(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
 name        = "%[1]s"
 description = "tested by terraform"
}

resource "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_user.test.id
  status  = "inactive"
}
`, userName)
}

func testAccAccessKeyV5_update(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
 name        = "%[1]s"
 description = "tested by terraform"
}

resource "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_user.test.id
  status  = "active"
}
`, userName)
}

func testAccV5AccessKeyImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		userId := rs.Primary.Attributes["user_id"]
		if userId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<user_id>/<id>', but got '%s/%s'",
				userId, rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", userId, rs.Primary.ID), nil
	}
}
