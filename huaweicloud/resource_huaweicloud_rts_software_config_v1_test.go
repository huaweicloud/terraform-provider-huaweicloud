package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
)

func TestAccRtsSoftwareConfigV1_basic(t *testing.T) {
	var config softwareconfig.SoftwareConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsSoftwareConfigV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1Exists("huaweicloud_rts_software_config_v1.config_1", &config),
					resource.TestCheckResourceAttr(
						"huaweicloud_rts_software_config_v1.config_1", "name", "huaweicloud-config"),
					resource.TestCheckResourceAttr(
						"huaweicloud_rts_software_config_v1.config_1", "group", "script"),
				),
			},
		},
	})
}

func TestAccRtsSoftwareConfigV1_timeout(t *testing.T) {
	var config softwareconfig.SoftwareConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsSoftwareConfigV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1Exists("huaweicloud_rts_software_config_v1.config_1", &config),
				),
			},
		},
	})
}

func testAccCheckRtsSoftwareConfigV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud orchestration client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_rts_software_config_v1" {
			continue
		}

		_, err := softwareconfig.Get(orchestrationClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("RTS Software Config still exists")
		}
	}

	return nil
}

func testAccCheckRtsSoftwareConfigV1Exists(n string, configs *softwareconfig.SoftwareConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud orchestration client: %s", err)
		}

		found, err := softwareconfig.Get(orchestrationClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("RTS Software Config not found")
		}

		*configs = *found

		return nil
	}
}

const testAccRtsSoftwareConfigV1_basic = `
resource "huaweicloud_rts_software_config_v1" "config_1" {
  name = "huaweicloud-config"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values=[{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
}
`

const testAccRtsSoftwareConfigV1_timeout = `
resource "huaweicloud_rts_software_config_v1" "config_1" {
  name = "huaweicloud-config"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values=[{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
