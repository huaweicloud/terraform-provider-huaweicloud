package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCBRV3Policy_basic(t *testing.T) {
	var asPolicy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "MO,TU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
				),
			},
			{
				Config: testCBRV3Policy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "SA,SU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "20:00"),
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

func TestAccCBRV3Policy_replication(t *testing.T) {
	var asPolicy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckDestProject(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRV3PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_replication(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "destination_region", HW_DEST_REGION),
					resource.TestCheckResourceAttr(resourceName, "destination_project_id", HW_DEST_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
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

func testAccCheckCBRV3PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CbrV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating Huaweicloud CBR client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cbr_policy" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("policy still exists")
		}
	}
	return nil
}

func testAccCheckCBRV3PolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CbrV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating Huaweicloud CBR client: %s", err)
		}

		found, err := policies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		logp.Printf("[DEBUG] test found is: %#v", found)
		policy = found

		return nil
	}
}

func testCBRV3Policy_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name        = "%s"
  type        = "backup"
  time_period = 20

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}
`, rName)
}

func testCBRV3Policy_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name            = "%s-update"
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
`, rName)
}

func testCBRV3Policy_replication(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name                   = "%s"
  type                   = "replication"
  destination_region     = "%s"
  destination_project_id = "%s"
  time_period            = 20

  backup_cycle {
    interval        = 5
    execution_times = ["06:00", "18:00"]
  }
}
`, rName, HW_DEST_REGION, HW_DEST_PROJECT_ID)
}
