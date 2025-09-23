package cts

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCtsTraces_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cts_traces.filter_by_resource_type"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	currentTime := time.Now().UTC()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")

	laterTime := currentTime.Add(3 * time.Minute)
	laterTimeString := laterTime.Format("2006-01-02 15:04:05")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCtsTraces_basic(rName, currentTimeString, laterTimeString),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.trace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.trace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.trace_rating"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.trace_type"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.record_time"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.operation_id"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
					resource.TestCheckOutput("is_trace_rating_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCtsTraces_basic(name, from, to string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id            = huaweicloud_cts_data_tracker.tracker.id
  resource_type = "tracker"
  trace_rating  = "normal"
}

data "huaweicloud_cts_traces" "test" {
  trace_type = "system"
  from       = "%[2]s"
  to         = "%[3]s"

  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cts_traces.test.traces) >= 1
}

data "huaweicloud_cts_traces" "filter_by_id" {
  trace_type  = "system"
  from        = "%[2]s"
  to          = "%[3]s"
  resource_id = local.id
  
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cts_traces.filter_by_id.traces) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_traces.filter_by_id.traces[*] : v.resource_id == local.id]
  )
}

data "huaweicloud_cts_traces" "filter_by_resource_type" {
  trace_type    = "system"
  from          = "%[2]s"
  to            = "%[3]s"
  resource_type = local.resource_type
  
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_resource_type_filter_useful" {
  value = length(data.huaweicloud_cts_traces.filter_by_resource_type.traces) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_traces.filter_by_resource_type.traces[*] : v.resource_type == local.resource_type]
  )
}

data "huaweicloud_cts_traces" "filter_by_trace_rating" {
  trace_type    = "system"
  from          = "%[2]s"
  to            = "%[3]s"
  trace_rating  = local.trace_rating
  
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_trace_rating_filter_useful" {
  value = length(data.huaweicloud_cts_traces.filter_by_trace_rating.traces) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_traces.filter_by_trace_rating.traces[*] : v.trace_rating == local.trace_rating]
  )
}
`, testDataSourceCtsTraces_base(name), from, to)
}

func testDataSourceCtsTraces_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%[1]s"
  acl           = "public-read"
  force_destroy = true
}
	  
resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name        = huaweicloud_obs_bucket.bucket.bucket
  file_prefix        = "cts"
  lts_enabled        = true
  compress_type      = "gzip"
  is_sort_by_service = false
}

resource "huaweicloud_obs_bucket" "data_bucket" {
  bucket = "%[1]sdata"
  acl    = "public-read"
}

resource "huaweicloud_cts_data_tracker" "tracker" {
  name        = "%[1]s-data"
  data_bucket = huaweicloud_obs_bucket.data_bucket.bucket
  lts_enabled = true
}
`, name)
}
