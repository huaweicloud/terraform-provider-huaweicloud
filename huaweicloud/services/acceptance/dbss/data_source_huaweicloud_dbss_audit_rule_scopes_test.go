package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAuditRuleScopes_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_audit_rule_scopes.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDbssInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAuditRuleScopes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "scopes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scopes.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scopes.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scopes.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scopes.0.all_audit"),
				),
			},
		},
	})
}

func testDataSourceAuditRuleScopes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dbss_audit_rule_scopes" "test" {
  instance_id = "%s"
}
`, acceptance.HW_DBSS_INSATNCE_ID)
}
