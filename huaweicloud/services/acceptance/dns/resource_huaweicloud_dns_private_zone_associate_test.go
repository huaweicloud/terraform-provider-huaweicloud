package dns

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDNSPrivateZoneAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dns_region", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}

	httpUrl := "v2/zones/{zone_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{zone_id}", state.Primary.Attributes["zone_id"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("routers[?router_id=='%s']|[0]", state.Primary.Attributes["router_id"])
	router := utils.PathSearch(searchPath, getRespBody, nil)
	if router == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return router, nil
}

func TestAccDNSPrivateZoneAssociate_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = fmt.Sprintf("acpttest-zone-%s.com", acctest.RandString(5))
		rName = "huaweicloud_dns_private_zone_associate.test"
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSPrivateZoneAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSPrivateZoneAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "zone_id", "huaweicloud_dns_zone.private", "id"),
					resource.TestCheckResourceAttrPair(rName, "router_id", "huaweicloud_vpc.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDNSPrivateZoneAssociate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count = 2

  name = "%[1]s_${count.index}"
  cidr = cidrsubnet("192.168.0.0/16", 4, count.index)
}

resource "huaweicloud_dns_zone" "private" {
  name        = "%[1]s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"
  status      = "ENABLE"

  router {
    router_id = huaweicloud_vpc.test[0].id
  }

  proxy_pattern = "RECURSIVE"

  lifecycle {
    ignore_changes = [router]
  }
}

resource "huaweicloud_dns_private_zone_associate" "test" {
  zone_id   = huaweicloud_dns_zone.private.id
  router_id = huaweicloud_vpc.test[1].id
}`, rName)
}
