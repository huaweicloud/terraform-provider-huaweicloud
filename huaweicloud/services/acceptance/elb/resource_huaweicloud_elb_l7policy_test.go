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
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_pool_id",
						"huaweicloud_elb_pool.test", "id"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.host", "testhost.com"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.path", "/test_path"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.query", "test_query"),
				),
			},
			{
				Config: testAccCheckElbV3L7PolicyConfig_basic_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_pool_id",
						"huaweicloud_elb_pool.test_update", "id"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.host", "testhostupdate.com"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.path", "/test_path_update"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.query", "test_query_update"),
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

func TestAccElbV3L7Policy_url(t *testing.T) {
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
				Config: testAccCheckElbV3L7PolicyConfig_url(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_URL"),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.status_code", "301"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.host", "testhost.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.port", "6666"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.path", "/test_policy"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.query", "abcd"),
				),
			},
			{
				Config: testAccCheckElbV3L7PolicyConfig_url_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_URL"),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.status_code", "302"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.host", "testhostupdate.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.port", "8888"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.path", "/test_policy_update"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url_config.0.query", "efgh"),
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

func TestAccElbV3L7Policy_fixed_response(t *testing.T) {
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
				Config: testAccCheckElbV3L7PolicyConfig_fixed_response(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "FIXED_RESPONSE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "fixed_response_config.0.status_code", "200"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.content_type", "application/json"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.message_body", "it is a test"),
				),
			},
			{
				Config: testAccCheckElbV3L7PolicyConfig_fixed_response_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "FIXED_RESPONSE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "fixed_response_config.0.status_code", "202"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.content_type", "text/css"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.message_body", "it is a test update"),
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

func testAccCheckElbV3L7PolicyConfig_base(rName string) string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_elb_listener" "test" {
  name                        = "%s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  forward_eip                 = true
  idle_timeout                = 62
  request_timeout             = 63
  response_timeout            = 64

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbV3LoadBalancerConfig_basic(rName), rName)
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
  priority         = 20
  listener_id      = huaweicloud_elb_listener.test.id
  redirect_pool_id = huaweicloud_elb_pool.test.id

  redirect_pools_extend_config {
    rewrite_url_enabled = true

    rewrite_url_config {
      host  = "testhost.com"
      path  = "/test_path"
      query = "test_query"
    }
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), rName)
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
  priority         = 50
  listener_id      = huaweicloud_elb_listener.test.id
  redirect_pool_id = huaweicloud_elb_pool.test_update.id

  redirect_pools_extend_config {
    rewrite_url_enabled = true

    rewrite_url_config {
      host  = "testhostupdate.com"
      path  = "/test_path_update"
      query = "test_query_update"
    }
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), updateName)
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
`, testAccCheckElbV3L7PolicyConfig_base(rName), testAccElbV3CertificateConfig_basic(rName), rName)
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
`, testAccCheckElbV3L7PolicyConfig_base(rName), testAccElbV3CertificateConfig_basic(rName), rName, updateName)
}

func testAccCheckElbV3L7PolicyConfig_url(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_l7policy" "test" {
  name        = "%[2]s"
  description = "test description"
  action      = "REDIRECT_TO_URL"
  priority    = 20
  listener_id = huaweicloud_elb_listener.test.id

  redirect_url_config {
    protocol    = "HTTP"
    host        = "testhost.com"
    port        = "6666"
    path        = "/test_policy"
    query       = "abcd"
    status_code = "301"
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), rName)
}

func testAccCheckElbV3L7PolicyConfig_url_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_l7policy" "test" {
  name        = "%[2]s"
  description = "test description update"
  action      = "REDIRECT_TO_URL"
  priority    = 50
  listener_id = huaweicloud_elb_listener.test.id

  redirect_url_config {
    protocol    = "HTTPS"
    host        = "testhostupdate.com"
    port        = "8888"
    path        = "/test_policy_update"
    query       = "efgh"
    status_code = "302"
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), updateName)
}

func testAccCheckElbV3L7PolicyConfig_fixed_response(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_l7policy" "test" {
  name        = "%[2]s"
  description = "test description"
  action      = "FIXED_RESPONSE"
  priority    = 20
  listener_id = huaweicloud_elb_listener.test.id

  fixed_response_config {
    status_code  = "200"
    content_type = "application/json"
    message_body = "it is a test"
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), rName)
}

func testAccCheckElbV3L7PolicyConfig_fixed_response_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_l7policy" "test" {
  name        = "%[2]s"
  description = "test description update"
  action      = "FIXED_RESPONSE"
  priority    = 50
  listener_id = huaweicloud_elb_listener.test.id

  fixed_response_config {
    status_code  = "202"
    content_type = "text/css"
    message_body = "it is a test update"
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), updateName)
}
