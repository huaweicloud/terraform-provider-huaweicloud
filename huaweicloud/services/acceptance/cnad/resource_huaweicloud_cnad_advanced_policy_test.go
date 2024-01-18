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

func getPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                 = acceptance.HW_REGION_NAME
		getPolicyDetailHttpUrl = "v1/cnad/policies/{policy_id}"
		getPolicyDetailProduct = "aad"
	)
	getPolicyDetailClient, err := cfg.NewServiceClient(getPolicyDetailProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CNAD Client: %s", err)
	}

	getPolicyDetailPath := getPolicyDetailClient.Endpoint + getPolicyDetailHttpUrl
	getPolicyDetailPath = strings.ReplaceAll(getPolicyDetailPath, "{policy_id}", state.Primary.ID)

	getPolicyDetailOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getPolicyDetailResp, err := getPolicyDetailClient.Request("GET", getPolicyDetailPath, &getPolicyDetailOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CNAD advanced policy: %s", err)
	}

	return utils.FlattenResponse(getPolicyDetailResp)
}

func TestAccPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cnad_advanced_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPolicyResourceFunc,
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
				Config: testCNADAdvancedPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_cnad_advanced_instances.test",
						"instances.0.instance_id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "threshold", "100"),
					resource.TestCheckResourceAttr(rName, "udp", "block"),
					resource.TestCheckResourceAttr(rName, "block_location.#", "0"),
					resource.TestCheckResourceAttr(rName, "block_protocol.#", "1"),
					resource.TestCheckResourceAttr(rName, "connection_protection_list.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "connection_protection"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint_count"),
					resource.TestCheckResourceAttrSet(rName, "port_block_count"),
					resource.TestCheckResourceAttrSet(rName, "watermark_count"),
				),
			},
			{
				Config: testCNADAdvancedPolicy_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "threshold", "200"),
					resource.TestCheckResourceAttr(rName, "udp", "unblock"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"udp",
				},
			},
		},
	})
}

const testCNADAdvancedPolicy_base = `
data "huaweicloud_cnad_advanced_instances" "test" {}
`

func testCNADAdvancedPolicy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cnad_advanced_policy" "test" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
  name        = "%s"
  threshold   = 100
  udp         = "block"
}
`, testCNADAdvancedPolicy_base, name)
}

func testCNADAdvancedPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cnad_advanced_policy" "test" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
  name        = "%s_update"
  threshold   = 200
  udp         = "unblock"
}
`, testCNADAdvancedPolicy_base, name)
}
