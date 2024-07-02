package apig

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAclPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return acls.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

// All elements' length are same.
// generateRandomStringArray is a method Used to generate the domain names and the domain IDs, and the name cannot start with a digit.
func generateRandomStringArray(count, strLen int) []string {
	if count < 1 || strLen < 1 {
		return nil
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = acctest.RandStringFromCharSet(strLen, "abcdef")
	}
	return result
}

func TestAccAclPolicy_basic(t *testing.T) {
	var (
		policy acls.Policy

		rName1 = "huaweicloud_apig_acl_policy.ip_rule"
		rName2 = "huaweicloud_apig_acl_policy.domain_rule"
		rName3 = "huaweicloud_apig_acl_policy.domain_id_rule"
		name   = acceptance.RandomAccResourceName() // The length is 13.

		basicDomainNames  = strings.Join(generateRandomStringArray(2, 4), ",")
		updateDomainNames = strings.Join(generateRandomStringArray(2, 4), ",")
		basicDomainIds    = strings.Join(generateRandomStringArray(2, 32), ",")
		updateDomainIds   = strings.Join(generateRandomStringArray(2, 32), ",")

		rc1 = acceptance.InitResourceCheck(rName1, &policy, getAclPolicyFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &policy, getAclPolicyFunc)
		rc3 = acceptance.InitResourceCheck(rName3, &policy, getAclPolicyFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigAclPolicy_basic_step1(name, basicDomainNames, basicDomainIds),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_rule_ip"),
					resource.TestCheckResourceAttr(rName1, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName1, "entity_type", "IP"),
					resource.TestCheckResourceAttr(rName1, "value", "10.201.33.4,10.30.2.15"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_rule_domain"),
					resource.TestCheckResourceAttr(rName2, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName2, "entity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(rName2, "value", basicDomainNames),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "name", name+"_rule_domain_id"),
					resource.TestCheckResourceAttr(rName3, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName3, "entity_type", "DOMAIN_ID"),
					resource.TestCheckResourceAttr(rName3, "value", basicDomainIds),
				),
			},
			{
				Config: testAccApigAclPolicy_basic_step2(name, updateDomainNames, updateDomainIds),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_rule_ip_update"),
					resource.TestCheckResourceAttr(rName1, "type", "DENY"),
					resource.TestCheckResourceAttr(rName1, "entity_type", "IP"),
					resource.TestCheckResourceAttr(rName1, "value", "10.201.33.8,10.30.2.23"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_rule_domain_update"),
					resource.TestCheckResourceAttr(rName2, "type", "DENY"),
					resource.TestCheckResourceAttr(rName2, "entity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(rName2, "value", updateDomainNames),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "name", name+"_rule_domain_id_update"),
					resource.TestCheckResourceAttr(rName3, "type", "DENY"),
					resource.TestCheckResourceAttr(rName3, "entity_type", "DOMAIN_ID"),
					resource.TestCheckResourceAttr(rName3, "value", updateDomainIds),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName3),
			},
		},
	})
}

func testAccAclPolicyImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccApigAclPolicy_base() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccApigAclPolicy_basic_step1(name, domainNames, domainIds string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_acl_policy" "ip_rule" {
  instance_id = local.instance_id
  name        = "%[2]s_rule_ip"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}

resource "huaweicloud_apig_acl_policy" "domain_rule" {
  instance_id = local.instance_id
  name        = "%[2]s_rule_domain"
  type        = "PERMIT"
  entity_type = "DOMAIN"
  value       = "%[3]s"
}

resource "huaweicloud_apig_acl_policy" "domain_id_rule" {
  instance_id = local.instance_id
  name        = "%[2]s_rule_domain_id"
  type        = "PERMIT"
  entity_type = "DOMAIN_ID"
  value       = "%[4]s"
}
`, testAccApigAclPolicy_base(), name, domainNames, domainIds)
}

func testAccApigAclPolicy_basic_step2(name, domainNames, domainIds string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_acl_policy" "ip_rule" {
  instance_id = local.instance_id
  name        = "%[2]s_rule_ip_update"
  type        = "DENY"
  entity_type = "IP"
  value       = "10.201.33.8,10.30.2.23"
}

resource "huaweicloud_apig_acl_policy" "domain_rule" {
  instance_id = local.instance_id
  name        = "%[2]s_rule_domain_update"
  type        = "DENY"
  entity_type = "DOMAIN"
  value       = "%[3]s"
}

resource "huaweicloud_apig_acl_policy" "domain_id_rule" {
  instance_id = local.instance_id
  name        = "%[2]s_rule_domain_id_update"
  type        = "DENY"
  entity_type = "DOMAIN_ID"
  value       = "%[4]s"
}
`, testAccApigAclPolicy_base(), name, domainNames, domainIds)
}
