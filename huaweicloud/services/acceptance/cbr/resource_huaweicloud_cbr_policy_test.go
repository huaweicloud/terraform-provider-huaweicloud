package cbr

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CbrV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CBR client: %s", err)
	}
	return policies.Get(c, state.Primary.ID).Extract()
}

func TestAccCBRV3Policy_basic(t *testing.T) {
	var asPolicy policies.Policy
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&asPolicy,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "MO,TU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
				),
			},
			{
				Config: testCBRV3Policy_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
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

func TestAccCBRV3Policy_retention(t *testing.T) {
	var asPolicy policies.Policy
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&asPolicy,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_retention(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "15"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.daily", "10"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.weekly", "10"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.monthly", "1"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "SA,SU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "20:00"),
				),
			},
			{
				Config: testCBRV3Policy_retention_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "35"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.daily", "20"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.weekly", "20"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.monthly", "6"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.yearly", "1"),
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
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cbr_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&asPolicy,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_replication(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "destination_region", acceptance.HW_DEST_REGION),
					resource.TestCheckResourceAttr(resourceName, "destination_project_id", acceptance.HW_DEST_PROJECT_ID),
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

func testCBRV3Policy_retention(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name            = "%s"
  type            = "backup"
  backup_quantity = 15

  time_zone       = "UTC+08:00"
  long_term_retention {
    daily   = 10
    weekly  = 10
    monthly = 1
  }

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
`, rName)
}

func testCBRV3Policy_retention_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name            = "%s"
  type            = "backup"
  backup_quantity = 35

  time_zone       = "UTC+08:00"
  long_term_retention {
    daily   = 20
    weekly  = 20
    monthly = 6
    yearly  = 1
  }

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
`, rName, acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}
