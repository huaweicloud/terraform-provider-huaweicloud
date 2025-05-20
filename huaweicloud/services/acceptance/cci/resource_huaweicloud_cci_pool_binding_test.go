package cci

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

func getPoolBindingResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}

	getPoolBindingHttpUrl := "apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings/{name}"
	getPoolBindingPath := client.Endpoint + getPoolBindingHttpUrl
	getPoolBindingPath = strings.ReplaceAll(getPoolBindingPath, "{namespace}", state.Primary.Attributes["namespace"])
	getPoolBindingPath = strings.ReplaceAll(getPoolBindingPath, "{name}", state.Primary.Attributes["name"])
	getPoolBindingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPoolBindingResp, err := client.Request("GET", getPoolBindingPath, &getPoolBindingOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getPoolBindingResp)
}

func TestAccPoolBinding_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cci_pool_binding.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		obj,
		getPoolBindingResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPoolBinding_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "api_version"),
					resource.TestCheckResourceAttrSet(resourceName, "kind"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "finalizers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPoolBindingStateFunc(resourceName),
			},
		},
	})
}

func testAccPoolBindingStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["namespace"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("the namespace (%s) or name(%s) is nil",
				rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccPoolBinding_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cci_pool_binding" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"
  
  pool_ref {
    id = huaweicloud_elb_pool.test.id
  }
  
  target_ref {
    group = huaweicloud_cciv2_service.test.api_version
    kind  = huaweicloud_cciv2_service.test.kind
    name  = huaweicloud_cciv2_service.test.name
    port  = tolist(huaweicloud_cciv2_service.test.ports)[0].port
  }
}
`, testAccPoolBinding_base(rName), rName)
}

func testAccPoolBinding_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]
}

resource "huaweicloud_elb_pool" "test" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  type        = "instance"
  vpc_id      = huaweicloud_vpc.test.id
  description = "test pool description"
}

resource "huaweicloud_cciv2_service" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  annotations = {
    "kubernetes.io/elb.class"         = "elb",
    "kubernetes.io/elb.id"            = huaweicloud_elb_loadbalancer.test.id,
    "kubernetes.io/elb.protocol-port" = "http:2222",
  }

  ports {
    name         = "service-1"
    protocol     = "TCP"
    port         = 2222
    target_port  = 80
  }

  selector = {
    app = "test1"
  }

  type = "LoadBalancer"

  lifecycle {
    ignore_changes = [
      annotations, ports,
    ]
  }
}
`, testAccV2Network_basic(rName), rName)
}
