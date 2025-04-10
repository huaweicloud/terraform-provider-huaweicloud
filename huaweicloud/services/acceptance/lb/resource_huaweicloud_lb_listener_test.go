package lb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v2/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getL7ListenerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/elb/listeners/{listener_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{listener_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccLBV2Listener_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_listener.listener_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getL7ListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2ListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "connection_limit", "-1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccLBV2ListenerConfig_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

func TestAccLBV2Listener_https(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_listener.listener_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&listener,
		getL7ListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2ListenerConfig_https(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "http2_enable", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"huaweicloud_lb_loadbalancer.loadbalancer_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "default_tls_container_ref",
						"huaweicloud_lb_certificate.certificate_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "client_ca_tls_container_ref",
						"huaweicloud_lb_certificate.certificate_client", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "sni_container_refs.0",
						"huaweicloud_lb_certificate.certificate_1", "id"),
					resource.TestCheckResourceAttr(resourceName, "tls_ciphers_policy", "tls-1-1"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "insert_headers.0.x_forwarded_elb_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "insert_headers.0.x_forwarded_host", "true"),
				),
			},
			{
				Config: testAccLBV2ListenerConfig_https_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "http2_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls_ciphers_policy", "tls-1-2"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "insert_headers.0.x_forwarded_elb_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "insert_headers.0.x_forwarded_host", "false"),
				),
			},
		},
	})
}

func testAccLBV2ListenerConfig_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  description   = "created by acceptance test"
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, common.TestVpc(rName), rName)
}

func testAccLBV2ListenerConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  description     = "created by acceptance test"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccLBV2ListenerConfig_base(rName), rName)
}

func testAccLBV2ListenerConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id

  tags = {
    foo   = "bar"
    owner = "terraform_update"
  }
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName), rNameUpdate)
}

func testAccLBV2ListenerConfig_https(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_lb_listener" "listener_1" {
  name                        = "%[4]s"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  loadbalancer_id             = huaweicloud_lb_loadbalancer.loadbalancer_1.id
  default_tls_container_ref   = huaweicloud_lb_certificate.certificate_1.id
  client_ca_tls_container_ref = huaweicloud_lb_certificate.certificate_client.id
  sni_container_refs          = [huaweicloud_lb_certificate.certificate_1.id]
  tls_ciphers_policy          = "tls-1-1"
  protection_status           = "consoleProtection"
  protection_reason           = "test protection reason"

  insert_headers {
    x_forwarded_elb_ip = true
    x_forwarded_host   = true
  }
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName),
		testAccLBV2CertificateConfig_client(rName), rName)
}

func testAccLBV2ListenerConfig_https_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_lb_listener" "listener_1" {
  name                        = "%[4]s"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  loadbalancer_id             = huaweicloud_lb_loadbalancer.loadbalancer_1.id
  default_tls_container_ref   = huaweicloud_lb_certificate.certificate_1.id
  client_ca_tls_container_ref = huaweicloud_lb_certificate.certificate_client.id
  sni_container_refs          = [huaweicloud_lb_certificate.certificate_1.id]
  http2_enable                = true
  tls_ciphers_policy          = "tls-1-2"
  protection_status           = "nonProtection"

  insert_headers {
    x_forwarded_elb_ip = false
    x_forwarded_host   = false
  }
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName),
		testAccLBV2CertificateConfig_client(rName), rNameUpdate)
}
