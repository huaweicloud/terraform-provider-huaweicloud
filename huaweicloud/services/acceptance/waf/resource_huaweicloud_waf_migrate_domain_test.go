package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMigrateDomain_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF domain (the client protocol is HTTPS) under a enterprise project.
			// Prepare a WAF certificate and a WAF policy under destination enterprise project.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafDomainId(t)
			acceptance.TestAccPreCheckWafPolicyId(t)
			acceptance.TestAccPreCheckWafCertID(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMigrateDomain_basic(),
			},
		},
	})
}

func testAccMigrateDomain_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_migrate_domain" "test" {
  enterprise_project_id        = "%[1]s"
  target_enterprise_project_id = "%[2]s"
  policy_id                    = "%[3]s"
  host_ids                     = ["%[4]s"]
  certificate_id               = "%[5]s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST,
		acceptance.HW_WAF_POLICY_ID, acceptance.HW_WAF_DOMAIN_ID, acceptance.HW_WAF_CERTIFICATE_ID)
}
