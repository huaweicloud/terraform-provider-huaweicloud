package lts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataContextLogs_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_lts_context_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataContextLogs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "logs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.line_num"),
					resource.TestCheckOutput("is_filter_success", "true"),
				),
			},
		},
	})
}

func testAccDataContextLogs_basic(name string) string {
	currentTime := time.Now()
	startTime := currentTime.Format(time.RFC3339)
	endTime := currentTime.Add(1 * time.Hour).Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_logs" "test" {
  depends_on = [null_resource.test]

  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = "%[2]s"
  end_time      = "%[3]s"
}

locals {
  logs = data.huaweicloud_lts_logs.test.logs
}

data "huaweicloud_lts_context_logs" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id

  # Get the line number of the log in the middle of the log list. Verify 'forwards_count' and 'backwards_count'.
  line_num = try(local.logs[floor(length(local.logs) / 2)].line_num, null)
}

output "is_filter_success" {
  value = (
    data.huaweicloud_lts_context_logs.test.forwards_count > 0 &&
    data.huaweicloud_lts_context_logs.test.backwards_count > 0 &&
    data.huaweicloud_lts_context_logs.test.total_count > 0
  )
}
`, testAccDataLogs_base(name), startTime, endTime)
}
