package cnad

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

func getBlackWhiteListResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region               = acceptance.HW_REGION_NAME
		getBlackWhiteHttpUrl = "v1/cnad/policies/{policy_id}"
		getBlackWhiteProduct = "aad"
	)
	getBlackWhiteClient, err := cfg.NewServiceClient(getBlackWhiteProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CNAD Client: %s", err)
	}

	getBlackWhitePath := getBlackWhiteClient.Endpoint + getBlackWhiteHttpUrl
	getBlackWhitePath = strings.ReplaceAll(getBlackWhitePath, "{policy_id}", state.Primary.ID)
	getBlackWhiteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getBlackWhiteResp, err := getBlackWhiteClient.Request("GET", getBlackWhitePath, &getBlackWhiteOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CNAD advanced policy black and white IP list: %s", err)
	}

	respBody, err := utils.FlattenResponse(getBlackWhiteResp)
	if err != nil {
		return nil, err
	}

	blackIpList := utils.PathSearch("pop_policy.bw_list.black_ip_list", respBody, nil)
	whiteIpList := utils.PathSearch("pop_policy.bw_list.white_ip_list", respBody, nil)
	if blackIpList == nil && whiteIpList == nil {
		// the black IP list and white IP list all deleted
		return nil, fmt.Errorf("the black and white IP list of the policy is empty")
	}

	return respBody, nil
}

func TestAccBlackWhiteList_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cnad_advanced_black_white_list.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBlackWhiteListResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCNADInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBlackWhiteList_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_cnad_advanced_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "black_ip_list.#", "2"),
					resource.TestCheckResourceAttr(rName, "white_ip_list.#", "2"),
				),
			},
			{
				Config: testBlackWhiteList_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "black_ip_list.0", "192.168.0.3"),
					resource.TestCheckResourceAttr(rName, "white_ip_list.#", "0"),
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

func testBlackWhiteList_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cnad_advanced_instances" "test" {}

resource "huaweicloud_cnad_advanced_policy" "test" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
  name        = "%s"
  threshold   = 100
  udp         = "block"
}
`, name)
}

func testBlackWhiteList_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cnad_advanced_black_white_list" "test" {
  policy_id     = huaweicloud_cnad_advanced_policy.test.id
  black_ip_list = ["192.168.0.1", "117.78.10.0/14"]
  white_ip_list = ["2001:1234:2234:abcd::1", "2001:1234:2234:abcd::1/64"]
}
`, testBlackWhiteList_base(name))
}

func testBlackWhiteList_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cnad_advanced_black_white_list" "test" {
  policy_id     = huaweicloud_cnad_advanced_policy.test.id
  black_ip_list = ["192.168.0.3"]
  white_ip_list = []
}
`, testBlackWhiteList_base(name))
}
