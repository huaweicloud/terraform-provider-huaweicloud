package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDryRunPolicyAttachedEntities_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_organizations_dry_run_policy_attached_entities.test"
		dc  = acceptance.InitDataSourceCheck(all)
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
				Config: testAccDataDryRunPolicyAttachedEntities_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "entities.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "entities.0.id"),
					resource.TestCheckResourceAttrSet(all, "entities.0.type"),
					resource.TestCheckResourceAttrSet(all, "entities.0.name"),
				),
			},
		},
	})
}

func testAccDataDryRunPolicyAttachedEntities_base(name string) string {
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

resource "huaweicloud_organizations_dry_run_policy_entity_attach" "test" {
  policy_id = huaweicloud_organizations_dry_run_policy.test.id
  entity_id = data.huaweicloud_organizations_organization.test.root_id

  depends_on = [huaweicloud_organizations_policy_dry_run_configuration.test]
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccDataDryRunPolicyAttachedEntities_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_organizations_dry_run_policy_attached_entities" "test" {
  policy_id  = huaweicloud_organizations_dry_run_policy.test.id

  depends_on = [huaweicloud_organizations_dry_run_policy_entity_attach.test]
}
`, testAccDataDryRunPolicyAttachedEntities_base(name))
}
