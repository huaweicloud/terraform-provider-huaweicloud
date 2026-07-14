package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceTransfers_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_lts_transfers.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byLogGroupName   = "data.huaweicloud_lts_transfers.filter_by_log_group_name"
		dcByLogGroupName = acceptance.InitDataSourceCheck(byLogGroupName)

		byLogStreamName   = "data.huaweicloud_lts_transfers.filter_by_log_stream_name"
		dcByLogStreamName = acceptance.InitDataSourceCheck(byLogStreamName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceTransfers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "transfers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByLogGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_group_name_filter_useful", "true"),
					resource.TestMatchResourceAttr(byLogGroupName, "transfers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrPair(byLogGroupName, "transfers.0.log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "transfers.0.log_group_name", "huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_streams.#", "1"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "transfers.0.log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "transfers.0.log_streams.0.log_stream_name",
						"huaweicloud_lts_stream.test", "stream_name"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.#", "1"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.#", "1"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_period", "3"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "hour"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_bucket_name",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name",
						"lts_transfer_obs_"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC"),
					resource.TestCheckResourceAttr(byLogGroupName, "transfers.0.log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id",
						"Etc/GMT"),
					dcByLogStreamName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_stream_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceTransfers_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = replace("%[1]s", "_", "-")
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = huaweicloud_obs_bucket.test.bucket
      obs_dir_prefix_name = "lts_transfer_obs_"
      obs_prefix_name     = "obs_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }
  }
}`, name)
}

func testAccDatasourceTransfers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_transfers" "test" {
  depends_on = [huaweicloud_lts_transfer.test]
}

# Filter by log group name
locals {
  log_group_name = huaweicloud_lts_transfer.test.log_group_name
}

data "huaweicloud_lts_transfers" "filter_by_log_group_name" {
  log_group_name = local.log_group_name
}

locals {
  log_group_name_filter_result = [
    for v in data.huaweicloud_lts_transfers.filter_by_log_group_name.transfers[*].log_group_name : v == local.log_group_name
  ]
}

output "is_log_group_name_filter_useful" {
  value = length(local.log_group_name_filter_result) > 0 && alltrue(local.log_group_name_filter_result)
}

# Filter by log stream name
locals {
  log_stream_name = huaweicloud_lts_transfer.test.log_streams[0].log_stream_name
}

data "huaweicloud_lts_transfers" "filter_by_log_stream_name" {
  log_stream_name = local.log_stream_name
}

locals {
  log_stream_name_filter_result = [
    for v in data.huaweicloud_lts_transfers.filter_by_log_stream_name.transfers[*].log_streams[0].log_stream_name : v == local.log_stream_name
  ]
}

output "is_log_stream_name_filter_useful" {
  value = length(local.log_stream_name_filter_result) > 0 && alltrue(local.log_stream_name_filter_result)
}
`, testAccDatasourceTransfers_base(name))
}
