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

func getRabbitmqQueueResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return rabbitmq.GetRabbitmqQueue(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["vhost"],
		state.Primary.Attributes["name"])
}

func TestAccRabbitmqQueue_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_queue.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqQueue_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "durable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "vhost", "huaweicloud_dms_rabbitmq_vhost.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_dms_rabbitmq_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "dead_letter_exchange", "amq.direct"),
					resource.TestCheckResourceAttr(resourceName, "dead_letter_routing_key", "binding"),
					resource.TestCheckResourceAttr(resourceName, "message_ttl", "4"),
					resource.TestCheckResourceAttr(resourceName, "lazy_mode", "lazy"),
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

func testRabbitmqQueue_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rabbitmq_queue" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id             = huaweicloud_dms_rabbitmq_instance.test.id
  vhost                   = huaweicloud_dms_rabbitmq_vhost.test.name
  name                    = "%[2]s"
  auto_delete             = false
  durable                 = true
  dead_letter_exchange    = "amq.direct"
  dead_letter_routing_key = "binding"
  message_ttl             = 4
  lazy_mode               = "lazy"
}
`, testRabbitmqVhost_basic(rName), rName)
}

func TestAccRabbitmqQueue_special_charcters(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dms_rabbitmq_queue.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqQueue_special_charcters(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "/test%Queue|-_"),
					resource.TestCheckResourceAttr(resourceName, "vhost", "__F_SLASH__test%25Vhost%7C-_"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_dms_rabbitmq_instance.test", "id"),
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

func testRabbitmqQueue_special_charcters() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rabbitmq_queue" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = urlencode(replace(huaweicloud_dms_rabbitmq_vhost.test.name, "/", "__F_SLASH__"))
  name        = "/test%%Queue|-_"
  auto_delete = false
  durable     = true
}
`, testRabbitmqVhost_special_charcters())
}
