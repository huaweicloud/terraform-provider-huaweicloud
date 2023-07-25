package waf

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

func getRuleGeolocationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", state.Primary.Attributes["policy_id"])
	getPath = strings.ReplaceAll(getPath, "{rule_id}", state.Primary.ID)

	queryParam := ""
	if epsID := state.Primary.Attributes["enterprise_project_id"]; epsID != "" {
		queryParam = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	getPath += queryParam

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving WAF geolocation access control rule: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccRuleGeolocation_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_geolocation_access_control.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleGeolocationResourceFunc,
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
				Config: testRuleGeolocation_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "geolocation", "FJ|JL|LN|GZ"),
					resource.TestCheckResourceAttr(rName, "action", "1"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
				),
			},
			{
				Config: testRuleGeolocation_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "geolocation", "FJ|JL|LN|HN"),
					resource.TestCheckResourceAttr(rName, "action", "0"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func TestAccRuleGeolocation_withEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_geolocation_access_control.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleGeolocationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRuleGeolocation_withEpsID(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "geolocation", "FJ|JL|LN|GZ"),
					resource.TestCheckResourceAttr(rName, "action", "1"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func testRuleGeolocation_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_geolocation_access_control" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  name        = "%s"
  geolocation = "FJ|JL|LN|GZ"
  action      = 1
  description = "test description"
}
`, testAccWafPolicyV1_basic(name), name)
}

func testRuleGeolocation_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_geolocation_access_control" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  name        = "%s_update"
  geolocation = "FJ|JL|LN|HN"
  action      = 0
  status      = 0
}
`, testAccWafPolicyV1_basic(name), name)
}

func testRuleGeolocation_withEpsID(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_geolocation_access_control" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  name                  = "%s"
  enterprise_project_id = "%s"
  geolocation           = "FJ|JL|LN|GZ"
  action                = 1
  description           = "test description"
}
`, testAccWafPolicyV1_basic_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
		name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
