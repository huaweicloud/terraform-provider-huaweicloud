package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getELBl7PolicyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	lbClient, err := c.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	return l7policies.Get(lbClient, state.Primary.ID).Extract()
}

func TestAccElbV3L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_elb_l7policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&l7Policy,
		getELBl7PolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElbV3L7PolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_pool_id",
						"huaweicloud_elb_pool.test", "id"),
				),
			},
			{
				Config: testAccCheckElbV3L7PolicyConfig_basic_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_pool_id",
						"huaweicloud_elb_pool.test_update", "id"),
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

func TestAccElbV3L7Policy_listener(t *testing.T) {
	var l7Policy l7policies.L7Policy
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_elb_l7policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&l7Policy,
		getELBl7PolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElbV3L7PolicyConfig_listener(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_LISTENER"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_listener_id",
						"huaweicloud_elb_listener.test_redirect", "id"),
				),
			},
			{
				Config: testAccCheckElbV3L7PolicyConfig_listener_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_LISTENER"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_listener_id",
						"huaweicloud_elb_listener.test_redirect_update", "id"),
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

func testAccCheckElbV3L7PolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  name            = "%[2]s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_l7policy" "test" {
  name             = "%[2]s"
  description      = "test description"
  action           = "REDIRECT_TO_POOL"
  listener_id      = huaweicloud_elb_listener.test.id
  redirect_pool_id = huaweicloud_elb_pool.test.id
}
`, testAccElbV3ListenerConfig_basic(rName), rName)
}

func testAccCheckElbV3L7PolicyConfig_basic_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  name            = "%[2]s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_pool" "test_update" {
  name            = "%[2]s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_l7policy" "test" {
  name             = "%[2]s"
  description      = "test description update"
  action           = "REDIRECT_TO_POOL"
  listener_id      = huaweicloud_elb_listener.test.id
  redirect_pool_id = huaweicloud_elb_pool.test_update.id
}
`, testAccElbV3ListenerConfig_basic(rName), updateName)
}

func testAccCheckElbV3L7PolicyConfig_listener(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_elb_listener" "test_redirect" {
  name                        = "%[3]s"
  protocol                    = "HTTPS"
  protocol_port               = 443
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  server_certificate          = huaweicloud_elb_certificate.test.id
}

resource "huaweicloud_elb_l7policy" "test" {
  name                 = "%[3]s"
  description          = "test description"
  action               = "REDIRECT_TO_LISTENER"
  listener_id          = huaweicloud_elb_listener.test.id
  redirect_listener_id = huaweicloud_elb_listener.test_redirect.id
}
`, testAccElbV3ListenerConfig_basic(rName), testAccElbV3CertificateConfig_basic(rName), rName)
}

func testAccCheckElbV3L7PolicyConfig_listener_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_elb_listener" "test_redirect" {
  name                        = "%[3]s"
  protocol                    = "HTTPS"
  protocol_port               = 443
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  server_certificate          = huaweicloud_elb_certificate.test.id
}

resource "huaweicloud_elb_listener" "test_redirect_update" {
  name                        = "%[4]s"
  protocol                    = "HTTPS"
  protocol_port               = 448
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  server_certificate          = huaweicloud_elb_certificate.test.id
}

resource "huaweicloud_elb_l7policy" "test" {
  name                 = "%[4]s"
  description          = "test description update"
  action               = "REDIRECT_TO_LISTENER"
  listener_id          = huaweicloud_elb_listener.test.id
  redirect_listener_id = huaweicloud_elb_listener.test_redirect_update.id
}
`, testAccElbV3ListenerConfig_basic(rName), testAccElbV3CertificateConfig_basic(rName), rName, updateName)
}
