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

func getRabbitmqExchangeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return rabbitmq.GetRabbitmqExchange(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["vhost"],
		state.Primary.Attributes["name"])
}

func TestAccRabbitmqExchange_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_exchange.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqExchangeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqExchange_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "direct"),
					resource.TestCheckResourceAttr(resourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "durable", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "vhost", "huaweicloud_dms_rabbitmq_vhost.test", "name"),
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

func testRabbitmqExchange_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_exchange" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
  name        = "%s"
  type        = "direct"
  auto_delete = false
  durable     = true
  internal    = false
}
`, testRabbitmqVhost_basic(rName), rName)
}

func TestAccRabbitmqExchange_special_charcters(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	var obj interface{}
	resourceName := "huaweicloud_dms_rabbitmq_exchange.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqExchangeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqExchange_special_charcters(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "/test%Exchange|-_"),
					resource.TestCheckResourceAttr(resourceName, "vhost", "__F_SLASH__test%25Vhost%7C-_"),
					resource.TestCheckResourceAttr(resourceName, "type", "x-delayed-message"),
					resource.TestCheckResourceAttrSet(resourceName, "arguments"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceExchangeOrQueueImportStateIDFunc(resourceName),
			},
		},
	})
}

func testRabbitmqExchange_special_charcters(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_vhost" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = "/test%%Vhost|-_"
}

resource "huaweicloud_dms_rabbitmq_exchange" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = urlencode(replace(huaweicloud_dms_rabbitmq_vhost.test.name, "/", "__F_SLASH__"))
  name        = "/test%%Exchange|-_"
  type        = "x-delayed-message"
  auto_delete = false

  arguments   = jsonencode({
    "x-delayed-type" = "header"
  })
}
`, testAccDmsRabbitmqInstance_amqp_single(name, false))
}

func testAccResourceExchangeOrQueueImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		vhost := rs.Primary.Attributes["vhost"]
		name := rs.Primary.Attributes["name"]
		if instanceID == "" || vhost == "" || name == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>,<vhost>,<name>', but got '%s,%s,%s'",
				instanceID, vhost, name)
		}
		return fmt.Sprintf("%s,%s,%s", instanceID, vhost, name), nil
	}
}

func TestAccRabbitmqExchange_bindings(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dms_rabbitmq_exchange.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqExchangeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqExchangeAssociate_special_characters(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "bindings.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceExchangeOrQueueImportStateIDFunc(resourceName),
			},
		},
	})
}
