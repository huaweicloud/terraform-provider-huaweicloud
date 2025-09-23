package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBl7PolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/l7policies/{l7policy_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{l7policy_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
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
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.key", "insert_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value_type", "SYSTEM_DEFINED"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value", "CLIENT-IP"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.remove_headers_config.0.configs.0.key", "remove_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.qps", "100"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.per_source_ip_qps", "10"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.burst", "20"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provisioning_status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.key", "insert_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value_type", "REFERENCE_HEADER"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.remove_headers_config.0.configs.0.key", "remove_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.qps", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.per_source_ip_qps", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.burst", "0"),
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

func TestAccElbV3L7Policy_pools_config(t *testing.T) {
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
				Config: testAccCheckElbV3L7PolicyConfig_pools_config(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_pools_config.0.pool_id",
						"huaweicloud_elb_pool.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "redirect_pools_config.0.weight", "10"),
					resource.TestCheckResourceAttr(resourceName, "redirect_pools_sticky_session_config.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "redirect_pools_sticky_session_config.0.timeout", "100"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.host", "testhost.com"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.path", "/test_path"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.query", "test_query"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.key", "insert_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value_type", "SYSTEM_DEFINED"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value", "CLIENT-IP"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.remove_headers_config.0.configs.0.key", "remove_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.qps", "100"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.per_source_ip_qps", "10"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.burst", "20"),
				),
			},
			{
				Config: testAccCheckElbV3L7PolicyConfig_pools_config_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "redirect_pools_config.0.pool_id",
						"huaweicloud_elb_pool.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "redirect_pools_config.0.weight", "20"),
					resource.TestCheckResourceAttr(resourceName, "redirect_pools_sticky_session_config.0.enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "redirect_pools_sticky_session_config.0.timeout", "200"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.host", "testhostupdate.com"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.path", "/test_path_update"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.rewrite_url_config.0.query", "test_query_update"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.key", "insert_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value_type", "REFERENCE_HEADER"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.remove_headers_config.0.configs.0.key", "remove_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.qps", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.per_source_ip_qps", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_pools_extend_config.0.traffic_limit_config.0.burst", "0"),
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
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.insert_headers_config.0.configs.0.key", "insert_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.insert_headers_config.0.configs.0.value_type", "SYSTEM_DEFINED"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.insert_headers_config.0.configs.0.value", "CLIENT-IP"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.remove_headers_config.0.configs.0.key", "remove_test_111"),
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
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.insert_headers_config.0.configs.0.key", "insert_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.insert_headers_config.0.configs.0.value_type", "REFERENCE_HEADER"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.insert_headers_config.0.configs.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName,
						"redirect_url_config.0.remove_headers_config.0.configs.0.key", "remove_test_222"),
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
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.insert_headers_config.0.configs.0.key", "insert_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.insert_headers_config.0.configs.0.value_type", "SYSTEM_DEFINED"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.insert_headers_config.0.configs.0.value", "CLIENT-IP"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.remove_headers_config.0.configs.0.key", "remove_test_111"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.traffic_limit_config.0.qps", "100"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.traffic_limit_config.0.per_source_ip_qps", "10"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.traffic_limit_config.0.burst", "20"),
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
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.insert_headers_config.0.configs.0.key", "insert_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.insert_headers_config.0.configs.0.value_type", "REFERENCE_HEADER"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.insert_headers_config.0.configs.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.remove_headers_config.0.configs.0.key", "remove_test_222"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.traffic_limit_config.0.qps", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.traffic_limit_config.0.per_source_ip_qps", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"fixed_response_config.0.traffic_limit_config.0.burst", "0"),
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

    insert_headers_config {
      configs {
        key        = "insert_test_111"
        value_type = "SYSTEM_DEFINED"
        value      = "CLIENT-IP"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_111"
      }
    }

    traffic_limit_config{
      qps               = 100
      per_source_ip_qps = 10
      burst             = 20
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

    insert_headers_config {
      configs {
        key        = "insert_test_222"
        value_type = "REFERENCE_HEADER"
        value      = "test_value"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_222"
      }
    }

    traffic_limit_config{
      qps               = 0
      per_source_ip_qps = 0
      burst             = 0
    }
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), updateName)
}

