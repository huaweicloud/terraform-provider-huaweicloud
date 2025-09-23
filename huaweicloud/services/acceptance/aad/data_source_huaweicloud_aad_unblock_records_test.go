package antiddos

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUnblockRecordsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_unblock_records.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare domain_id before running this test cases.
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUnblockRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.block_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.blocking_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.executor"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.sort_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblock_record.0.unblock_type"),
				),
			},
		},
	})
}

func testUnblockRecords_basic() string {
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf(`
data "huaweicloud_aad_unblock_records" "test" {
  domain_id  = "%[1]s"
  start_time = 0
  end_time   = "%[2]d"
}
`, acceptance.HW_DOMAIN_ID, millis)
}
