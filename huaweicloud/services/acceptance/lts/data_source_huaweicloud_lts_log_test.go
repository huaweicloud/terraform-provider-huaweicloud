package lts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataLogs_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_lts_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLabels   = "data.huaweicloud_lts_logs.filter_by_labels"
		dcByLabels = acceptance.InitDataSourceCheck(byLabels)

		byIsDesc   = "data.huaweicloud_lts_logs.filter_by_is_desc"
		dcByIsDesc = acceptance.InitDataSourceCheck(byIsDesc)

		byKeywordsHighlight   = "data.huaweicloud_lts_logs.filter_by_keywords_highlight"
		dcByKeywordsHighlight = acceptance.InitDataSourceCheck(byKeywordsHighlight)

		randomId, _ = uuid.GenerateUUID()
		currentTime = time.Now()
		startTime   = currentTime.Format(time.RFC3339)
		endTime     = currentTime.Add(1 * time.Hour).Format(time.RFC3339)
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
				Config:      testAccDataLogs_logGroupNotFound(randomId, startTime, endTime),
				ExpectError: regexp.MustCompile(`The log group does not existed`),
			},
			{
				Config:      testAccDataLogs_logStreamNotFound(name, randomId, startTime, endTime),
				ExpectError: regexp.MustCompile(`The log stream does not existed`),
			},
			{
				Config: testAccDataLogs_basic(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "logs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.line_num"),
					dcByLabels.CheckResourceExists(),
					resource.TestCheckOutput("is_labels_filter_useful", "true"),
					dcByIsDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_desc_filter_useful", "true"),
					dcByKeywordsHighlight.CheckResourceExists(),
					resource.TestCheckOutput("is_keywords_highlight_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataLogs_logGroupNotFound(randomId, startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_lts_logs" "test" {
  log_group_id  = "%[1]s"
  log_stream_id = "%[1]s"
  start_time    = "%[2]s"
  end_time      = "%[3]s"
}
`, randomId, startTime, endTime)
}

func testAccDataLogs_logStreamNotFound(randomId, name, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}

data "huaweicloud_lts_logs" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = "%[2]s"
  start_time    = "%[3]s"
  end_time      = "%[4]s"
}
`, name, randomId, startTime, endTime)
}

func testAccDataLogs_base(name string) string {
	return fmt.Sprintf(`
variable "script_content" {
  type    = string
  default = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[1]s"
  memory_size           = 128
  runtime               = "Python3.9"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  functiongraph_version = "v2"
  agency                = "%[2]s"
  enable_lts_log        = true
  log_group_id          = huaweicloud_lts_group.test.id
  log_group_name        = huaweicloud_lts_group.test.group_name
  log_stream_id         = huaweicloud_lts_stream.test.id
  log_stream_name       = huaweicloud_lts_stream.test.stream_name
}

resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "ACTIVE"

  event_data = jsonencode({
    name           = "%[1]s"
    schedule_type  = "Rate"
    sync_execution = false
    schedule       = "1m"
  })
}

# Waiting for the LTS log to be generated.
resource "null_resource" "test" {
  depends_on = [huaweicloud_fgs_function_trigger.test]

  provisioner "local-exec" {
    command = "sleep 180;"
  }
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccDataLogs_basic(name, startTime, endTime string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_logs" "test" {
  depends_on = [null_resource.test]

  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = "%[2]s"
  end_time      = "%[3]s"
}

# Filter by labels.
data "huaweicloud_lts_logs" "filter_by_labels" {
  depends_on = [null_resource.test]

  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = "%[2]s"
  end_time      = "%[3]s"

  labels = {
    function = huaweicloud_fgs_function.test.name
  }
}

locals {
  lables_filter_result = [for v in flatten(data.huaweicloud_lts_logs.filter_by_labels.logs[*].labels[*].function) :
  strcontains(v, huaweicloud_fgs_function.test.name)]
}

output "is_labels_filter_useful" {
  value = length(local.lables_filter_result) > 0 && alltrue(local.lables_filter_result)
}

# Filter by descending order.
data "huaweicloud_lts_logs" "filter_by_is_desc" {
  depends_on = [null_resource.test]

  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = "%[2]s"
  end_time      = "%[3]s"
  is_desc       = true
}

# Filter by ascending order.
data "huaweicloud_lts_logs" "filter_by_is_asc" {
  depends_on = [null_resource.test]

  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = "%[2]s"
  end_time      = "%[3]s"
  is_desc       = false
}

locals {
  sort_desc_filter_result = data.huaweicloud_lts_logs.filter_by_is_desc.logs[*].line_num
  sort_asc_filter_result  = data.huaweicloud_lts_logs.filter_by_is_asc.logs[*].line_num
}

output "is_desc_filter_useful" {
  value = (
    length(local.sort_desc_filter_result) == length(local.sort_asc_filter_result) &&
    length(local.sort_asc_filter_result) > 0 &&
    local.sort_desc_filter_result[0] == element(local.sort_asc_filter_result, -1)
  )
}

# Filter by highlight.
locals {
  keywords = try(split(" ", data.huaweicloud_lts_logs.test.logs[0].content)[1], null)
}

data "huaweicloud_lts_logs" "filter_by_keywords_highlight" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = "%[2]s"
  end_time      = "%[3]s"
  keywords      = local.keywords
  highlight     = true
}

locals {
  keywords_highlight_flter_result = [for v in data.huaweicloud_lts_logs.filter_by_keywords_highlight.logs[*].content :
  strcontains(v, format("<HighLightTag>%%s</HighLightTag>", local.keywords))]
}

output "is_keywords_highlight_filter_useful" {
  value = length(local.keywords_highlight_flter_result) > 0 && alltrue(local.keywords_highlight_flter_result)
}
`, testAccDataLogs_base(name), startTime, endTime)
}
