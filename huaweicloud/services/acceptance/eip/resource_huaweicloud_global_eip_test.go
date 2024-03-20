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

func getGEIPResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geip", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GEIP client: %s", err)
	}
	getGEIPHttpUrl := "v3/{domain_id}/global-eips/{id}"
	getGEIPPath := client.Endpoint + getGEIPHttpUrl
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{domain_id}", cfg.DomainID)
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{id}", state.Primary.ID)

	getGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGEIPResp, err := client.Request("GET", getGEIPPath, &getGEIPOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving global EIP: %s", err)
	}
	return utils.FlattenResponse(getGEIPResp)
}

func TestAccGEIP_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_global_eip.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGEIPResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGEIP_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "isp"),
					resource.TestCheckResourceAttrSet(rName, "ip_version"),
					resource.TestCheckResourceAttrSet(rName, "ip_address"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAccGEIP_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func testAccGEIP_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_global_eip" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = "%s"

  tags = {
    foo = "bar"
  }
}
`, testAccInternetBandwidth_basic(name), name)
}

func testAccGEIP_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_global_eip" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = "%s-update"
  description           = "test-update"

  tags = {
    key = "value"
  }
}
`, testAccInternetBandwidth_basic(name), name)
}
