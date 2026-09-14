package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func getVaultPoliciesAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cbr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR client: %s", err)
	}
	return cbr.GetVaultAssociatedPolicies(client, state.Primary.ID)
}

// Make sure at least one resource backup is existed before the test.
func TestAccVaultPoliciesAssociate_basic(t *testing.T) {
	var (
		obj          interface{}
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_cbr_vault_policies_associate.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getVaultPoliciesAssociateResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVaultPoliciesAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vault_id", "huaweicloud_cbr_vault.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "policies.*.id",
						"huaweicloud_cbr_policy.backup.0", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "policies.*.id",
						"huaweicloud_cbr_policy.replication.0", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "policies.*.destination_vault_id",
						"huaweicloud_cbr_vault.destination", "id"),
				),
			},
			{
				Config: testAccVaultPoliciesAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vault_id", "huaweicloud_cbr_vault.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "policies.*.id",
						"huaweicloud_cbr_policy.backup.1", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "policies.*.id",
						"huaweicloud_cbr_policy.replication.1", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "policies.*.destination_vault_id",
						"huaweicloud_cbr_vault.destination", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVaultPoliciesAssociate_basic_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_cbr_policy" "backup" {
  count = 2

  name            = format("%[2]s_%%d", count.index)
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    interval        = 1
    execution_times = ["06:00"]
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  protection_type       = "backup"
  size                  = 100
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_cbr_vault" "destination" {
  depends_on = [huaweicloud_cbr_vault.test]

  region                = "%[3]s" # Vault in the destination region which used to store the replicated data.
  name                  = "%[2]s_destination"
  type                  = "server"
  protection_type       = "replication"
  size                  = 100
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_cbr_policy" "replication" {
  depends_on = [
  	huaweicloud_cbr_vault.destination,
  	huaweicloud_cbr_policy.backup,
  ]

  count = 2

  name                   = format("%[2]s_replication_%%d", count.index)
  type                   = "replication"
  destination_region     = "%[3]s"
  destination_project_id = "%[4]s"
  time_period            = 20

  backup_cycle {
    interval        = 1
    execution_times = ["08:00"]
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID, name, acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}

func testAccVaultPoliciesAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault_policies_associate" "test" {
  vault_id = huaweicloud_cbr_vault.test.id

  policies {
    id = huaweicloud_cbr_policy.backup[0].id
  }
  policies {
    id                   = huaweicloud_cbr_policy.replication[0].id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
`, testAccVaultPoliciesAssociate_basic_base(name))
}

func testAccVaultPoliciesAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault_policies_associate" "test" {
  vault_id = huaweicloud_cbr_vault.test.id

  policies {
    id = huaweicloud_cbr_policy.backup[1].id
  }
  policies {
    id                   = huaweicloud_cbr_policy.replication[1].id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
`, testAccVaultPoliciesAssociate_basic_base(name))
}
