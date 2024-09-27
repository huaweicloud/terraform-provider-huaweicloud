package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAuditDataMaskingRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_audit_data_masking_rules.test"
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAuditDataMaskingRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.regex"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.mask_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.operate_time"),
				),
			},
		},
	})
}

func testDataSourceAuditDataMaskingRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dbss_audit_data_masking_rules" "test" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
}
`, testInstance_basic(name))
}
