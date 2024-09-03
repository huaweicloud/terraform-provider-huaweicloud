package dms

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

func getRabbitmqQueueResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	getHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues/{queue}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{vhost}", state.Primary.Attributes["vhost"])
	getPath = strings.ReplaceAll(getPath, "{queue}", state.Primary.Attributes["name"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving queue: %s", err)
	}

	return utils.FlattenResponse(getResp)
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
