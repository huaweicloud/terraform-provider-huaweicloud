package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getInstanceIngressPortFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	ingressPortId := state.Primary.ID

	return apig.GetInstanceIngressPortById(client, instanceId, ingressPortId)
}

func TestAccInstanceIngressPort_basic(t *testing.T) {
	var (
		ingressPort interface{}

		protocolHTTP    = "huaweicloud_apig_instance_ingress_port.http"
		rcProtocolHTTP  = acceptance.InitResourceCheck(protocolHTTP, &ingressPort, getInstanceIngressPortFunc)
		protocolHTTPS   = "huaweicloud_apig_instance_ingress_port.https"
		rcProtocolHTTPS = acceptance.InitResourceCheck(protocolHTTPS, &ingressPort, getInstanceIngressPortFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcProtocolHTTP.CheckResourceDestroy(),
			rcProtocolHTTPS.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceIngressPort_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rcProtocolHTTP.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolHTTP, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(protocolHTTP, "port", "8080"),
					resource.TestCheckResourceAttrSet(protocolHTTP, "status"),
				),
			},
			{
				ResourceName:      protocolHTTP,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceIngressPortIDImportStateFunc(protocolHTTP),
			},
			{
				ResourceName:      protocolHTTP,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceIngressPortProtocolPortImportStateFunc(protocolHTTP),
			},
			{
				Config: testAccInstanceIngressPort_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rcProtocolHTTPS.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolHTTPS, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(protocolHTTPS, "port", "8443"),
					resource.TestCheckResourceAttrSet(protocolHTTPS, "status"),
				),
			},
			{
				ResourceName:      protocolHTTPS,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceIngressPortIDImportStateFunc(protocolHTTPS),
			},
			{
				ResourceName:      protocolHTTPS,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceIngressPortProtocolPortImportStateFunc(protocolHTTPS),
			},
		},
	})
}

func testAccInstanceIngressPortIDImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		ingressPortId := rs.Primary.ID
		if instanceId == "" || ingressPortId == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>', but '%s/%s'",
				instanceId, ingressPortId)
		}
		return fmt.Sprintf("%s/%s", instanceId, ingressPortId), nil
	}
}

func testAccInstanceIngressPortProtocolPortImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		protocol := rs.Primary.Attributes["protocol"]
		port := rs.Primary.Attributes["port"]
		if instanceId == "" || protocol == "" || port == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<protocol>/<port>', but '%s/%s/%s'",
				instanceId, protocol, port)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, protocol, port), nil
	}
}

func testAccInstanceIngressPort_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_instance_ingress_port" "http" {
  instance_id = local.instance_id
  protocol    = "HTTP"
  port        = 8080
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccInstanceIngressPort_basic_step2() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_instance_ingress_port" "https" {
  instance_id = local.instance_id
  protocol    = "HTTPS"
  port        = 8443
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
