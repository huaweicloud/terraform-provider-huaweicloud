package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
)

func getDryRunPolicyEntityAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetAttachedEntityForDryRunPolicy(client, state.Primary.Attributes["policy_id"],
		state.Primary.Attributes["entity_id"])
}

func TestAccDryRunPolicyEntityAttach_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_organizations_dry_run_policy_entity_attach.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDryRunPolicyEntityAttachResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDryRunPolicyEntityAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_organizations_dry_run_policy.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "entity_id",
						"data.huaweicloud_organizations_organization.test", "root_id"),
					resource.TestCheckResourceAttrSet(rName, "entity_name"),
					resource.TestCheckResourceAttrSet(rName, "entity_type"),
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

func testAccDryRunPolicyEntityAttach_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_dry_run_policy" "test" {
  name    = "%[2]s"
  type    = "service_control_policy"
  content = jsonencode({
    Version = "5.0",

    Statement = [
      {
        Effect = "Deny",
        Action = []
      }
    ]
  })
}

resource "huaweicloud_obs_bucket" "test" {
  bucket                = "%[2]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  force_destroy         = true
}

resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%[2]s"
  policy_names = ["OBSFullAccessPolicy"]
  trust_policy = jsonencode(
    {
      Version = "5.0"
      Statement = [
        {
          Action    = ["sts:agencies:assume"]
          Effect    = "Allow"
          Principal = {
            Service = ["service.Organizations"]
          }
        },
      ]
    }
  )
}

resource "huaweicloud_organizations_policy_dry_run_configuration" "test" {
  root_id     = data.huaweicloud_organizations_organization.test.root_id
  policy_type = "service_control_policy"
  status      = "enabled"
  bucket_name = huaweicloud_obs_bucket.test.bucket
  region_id   = huaweicloud_obs_bucket.test.region
  agency_name = huaweicloud_identity_trust_agency.test.name
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccDryRunPolicyEntityAttach_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_dry_run_policy_entity_attach" "test" {
  policy_id = huaweicloud_organizations_dry_run_policy.test.id
  entity_id = data.huaweicloud_organizations_organization.test.root_id

  depends_on = [huaweicloud_organizations_policy_dry_run_configuration.test]
}
`, testAccDryRunPolicyEntityAttach_base(name))
}