func testAccCheckElbV3L7PolicyConfig_pools_config(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  count = 2

  name            = "%[2]s_${count.index}"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_l7policy" "test" {
  name        = "%[2]s"
  description = "test description"
  action      = "REDIRECT_TO_POOL"
  priority    = 20
  listener_id = huaweicloud_elb_listener.test.id

  redirect_pools_config {
    pool_id = huaweicloud_elb_pool.test[0].id
    weight  = 10
  }

  redirect_pools_sticky_session_config{
    enable  = true
    timeout = 100
  }

  redirect_pools_extend_config {
    rewrite_url_enabled = true

    rewrite_url_config {
      host  = "testhost.com"
      path  = "/test_path"
      query = "test_query"
    }

    insert_headers_config {
      configs {
        key        = "insert_test_111"
        value_type = "SYSTEM_DEFINED"
        value      = "CLIENT-IP"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_111"
      }
    }

    traffic_limit_config{
      qps               = 100
      per_source_ip_qps = 10
      burst             = 20
    }
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), rName)
}

func testAccCheckElbV3L7PolicyConfig_pools_config_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  count = 2

  name            = "%[2]s_${count.index}"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_l7policy" "test" {
  name        = "%[2]s"
  description = "test description update"
  action      = "REDIRECT_TO_POOL"
  priority    = 20
  listener_id = huaweicloud_elb_listener.test.id

  redirect_pools_config {
    pool_id = huaweicloud_elb_pool.test[1].id
    weight  = 20
  }

  redirect_pools_sticky_session_config{
    enable  = false
    timeout = 200
  }

  redirect_pools_extend_config {
    rewrite_url_enabled = true

    rewrite_url_config {
      host  = "testhostupdate.com"
      path  = "/test_path_update"
      query = "test_query_update"
    }

    insert_headers_config {
      configs {
        key        = "insert_test_222"
        value_type = "REFERENCE_HEADER"
        value      = "test_value"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_222"
      }
    }

    traffic_limit_config{
      qps               = 0
      per_source_ip_qps = 0
      burst             = 0
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
  server_certificate          = huaweicloud_elb_certificate.server.id
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
  server_certificate          = huaweicloud_elb_certificate.server.id
}

resource "huaweicloud_elb_listener" "test_redirect_update" {
  name                        = "%[4]s"
  protocol                    = "HTTPS"
  protocol_port               = 448
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  server_certificate          = huaweicloud_elb_certificate.server.id
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

    insert_headers_config {
      configs {
        key        = "insert_test_111"
        value_type = "SYSTEM_DEFINED"
        value      = "CLIENT-IP"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_111"
      }
    }
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

    insert_headers_config {
      configs {
        key        = "insert_test_222"
        value_type = "REFERENCE_HEADER"
        value      = "test_value"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_222"
      }
    }
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

    insert_headers_config {
      configs {
        key        = "insert_test_111"
        value_type = "SYSTEM_DEFINED"
        value      = "CLIENT-IP"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_111"
      }
    }

    traffic_limit_config{
      qps               = 100
      per_source_ip_qps = 10
      burst             = 20
    }
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

    insert_headers_config {
      configs {
        key        = "insert_test_222"
        value_type = "REFERENCE_HEADER"
        value      = "test_value"
      }
    }

    remove_headers_config {
      configs {
        key = "remove_test_222"
      }
    }

    traffic_limit_config{
      qps               = 0
      per_source_ip_qps = 0
      burst             = 0
    }
  }
}
`, testAccCheckElbV3L7PolicyConfig_base(rName), updateName)
}
