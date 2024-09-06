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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRuleAntiCrawlerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}"
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

	getRuleOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getRuleResp, err := client.Request("GET", getPath, &getRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving WAF anti crawler rule: %s", err)
	}
	return utils.FlattenResponse(getRuleResp)
}

func TestAccRuleAntiCrawler_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_anti_crawler.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleAntiCrawlerResourceFunc,
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
				Config: testRuleAntiCrawler_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protection_mode", "anticrawler_specific_url"),
					resource.TestCheckResourceAttr(rName, "priority", "0"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "user-agent"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "TR"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "url"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "equal"),
					resource.TestCheckResourceAttr(rName, "conditions.1.content", "/test/path"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testRuleAntiCrawler_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "priority", "65535"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "user-agent"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "suffix_any"),
					resource.TestCheckResourceAttrPair(rName, "conditions.0.reference_table_id",
						"huaweicloud_waf_reference_table.ref_table", "id"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "user-agent"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "prefix"),
					resource.TestCheckResourceAttr(rName, "conditions.1.content", "RF"),
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

func TestAccRuleAntiCrawler_exceptProtectionMode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_anti_crawler.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleAntiCrawlerResourceFunc,
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
				Config: testRuleAntiCrawler_excepProtectionMode(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protection_mode", "anticrawler_except_url"),
					resource.TestCheckResourceAttr(rName, "priority", "500"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "url"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "not_contain_all"),
					resource.TestCheckResourceAttrPair(rName, "conditions.0.reference_table_id",
						"huaweicloud_waf_reference_table.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testRuleAntiCrawler_excepProtectionMode_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "priority", "1000"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "user-agent"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "TR"),
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

func TestAccRuleAntiCrawler_withEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_anti_crawler.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleAntiCrawlerResourceFunc,
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
				Config: testRuleAntiCrawler_withEpsID(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protection_mode", "anticrawler_except_url"),
					resource.TestCheckResourceAttr(rName, "priority", "500"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "user-agent"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "TR"),
					resource.TestCheckResourceAttrSet(rName, "status"),
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

func testRuleAntiCrawler_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_anti_crawler" "test" {
  policy_id       = huaweicloud_waf_policy.policy_1.id
  name            = "%s"
  protection_mode = "anticrawler_specific_url"
  priority        = 0
  description     = "test description"

  conditions {
    field   = "user-agent"
    logic   = "contain"
    content = "TR"
  }

  conditions {
    field   = "url"
    logic   = "equal"
    content = "/test/path"
  }
}
`, testAccWafPolicyV1_basic(name), name)
}

func testRuleAntiCrawler_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name        = "%s"
  type        = "user-agent"
  description = "test user agent"

  conditions = [
    "UA"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}

resource "huaweicloud_waf_rule_anti_crawler" "test" {
  policy_id       = huaweicloud_waf_policy.policy_1.id
  name            = "%s_update"
  protection_mode = "anticrawler_specific_url"
  priority        = 65535
  description     = "test description update"

  conditions {
    field              = "user-agent"
    logic              = "suffix_any"
    reference_table_id = huaweicloud_waf_reference_table.ref_table.id
  }

  conditions {
    field   = "user-agent"
    logic   = "prefix"
    content = "RF"
  }
}
`, testAccWafPolicyV1_basic(name), name, name)
}

func testRuleAntiCrawler_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_instance" "test" {
  name               = "%[2]s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

resource "huaweicloud_waf_policy" "test" {
  name  = "%[2]s"
  level = 1

  depends_on = [
    huaweicloud_waf_dedicated_instance.test
  ]
}

resource "huaweicloud_waf_reference_table" "test" {
  name        = "%[2]s"
  type        = "url"
  description = "test url"

  conditions = [
    "/test/path"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.test
  ]
}
`, common.TestBaseComputeResources(name), name)
}

func testRuleAntiCrawler_excepProtectionMode(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_anti_crawler" "test" {
  policy_id       = huaweicloud_waf_policy.test.id
  name            = "%[2]s"
  protection_mode = "anticrawler_except_url"
  priority        = 500
  description     = "test description"

  conditions {
    field              = "url"
    logic              = "not_contain_all"
    reference_table_id = huaweicloud_waf_reference_table.test.id
  }
}
`, testRuleAntiCrawler_base(name), name)
}

func testRuleAntiCrawler_excepProtectionMode_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_anti_crawler" "test" {
  policy_id       = huaweicloud_waf_policy.test.id
  name            = "%[2]s_update"
  protection_mode = "anticrawler_except_url"
  priority        = 1000
  description     = "test description update"

  conditions {
    field   = "user-agent"
    logic   = "contain"
    content = "TR"
  }
}
`, testRuleAntiCrawler_base(name), name)
}

func testRuleAntiCrawler_withEpsID(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_anti_crawler" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  name                  = "%s"
  protection_mode       = "anticrawler_except_url"
  priority              = 500
  description           = "test description"
  enterprise_project_id = "%s"

  conditions {
    field   = "user-agent"
    logic   = "contain"
    content = "TR"
  }
}
`, testAccWafPolicyV1_basic_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
