package rabbitmq

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRabbitmqExchangeAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	getHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{vhost}", state.Primary.Attributes["vhost"])
	getPath = strings.ReplaceAll(getPath, "{exchange}", state.Primary.Attributes["exchange"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the exchange association infos: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the exchanges association infos: %s", err)
	}

	searchPath := fmt.Sprintf("items[?destination_type=='%s']|[?destination=='%s']",
		strings.ToLower(state.Primary.Attributes["destination_type"]), state.Primary.Attributes["destination"])
	associations := utils.PathSearch(searchPath, getRespBody, make([]interface{}, 0)).([]interface{})

	routingKey := state.Primary.Attributes["routing_key"]

	for _, association := range associations {
		if routingKey == utils.PathSearch("routing_key", association, "").(string) {
			return association, nil
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccRabbitmqExchangeAssociate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_exchange_associate.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqExchangeAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqExchangeAssociate_basic(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "properties_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceExchangeAssociateImportStateIDFunc(resourceName),
			},
			{
				Config: testRabbitmqExchangeAssociate_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "properties_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceExchangeAssociateImportStateIDFunc(resourceName),
			},
		},
	})
}

func testRabbitmqExchangeAssociate_basic(rName, routingKey string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_dms_rabbitmq_exchange_associate" "test" {
  depends_on = [
    huaweicloud_dms_rabbitmq_vhost.test,
    huaweicloud_dms_rabbitmq_exchange.test,
    huaweicloud_dms_rabbitmq_queue.test
  ]

  instance_id      = huaweicloud_dms_rabbitmq_instance.test.id
  vhost            = huaweicloud_dms_rabbitmq_vhost.test.name
  exchange         = huaweicloud_dms_rabbitmq_exchange.test.name
  destination_type = "Queue"
  destination      = huaweicloud_dms_rabbitmq_queue.test.name
  routing_key      = "%[3]s"
}
`, testRabbitmqVhost_basic(rName), testRabbitmqExchangeAssociate_base_queue_and_exchange(rName), routingKey)
}

func testRabbitmqExchangeAssociate_base_queue_and_exchange(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_rabbitmq_exchange" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
  name        = "%[1]s"
  type        = "direct"
  auto_delete = false
}

resource "huaweicloud_dms_rabbitmq_queue" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
  name        = "%[1]s"
  auto_delete = false
}
`, rName)
}

func testAccResourceExchangeAssociateImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		vhost := rs.Primary.Attributes["vhost"]
		exchange := rs.Primary.Attributes["exchange"]
		destinationType := rs.Primary.Attributes["destination_type"]
		destination := rs.Primary.Attributes["destination"]
		routingKey := rs.Primary.Attributes["routing_key"]

		if instanceID == "" || vhost == "" || exchange == "" || destinationType == "" || destination == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<instance_id>,<vhost>,<exchange>,<destination_type>,<destination>,', but got '%s,%s,%s,%s,%s'",
				instanceID, vhost, exchange, destinationType, destination)
		}

		if routingKey == "" {
			return fmt.Sprintf("%s,%s,%s,%s,%s", instanceID, vhost, exchange, destinationType, destination), nil
		}

		return fmt.Sprintf("%s,%s,%s,%s,%s,%s", instanceID, vhost, exchange, destinationType, destination, routingKey), nil
	}
}

func TestAccRabbitmqExchangeAssociate_special_characters(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_exchange_associate.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqExchangeAssociateResourceFunc,
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
					resource.TestCheckResourceAttrSet(resourceName, "properties_key"),
					resource.TestCheckResourceAttr(resourceName, "routing_key", "/test%encode|\\"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceExchangeAssociateImportStateIDFunc(resourceName),
			},
		},
	})
}

func testRabbitmqExchangeAssociate_special_characters(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  special_characters_name = "/test%%special|characters_-"
}

resource "huaweicloud_dms_rabbitmq_vhost" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = local.special_characters_name
}

resource "huaweicloud_dms_rabbitmq_exchange" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = urlencode(replace(huaweicloud_dms_rabbitmq_vhost.test.name, "/", "__F_SLASH__"))
  name        = local.special_characters_name
  type        = "direct"
  auto_delete = false
}

resource "huaweicloud_dms_rabbitmq_queue" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = urlencode(replace(huaweicloud_dms_rabbitmq_vhost.test.name, "/", "__F_SLASH__"))
  name        = local.special_characters_name
  auto_delete = false
}

resource "huaweicloud_dms_rabbitmq_exchange_associate" "test" {
  depends_on = [
    huaweicloud_dms_rabbitmq_vhost.test,
    huaweicloud_dms_rabbitmq_exchange.test,
    huaweicloud_dms_rabbitmq_queue.test
  ]

  instance_id      = huaweicloud_dms_rabbitmq_instance.test.id
  vhost            = urlencode(replace(huaweicloud_dms_rabbitmq_vhost.test.name, "/", "__F_SLASH__"))
  exchange         = urlencode(replace(huaweicloud_dms_rabbitmq_exchange.test.name, "/", "__F_SLASH__"))
  destination_type = "queue"
  destination      = huaweicloud_dms_rabbitmq_queue.test.name
  routing_key      = "/test%%encode|\\"
}
`, testAccDmsRabbitmqInstance_newFormat_single(rName))
}
