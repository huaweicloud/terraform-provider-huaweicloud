package cfw

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getProtectionRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getProtectionRule: Query the CFW Protection Rule detail
	var (
		getProtectionRuleHttpUrl = "v1/{project_id}/acl-rules"
		getProtectionRuleProduct = "cfw"
	)
	getProtectionRuleClient, err := conf.NewServiceClient(getProtectionRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ProtectionRule Client: %s", err)
	}

	getProtectionRulePath := getProtectionRuleClient.Endpoint + getProtectionRuleHttpUrl
	getProtectionRulePath = strings.ReplaceAll(getProtectionRulePath, "{project_id}", getProtectionRuleClient.ProjectID)

	getProtectionRulequeryParams := buildGetProtectionRuleQueryParams(state)
	getProtectionRulePath += getProtectionRulequeryParams

	getPotectionRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getProtectionRuleResp, err := getProtectionRuleClient.Request("GET", getProtectionRulePath, &getPotectionRulesOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving protection rule: %s", err)
	}

	getProtectionRuleRespBody, err := utils.FlattenResponse(getProtectionRuleResp)
	if err != nil {
		return nil, err
	}

	rules, err := jmespath.Search("data.records", getProtectionRuleRespBody)
	if err != nil {
		diag.Errorf("error parsing data.records from response= %#v", getProtectionRuleRespBody)
	}

	return cfw.FilterRules(rules.([]interface{}), state.Primary.ID)
}

func TestAccProtectionRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_protection_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectionRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProtectionRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
					resource.TestCheckResourceAttr(rName, "action_type", "0"),
					resource.TestCheckResourceAttr(rName, "long_connect_enable", "0"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.address", "1.1.1.2"),
					resource.TestCheckResourceAttrSet(rName, "rule_hit_count"),
				),
			},
			{
				Config: testProtectionRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(rName, "action_type", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "2.2.2.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.address", "2.2.2.2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testProtectionRuleImportState(rName),
				ImportStateVerifyIgnore: []string{
					"sequence", "type",
				},
			},
		},
	})
}

func testProtectionRule_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type    = 0
    address = "1.1.1.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testProtectionRule_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%s-update"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test update"
  type                = 0
  address_type        = 0
  action_type         = 1
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "2.2.2.1"
  }

  destination {
    type    = 0
    address = "2.2.2.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func buildGetProtectionRuleQueryParams(state *terraform.ResourceState) string {
	res := "?offset=0&limit=1024"
	res = fmt.Sprintf("%s&object_id=%v", res, state.Primary.Attributes["object_id"])

	return res
}

func testProtectionRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["object_id"] == "" {
			return "", fmt.Errorf("Attribute (object_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["object_id"] + "/" +
			rs.Primary.ID, nil
	}
}
