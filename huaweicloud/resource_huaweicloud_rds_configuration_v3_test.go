package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rds/v3/configurations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccRdsConfigurationV3_basic(t *testing.T) {
	var config configurations.Configuration
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-acc-test-%s-update", acctest.RandString(5))
	resourceName := "huaweicloud_rds_parametergroup.pg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsConfigV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsConfigV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigV3Exists(resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "description_1"),
				),
			},
			{
				Config: testAccRdsConfigV3_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigV3Exists(resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update"),
				),
			},
		},
	})
}

func testAccCheckRdsConfigV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	rdsClient, err := config.RdsV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud RDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_rds_parametergroup" {
			continue
		}

		_, err := configurations.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Rds configuration still exists")
		}
	}

	return nil
}

func testAccCheckRdsConfigV3Exists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		rdsClient, err := config.RdsV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud RDS client: %s", err)
		}

		found, err := configurations.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("Rds configuration not found")
		}

		*configuration = *found

		return nil
	}
}

func testAccRdsConfigV3_basic(rName string) string {
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
    version = "5.6"
  }
}
`, rName)
}

func testAccRdsConfigV3_update(updateName string) string {
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
    version = "5.6"
  }
}
`, updateName)
}
