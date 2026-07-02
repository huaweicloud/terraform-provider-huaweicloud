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

func getV5PolicyAgencyAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5AgencyAttachedPolicy(client, state.Primary.Attributes["agency_id"], state.Primary.Attributes["policy_id"])
}

func TestAccV5PolicyAgencyAttach_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_policy_agency_attach.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5PolicyAgencyAttachResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5PolicyAgencyAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_identity_policy.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "agency_id", "huaweicloud_identity_agency.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "policy_name", "huaweicloud_identity_policy.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "policy_urn", "huaweicloud_identity_policy.test", "urn"),
					resource.TestMatchResourceAttr(rName, "attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
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

func testAccV5PolicyAgencyAttach_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%[1]s"
  delegated_domain_name = "%[2]s"
  description           = "Created by terraform acceptance test"
  duration              = "30"
}

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

resource "huaweicloud_identityv5_policy_agency_attach" "test" {
  policy_id = huaweicloud_identity_policy.test.id
  agency_id = huaweicloud_identity_agency.test.id
}
`, name, acceptance.HW_DOMAIN_NAME)
}
