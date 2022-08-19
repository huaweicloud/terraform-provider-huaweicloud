package antiddos

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chnsz/golangsdk/openstack/aad/v1/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getForwardRuleFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.AadV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CloudTable v2 client: %s", err)
	}
	port := state.Primary.Attributes["forward_port"]
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	return aad.GetForwardRuleFromServer(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["ip"], state.Primary.Attributes["forward_protocol"], portNum)
}

func TestAccForwardRule_basic(t *testing.T) {
	var rule rules.Rule
	resourceName := "huaweicloud_aad_forward_rule.test"
	randomPort := acctest.RandIntRange(1, 65535)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&rule,
		getForwardRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAadForwardRule(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccForwardRule_basic(randomPort),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_AAD_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "ip", acceptance.HW_AAD_IP_ADDRESS),
					resource.TestCheckResourceAttr(resourceName, "forward_protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", strconv.Itoa(randomPort)),
					resource.TestCheckResourceAttr(resourceName, "source_port", strconv.Itoa(randomPort)),
					resource.TestCheckResourceAttr(resourceName, "source_ip", "1.1.1.1,2.2.2.2"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "lb_method"),
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

func testAccForwardRule_basic(port int) string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_forward_rule" "test" {
  instance_id      = "%[1]s"
  ip               = "%[2]s"
  forward_protocol = "udp"
  forward_port     = %[3]d
  source_port      = %[3]d
  source_ip        = "1.1.1.1,2.2.2.2"
}
`, acceptance.HW_AAD_INSTANCE_ID, acceptance.HW_AAD_IP_ADDRESS, port)
}
