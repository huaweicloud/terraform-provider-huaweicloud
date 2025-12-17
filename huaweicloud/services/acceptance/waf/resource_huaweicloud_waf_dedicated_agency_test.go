package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

func getDedicatedAgencyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	return waf.QueryDedicatedAgency(client)
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDedicatedAgency_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_waf_dedicated_agency.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDedicatedAgencyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDedicatedAgency_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrSet(rName, "duration"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "is_valid"),
					resource.TestCheckResourceAttrSet(rName, "role_list.0.description"),
					resource.TestCheckResourceAttrSet(rName, "role_list.0.catalog"),
					resource.TestCheckResourceAttrSet(rName, "role_list.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "role_list.0.is_granted"),
				),
			},
			{
				Config: testDedicatedAgency_basic_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateVerifyIgnore: []string{
					"role_name_list",
					"purged",
				},
			},
		},
	})
}

const testDedicatedAgency_basic = `
resource "huaweicloud_waf_dedicated_agency" "test" {
  role_name_list = ["evs_to_waf_operate_policy", "vpc_to_waf_operate_policy", "ecs_to_waf_operate_policy"]
  purged         = true
}`

const testDedicatedAgency_basic_update = `
resource "huaweicloud_waf_dedicated_agency" "test" {
  role_name_list = ["evs_to_waf_operate_policy"]
  purged         = true
}`
