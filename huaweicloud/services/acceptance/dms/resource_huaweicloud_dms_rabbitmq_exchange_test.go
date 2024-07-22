package dms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
)

func getRabbitmqExchangeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<vhost>/<name>")
	}
	instanceID := parts[0]
	vhost := parts[1]
	name := parts[2]

	return dms.GetRabbitmqExchange(client, instanceID, vhost, name)
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
