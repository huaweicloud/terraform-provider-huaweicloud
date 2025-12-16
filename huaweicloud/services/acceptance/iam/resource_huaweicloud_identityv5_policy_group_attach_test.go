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

func getIdentityV5PolicyGroupAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetGroupAttachedIdentityV5Policy(client, state.Primary.Attributes["group_id"], state.Primary.Attributes["policy_id"])
}

func TestAccIdentityV5PolicyGroupAttach_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_policy_group_attach.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getIdentityV5PolicyGroupAttachResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV5PolicyGroupAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_identity_policy.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_identityv5_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "policy_name", "huaweicloud_identity_policy.test", "name"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestMatchResourceAttr(rName, "attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityV5PolicyGroupAttachImportState(rName),
			},
		},
	})
}

func testAccIdentityV5PolicyGroupAttach_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
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

resource "huaweicloud_identityv5_group" "test" {
  group_name = "%[1]s"
}

resource "huaweicloud_identityv5_policy_group_attach" "test" {
  policy_id = huaweicloud_identity_policy.test.id
  group_id  = huaweicloud_identityv5_group.test.id
}
`, rName)
}

func testAccIdentityV5PolicyGroupAttachImportState(rName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		policyId := rs.Primary.Attributes["policy_id"]
		groupId := rs.Primary.Attributes["group_id"]
		if policyId == "" || groupId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<group_id>', but got '%s/%s'",
				policyId, groupId)
		}

		return fmt.Sprintf("%s/%s", policyId, groupId), nil
	}
}
