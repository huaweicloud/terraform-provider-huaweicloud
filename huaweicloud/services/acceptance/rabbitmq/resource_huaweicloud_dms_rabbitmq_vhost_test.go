package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rabbitmq"
)

func getRabbitmqVhostResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return rabbitmq.GetRabbitmqVhost(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["name"])
}

func TestAccRabbitmqVhost_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_vhost.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqVhostResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqVhost_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "tracing"),
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

func testRabbitmqVhost_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_vhost" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = "%s"
}
`, testAccDmsRabbitmqInstance_newFormat_single(rName), rName)
}

func TestAccRabbitmqVhost_special_charcters(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dms_rabbitmq_vhost.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqVhostResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqVhost_special_charcters(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "/test%Vhost|-_"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceVhostImportStateIDFunc(resourceName),
			},
		},
	})
}

func testRabbitmqVhost_special_charcters() string {
	rNameWithDash := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_vhost" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = "/test%%Vhost|-_"
}
`, testAccDmsRabbitmqInstance_newFormat_single(rNameWithDash))
}

func testAccResourceVhostImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		name := rs.Primary.Attributes["name"]
		if instanceID == "" || name == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>,<name>', but got '%s,%s'",
				instanceID, name)
		}
		return fmt.Sprintf("%s,%s", instanceID, name), nil
	}
}
