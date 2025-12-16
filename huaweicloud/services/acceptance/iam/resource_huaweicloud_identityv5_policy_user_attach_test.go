package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getIdentityV5PolicyUserAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetUserAttachedIdentityV5Policy(client, state.Primary.Attributes["user_id"], state.Primary.Attributes["policy_id"])
}

func TestAccIdentityV5PolicyUserAttach_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_policy_user_attach.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getIdentityV5PolicyUserAttachResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV5PolicyUserAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_identity_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttrPair(rName, "policy_name", "huaweicloud_identity_policy.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "urn", "huaweicloud_identity_policy.test", "urn"),
					resource.TestMatchResourceAttr(rName, "attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityV5PolicyUserAttachImportState(rName),
			},
		},
	})
}

func testAccIdentityV5PolicyUserAttach_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  description     = "test policy for terraform"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Allow"
        }
      ]
      Version = "5.0"
    }
  )
}

resource "huaweicloud_identityv5_policy_user_attach" "test" {
  policy_id = huaweicloud_identity_policy.test.id
  user_id   = "%[2]s"
}
`, name, acceptance.HW_USER_ID)
}

func testAccIdentityV5PolicyUserAttachImportState(rName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		policyId := rs.Primary.Attributes["policy_id"]
		userId := rs.Primary.Attributes["user_id"]
		if policyId == "" || userId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<user_id>', but got '%s/%s'",
				policyId, userId)
		}

		return fmt.Sprintf("%s/%s", policyId, userId), nil
	}
}
