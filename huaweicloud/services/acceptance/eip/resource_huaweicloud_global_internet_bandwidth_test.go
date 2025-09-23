package eip

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

func getInternetBandwidthResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geip", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GEIP client: %s", err)
	}
	getInternetBandwidthHttpUrl := "v3/{domain_id}/geip/internet-bandwidths/{id}"
	getInternetBandwidthPath := client.Endpoint + getInternetBandwidthHttpUrl
	getInternetBandwidthPath = strings.ReplaceAll(getInternetBandwidthPath, "{domain_id}", cfg.DomainID)
	getInternetBandwidthPath = strings.ReplaceAll(getInternetBandwidthPath, "{id}", state.Primary.ID)

	getInternetBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getInternetBandwidthResp, err := client.Request("GET", getInternetBandwidthPath, &getInternetBandwidthOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving global internet bandwidth: %s", err)
	}
	return utils.FlattenResponse(getInternetBandwidthResp)
}

func TestAccInternetBandwidth_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_global_internet_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInternetBandwidthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInternetBandwidth_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "charge_mode", "95peak_guar"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "size", "300"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "ratio_95peak"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAccInternetBandwidth_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "size", "400"),
					resource.TestCheckResourceAttr(rName, "description", "test-update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
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

func testAccInternetBandwidth_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_global_eip_pools" "all" {
  access_site = "cn-north-beijing"
  ip_version  = 4
}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode           = "95peak_guar"
  size                  = 300
  isp                   = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name                  = "%s"
  type                  = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type

  tags = {
    foo = "bar"
  }
}
`, name)
}

func testAccInternetBandwidth_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_global_eip_pools" "all" {
  access_site = "cn-north-beijing"
  ip_version  = 4
}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode           = "95peak_guar"
  size                  = 400
  isp                   = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name                  = "%s-update"
  description           = "test-update"
  type                  = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type

  tags = {
    key = "value"
  }
}
`, name)
}
