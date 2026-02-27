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

func getV5AgencyBatchPoliciesAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.ListV5AgencyAssociatedPolicies(client, state.Primary.ID, nil)
}

func TestAccV5AgencyBatchPoliciesAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_identityv5_agency_batch_policies_associate.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5AgencyBatchPoliciesAssociateResourceFunc)

		name = acceptance.RandomAccResourceName()
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
				Config: testAccV5AgencyBatchPoliciesAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "agency_id", "huaweicloud_identity_agency.test", "id"),
					resource.TestCheckResourceAttr(rName, "policies.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "policies.0.id", "huaweicloud_identity_policy.test.0", "id"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.urn"),
					resource.TestMatchResourceAttr(rName, "policies.0.attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(rName, "policies.1.id", "huaweicloud_identity_policy.test.1", "id"),
					resource.TestCheckResourceAttrSet(rName, "policies.1.name"),
					resource.TestCheckResourceAttrSet(rName, "policies.1.urn"),
					resource.TestMatchResourceAttr(rName, "policies.1.attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV5AgencyBatchPoliciesAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "agency_id", "huaweicloud_identity_agency.test", "id"),
					resource.TestCheckResourceAttr(rName, "policies.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "policies.0.id", "huaweicloud_identity_policy.test.1", "id"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.urn"),
					resource.TestMatchResourceAttr(rName, "policies.0.attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(rName, "policies.1.id", "huaweicloud_identity_policy.test.2", "id"),
					resource.TestCheckResourceAttrSet(rName, "policies.1.name"),
					resource.TestCheckResourceAttrSet(rName, "policies.1.urn"),
					resource.TestMatchResourceAttr(rName, "policies.1.attached_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// After importing (and since origin is empty), the remote values ​​will be inconsistent with
					// the policies in the configuration, so it needs to be ignored.
					"policies",
					"policies_origin",
				},
			},
		},
	})
}

func testAccV5AgencyBatchPoliciesAssociate_basic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%[1]s"
  description           = "Created by terraform script"
  delegated_domain_name = "%[2]s"
}

resource "huaweicloud_identity_policy" "test" {
  count = 3

  name            = format("%[1]s_%%d", count.index)
  description     = "Created by terraform script"
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
`, name, acceptance.HW_DOMAIN_NAME)
}

func testAccV5AgencyBatchPoliciesAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_agency_batch_policies_associate" "test" {
  agency_id = huaweicloud_identity_agency.test.id

  dynamic "policies" {
    for_each = slice(huaweicloud_identity_policy.test[*].id, 0, 2)

    content {
      id = policies.value
    }
  }
}
`, testAccV5AgencyBatchPoliciesAssociate_basic_base(name))
}

func testAccV5AgencyBatchPoliciesAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_agency_batch_policies_associate" "test" {
  agency_id = huaweicloud_identity_agency.test.id

  dynamic "policies" {
    for_each = slice(huaweicloud_identity_policy.test[*].id, 1, 3)

    content {
      id = policies.value
    }
  }
}
`, testAccV5AgencyBatchPoliciesAssociate_basic_base(name))
}
