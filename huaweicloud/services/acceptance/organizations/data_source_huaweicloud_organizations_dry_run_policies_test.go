package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDryRunPolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_organizations_dry_run_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byPolicyType   = "data.huaweicloud_organizations_dry_run_policies.filter_by_policy_type"
		dcByPolicyType = acceptance.InitDataSourceCheck(byPolicyType)

		byAttachedEntityId   = "data.huaweicloud_organizations_dry_run_policies.filter_by_attached_entity_id"
		dcByAttachedEntityId = acceptance.InitDataSourceCheck(byAttachedEntityId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDryRunPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'policy_type' parameter.
					dcByPolicyType.CheckResourceExists(),
					resource.TestCheckOutput("is_policy_type_filter_useful", "true"),
					// Filter by 'attached_entity_id' parameter.
					dcByAttachedEntityId.CheckResourceExists(),
					resource.TestCheckOutput("attached_entity_id_filter_result", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(all, "policies.0.id"),
					resource.TestCheckResourceAttrSet(all, "policies.0.name"),
					resource.TestCheckResourceAttrSet(all, "policies.0.type"),
					resource.TestCheckResourceAttrSet(all, "policies.0.urn"),
					resource.TestCheckResourceAttrSet(all, "policies.0.is_builtin"),
					resource.TestCheckOutput("is_description_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataDryRunPolicies_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_organizations_organization" "test" {}

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
          Action = ["sts:agencies:assume"]
          Effect = "Allow"
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

resource "huaweicloud_organizations_dry_run_policy" "test" {
  name        = "%[2]s"
  type        = "service_control_policy"
  description = "Created by terraform script"
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

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[2]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}

resource "huaweicloud_organizations_dry_run_policy_entity_attach" "test" {
  policy_id = huaweicloud_organizations_dry_run_policy.test.id
  entity_id = huaweicloud_organizations_organizational_unit.test.id

  depends_on = [huaweicloud_organizations_policy_dry_run_configuration.test]
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccDataDryRunPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_organizations_dry_run_policies" "test" {
  depends_on = [huaweicloud_organizations_dry_run_policy_entity_attach.test]
}

# Filter by 'policy_type' parameter.
locals {
  policy_type = huaweicloud_organizations_dry_run_policy.test.type
}

data "huaweicloud_organizations_dry_run_policies" "filter_by_policy_type" {
  policy_type = local.policy_type

  depends_on = [huaweicloud_organizations_dry_run_policy.test]
}

locals {
  policy_type_filter_result = [for v in data.huaweicloud_organizations_dry_run_policies.filter_by_policy_type.policies[*].type :
  v == local.policy_type]
}

output "is_policy_type_filter_useful" {
  value = length(local.policy_type_filter_result) > 0 && alltrue(local.policy_type_filter_result)
}

# Filter by 'attached_entity_id' parameter.
locals {
  policy_id = huaweicloud_organizations_dry_run_policy.test.id
}

data "huaweicloud_organizations_dry_run_policies" "filter_by_attached_entity_id" {
  attached_entity_id = huaweicloud_organizations_organizational_unit.test.id

  depends_on = [huaweicloud_organizations_dry_run_policy_entity_attach.test]
}

locals {
  attached_entity_id_filter_result = [for v in data.huaweicloud_organizations_dry_run_policies.filter_by_attached_entity_id.policies :
  v if v.id == local.policy_id]
}

output "attached_entity_id_filter_result" {
  value = try(local.attached_entity_id_filter_result[0].id == local.policy_id, false)
}

output "is_description_set_and_valid" {
  value = try(local.attached_entity_id_filter_result[0].description == huaweicloud_organizations_dry_run_policy.test.description, false)
}
`, testAccDataDryRunPolicies_base(name))
}
