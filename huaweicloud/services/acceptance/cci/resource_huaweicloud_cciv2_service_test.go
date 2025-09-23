package cci

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getV2ServiceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}

	getServiceHttpUrl := "apis/cci/v2/namespaces/{namespace}/services/{name}"
	getServicePath := client.Endpoint + getServiceHttpUrl
	getServicePath = strings.ReplaceAll(getServicePath, "{namespace}", state.Primary.Attributes["namespace"])
	getServicePath = strings.ReplaceAll(getServicePath, "{name}", state.Primary.Attributes["name"])
	getServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getServiceResp, err := client.Request("GET", getServicePath, &getServiceOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getServiceResp)
}

func TestAccV2Service_basic(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_service.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getV2ServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Service_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "Service"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttr(resourceName, "ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "type", "LoadBalancer"),
					resource.TestCheckResourceAttr(resourceName, "ports.0.port", "86"),
					resource.TestCheckResourceAttr(resourceName, "ports.0.target_port", "65530"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.app", "test1"),
				),
			},
			{
				Config: testAccV2Service_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ports.0.port", "87"),
					resource.TestCheckResourceAttr(resourceName, "ports.0.target_port", "65529"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.app", "test2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func testAccV2Service_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_service" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  annotations = {
    "kubernetes.io/elb.class" = "elb",
    "kubernetes.io/elb.id"    = huaweicloud_elb_loadbalancer.test.id,
  }

  ports {
    name         = "service-%[2]s-port"
    app_protocol = "TCP"
    protocol     = "TCP"
    port         = 86
    target_port  = 65530
  }

  selector = {
    app = "test1"
  }

  type = "LoadBalancer"

  lifecycle {
    ignore_changes = [
      annotations,
    ]
  }
}
`, testAccV2ServiceBase(rName), rName)
}

func testAccV2Service_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_service" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  annotations = {
    "kubernetes.io/elb.class" = "elb",
    "kubernetes.io/elb.id"    = huaweicloud_elb_loadbalancer.test.id,
  }

  ports {
    name         = "service-%[2]s-port"
    app_protocol = "TCP"
    protocol     = "TCP"
    port         = 87
    target_port  = 65529
  }

  selector = {
    app = "test2"
  }

  type = "LoadBalancer"

  lifecycle {
    ignore_changes = [
      annotations,
    ]
  }
}
`, testAccV2ServiceBase(rName), rName)
}

func testAccV2ServiceBase(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[3]s"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]
}
`, common.TestBaseNetwork(rName), testAccV2Namespace_basic(rName), rName)
}
