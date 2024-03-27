package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/rds/v3/configurations"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsConfiguration_basic(t *testing.T) {
	var config configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	updateName := rName + "-update"
	resourceName := "huaweicloud_rds_parametergroup.pg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigExists(resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "description_1"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "8.0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccRdsConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigExists(resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update"),
				),
			},
		},
	})
}

func testAccCheckRdsConfigDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	rdsClient, err := config.RdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating RDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_rds_parametergroup" {
			continue
		}

		_, err := configurations.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("RDS configuration still exists")
		}
	}

	return nil
}

func testAccCheckRdsConfigExists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		rdsClient, err := config.RdsV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating RDS client: %s", err)
		}

		found, err := configurations.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("RDS configuration not found")
		}

		*configuration = *found

		return nil
	}
}

func testAccRdsConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_parametergroup" "pg_1" {
  name        = "%s"
  description = "description_1"

  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }
  datastore {
    type    = "mysql"
    version = "8.0"
  }
}
`, rName)
}

func testAccRdsConfig_update(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_parametergroup" "pg_1" {
  name        = "%s"
  description = "description_update"

  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }
  datastore {
    type    = "mysql"
    version = "8.0"
  }
}
`, updateName)
}
