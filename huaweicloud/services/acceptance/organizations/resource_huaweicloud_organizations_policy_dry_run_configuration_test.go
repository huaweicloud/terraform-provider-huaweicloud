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

func getPolicyDryRunConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetPolicyDryRunConfiguration(client, state.Primary.Attributes["root_id"],
		state.Primary.Attributes["policy_type"])
}

func TestAccPolicyDryRunConfiguration_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_organizations_policy_dry_run_configuration.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getPolicyDryRunConfigurationResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDryRunConfiguration_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "root_id",
						"data.huaweicloud_organizations_organization.test", "root_id"),
					resource.TestCheckResourceAttr(rName, "policy_type", "service_control_policy"),
					resource.TestCheckResourceAttr(rName, "status", "enabled"),
					resource.TestCheckResourceAttrPair(rName, "bucket_name",
						"huaweicloud_obs_bucket.test.0", "bucket"),
					resource.TestCheckResourceAttrPair(rName, "region_id",
						"huaweicloud_obs_bucket.test.0", "region"),
					resource.TestCheckResourceAttr(rName, "bucket_prefix", name),
					resource.TestCheckResourceAttrPair(rName, "agency_name",
						"huaweicloud_identity_trust_agency.test.0", "name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccPolicyDryRunConfiguration_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "enabled"),
					resource.TestCheckResourceAttrPair(rName, "bucket_name",
						"huaweicloud_obs_bucket.test.1", "bucket"),
					resource.TestCheckResourceAttrPair(rName, "region_id",
						"huaweicloud_obs_bucket.test.1", "region"),
					resource.TestCheckResourceAttr(rName, "bucket_prefix", ""),
					resource.TestCheckResourceAttrPair(rName, "agency_name",
						"huaweicloud_identity_trust_agency.test.1", "name"),
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

func testAccPolicyDryRunConfiguration_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_obs_bucket" "test" {
  count = 2

  bucket                = "%[1]s-${count.index}"
  enterprise_project_id = "%[2]s"
  force_destroy         = true
}

resource "huaweicloud_identity_trust_agency" "test" {
  count = 2

  name         = "%[1]s-${count.index}"
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
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccPolicyDryRunConfiguration_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_policy_dry_run_configuration" "test" {
  root_id       = data.huaweicloud_organizations_organization.test.root_id
  policy_type   = "service_control_policy"
  status        = "enabled"
  bucket_name   = try(huaweicloud_obs_bucket.test[0].bucket, null)
  region_id     = try(huaweicloud_obs_bucket.test[0].region, null)
  bucket_prefix = "%[2]s"
  agency_name   = try(huaweicloud_identity_trust_agency.test[0].name, null)
}
`, testAccPolicyDryRunConfiguration_base(name), name)
}

func testAccPolicyDryRunConfiguration_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_policy_dry_run_configuration" "test" {
  root_id     = data.huaweicloud_organizations_organization.test.root_id
  policy_type = "service_control_policy"
  status      = "enabled"
  bucket_name = try(huaweicloud_obs_bucket.test[1].bucket, null)
  region_id   = try(huaweicloud_obs_bucket.test[1].region, null)
  agency_name = try(huaweicloud_identity_trust_agency.test[1].name, null)
}
`, testAccPolicyDryRunConfiguration_base(name))
}
