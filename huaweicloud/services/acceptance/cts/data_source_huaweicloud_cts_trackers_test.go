package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCtsTrackers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cts_trackers.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCtsTrackers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "trackers.0.status"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCtsTrackers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id     = huaweicloud_cts_data_tracker.tracker.id
  name   = huaweicloud_cts_tracker.tracker.name
  type   = "data"
  status = huaweicloud_cts_data_tracker.tracker.status
}

data "huaweicloud_cts_trackers" "test" {
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cts_trackers.test.trackers) >= 1
}

data "huaweicloud_cts_trackers" "filter_by_id" {
  tracker_id = local.id

  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cts_trackers.filter_by_id.trackers) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_trackers.filter_by_id.trackers[*] : v.id == local.id]
  )
}

data "huaweicloud_cts_trackers" "filter_by_name" {
  name = local.name
  
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_cts_trackers.filter_by_name.trackers) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_trackers.filter_by_name.trackers[*] : v.name == local.name]
  )
}

data "huaweicloud_cts_trackers" "filter_by_type" {
  type = local.type
  
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_cts_trackers.filter_by_type.trackers) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_trackers.filter_by_type.trackers[*] : v.type == local.type]
  )
}

data "huaweicloud_cts_trackers" "filter_by_status" {
  status = local.status
  
  depends_on = [
    huaweicloud_cts_tracker.tracker,
    huaweicloud_cts_data_tracker.tracker,
  ]
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_cts_trackers.filter_by_status.trackers) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_trackers.filter_by_status.trackers[*] : v.status == local.status]
  )
}
`, testDataSourceCtsTrackers_base(name))
}

func testDataSourceCtsTrackers_base(name string) string {
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
