package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsAuditLogLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_audit_log_links.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
			acceptance.TestAccPreCheckDDSTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDdsAuditLogLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDdsAuditLogLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_audit_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

data "huaweicloud_dds_audit_log_links" "test" {
  instance_id = "%[1]s"
  ids         = data.huaweicloud_dds_audit_logs.test.audit_logs.*.id
}
`, acceptance.HW_DDS_INSTANCE_ID, acceptance.HW_DDS_START_TIME, acceptance.HW_DDS_END_TIME)
}
