package antiddos

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPolicyBlackWhiteRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region       = acceptance.HW_REGION_NAME
		product      = "aad"
		domainName   = state.Primary.Attributes["domain_name"]
		overseasType = state.Primary.Attributes["overseas_type"]
		ip           = state.Primary.Attributes["ip"]
		ruleType     = state.Primary.Attributes["type"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AAD client: %s", err)
	}

	return aad.GetPolicyBlackWhiteRule(client, domainName, convertStringtoInt(overseasType), ip, convertStringtoInt(ruleType))
}

func convertStringtoInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("[ERROR] convert the string %s to int failed.", s)
	}
	return i
}

func TestAccPolicyBlackWhiteRule_basic(t *testing.T) {
	var obj interface{}
	rscName := "huaweicloud_aad_policy_black_white_rule.test"
	rc := acceptance.InitResourceCheck(
		rscName,
		&obj,
		getPolicyBlackWhiteRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckAadDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPolicyBlackWhiteRule(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rscName, "domain_name", acceptance.HW_AAD_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rscName, "type", "1"),
					resource.TestCheckResourceAttr(rscName, "ip", "192.168.5.10"),
					resource.TestCheckResourceAttr(rscName, "overseas_type", "0"),
					resource.TestCheckResourceAttrSet(rscName, "domain_id"),
				),
			},
			{
				ResourceName:      rscName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPolicyBlackWhiteRuleImportState(rscName),
			},
		},
	})
}

func testPolicyBlackWhiteRule() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_policy_black_white_rule" "test" {
  domain_name   = "%s"
  ip            = "192.168.5.10"
  overseas_type = 0
  type          = 1
}
`, acceptance.HW_AAD_DOMAIN_NAME)
}

func testPolicyBlackWhiteRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		domainName := rs.Primary.Attributes["domain_name"]
		if domainName == "" {
			return "", fmt.Errorf("attribute (domain_name) of Resource (%s) not found", name)
		}

		overseasType := rs.Primary.Attributes["overseas_type"]
		if overseasType == "" {
			return "", fmt.Errorf("attribute (overseas_type) of Resource (%s) not found", name)
		}

		ip := rs.Primary.Attributes["ip"]
		if ip == "" {
			return "", fmt.Errorf("attribute (ip) of Resource (%s) not found", name)
		}

		ruleType := rs.Primary.Attributes["type"]
		if ruleType == "" {
			return "", fmt.Errorf("attribute (type) of Resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s/%s/%s", domainName, overseasType, ip, ruleType), nil
	}
}
