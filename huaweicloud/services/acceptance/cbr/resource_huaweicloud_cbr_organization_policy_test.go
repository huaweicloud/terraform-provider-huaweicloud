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

func getOrganizationPolicyResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := conf.NewServiceClient("cbr", region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR client: %s", err)
	}

	return cbr.GetOrganizationPolicyById(client, state.Primary.ID)
}

func TestAccOrganizationPolicy_basic(t *testing.T) {
	var (
		organizationPolicy interface{}
		name               = acceptance.RandomAccResourceName()
		updateName         = acceptance.RandomAccResourceName()
		resourceName       = "huaweicloud_cbr_organization_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&organizationPolicy,
		getOrganizationPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "operation_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "policy_name", fmt.Sprintf("%s_policy", name)),
					resource.TestCheckResourceAttr(resourceName, "policy_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.day_backups", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.max_backups", "30"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.month_backups", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.week_backups", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.year_backups", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.timezone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.full_backup_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "policy_trigger.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_trigger.0.properties.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_trigger.0.properties.0.pattern.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_trigger.0.properties.0.pattern.0",
						"FREQ=WEEKLY;BYDAY=WE,TH,FR;BYHOUR=16;BYMINUTE=00"),
					resource.TestCheckResourceAttrSet(resourceName, "effective_scope"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccOrganizationPolicy_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "operation_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "policy_name", fmt.Sprintf("%s_policy", updateName)),
					resource.TestCheckResourceAttr(resourceName, "policy_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.day_backups", "10"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.max_backups", "40"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.month_backups", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.week_backups", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.year_backups", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.timezone", "UTC+09:00"),
					resource.TestCheckResourceAttr(resourceName, "policy_operation_definition.0.full_backup_interval", "15"),
					resource.TestCheckResourceAttr(resourceName, "policy_trigger.0.properties.0.pattern.0",
						"FREQ=WEEKLY;BYDAY=SA,SU;BYHOUR=08;BYMINUTE=00"),
					resource.TestCheckResourceAttr(resourceName, "policy_trigger.0.properties.0.pattern.1",
						"FREQ=WEEKLY;BYDAY=WE,TH,FR,SA,TU,MO,SU;BYHOUR=14;BYMINUTE=00"),
					resource.TestCheckResourceAttrSet(resourceName, "effective_scope"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccOrganizationPolicyImportState(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_organization_policy" "test" {
  name           = "%[1]s"
  description    = "Created by terraform script"
  operation_type = "backup"
  policy_name    = "%[1]s_policy"
  policy_enabled = false

  policy_operation_definition {
    day_backups          = 5
    max_backups          = 30
    month_backups        = 1
    week_backups         = 2
    year_backups         = 1
    timezone             = "UTC+08:00"
    full_backup_interval = 10
  }

  policy_trigger {
    properties {
      pattern = ["FREQ=WEEKLY;BYDAY=WE,TH,FR;BYHOUR=16;BYMINUTE=00"]
    }
  }
}
`, name)
}

func testAccOrganizationPolicy_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_organization_policy" "test" {
  name           = "%[1]s"
  description    = "Updated by terraform script"
  operation_type = "backup"
  policy_name    = "%[1]s_policy"
  policy_enabled = true

  policy_operation_definition {
    day_backups          = 10
    max_backups          = 40
    month_backups        = 3
    week_backups         = 5
    year_backups         = 2
    timezone             = "UTC+09:00"
    full_backup_interval = 15
  }

  policy_trigger {
    properties {
      pattern = [
        "FREQ=WEEKLY;BYDAY=SA,SU;BYHOUR=08;BYMINUTE=00",
        "FREQ=WEEKLY;BYDAY=WE,TH,FR,SA,TU,MO,SU;BYHOUR=14;BYMINUTE=00",
      ]
    }
  }
}
`, name, name, name)
}

// testAccOrganizationPolicyImportState returns a function for importing a CBR organization policy by name
func testAccOrganizationPolicyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("attribute (name) of resource (%s) not found: %s", name, rs)
		}
		return rs.Primary.Attributes["name"], nil
	}
}
