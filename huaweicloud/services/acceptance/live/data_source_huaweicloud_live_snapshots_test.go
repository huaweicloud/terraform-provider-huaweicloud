package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSnapshots_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_live_snapshots.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byAppName   = "data.huaweicloud_live_snapshots.filter_by_app_name"
		dcByAppName = acceptance.InitDataSourceCheck(byAppName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveIngestDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSnapshots_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.call_back_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.frequency"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.storage_bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.storage_location"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.storage_mode"),

					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSnapshots_base() string {
	bucketName := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_live_bucket_authorization" "test" {
  depends_on = [huaweicloud_obs_bucket.test]

  bucket = "%[1]s"
}

resource "huaweicloud_live_snapshot" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test]

  domain_name       = "%[2]s"
  app_name          = "live"
  frequency         = 10
  storage_mode      = 0
  storage_bucket    = huaweicloud_obs_bucket.test.bucket
  storage_path      = "input/"
  call_back_enabled = "off"
}
`, bucketName, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testDataSourceSnapshots_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_live_snapshots" "test" {
  depends_on = [huaweicloud_live_snapshot.test]

  domain_name = "%[2]s"
}

# Filter by app_name
locals {
  app_name = data.huaweicloud_live_snapshots.test.snapshots[0].app_name
}

data "huaweicloud_live_snapshots" "filter_by_app_name" {
  depends_on = [huaweicloud_live_snapshot.test]

  domain_name = "%[2]s"
  app_name    = local.app_name
}

locals {
  app_name_filter_result = [
    for v in data.huaweicloud_live_snapshots.filter_by_app_name.snapshots[*].app_name : v == local.app_name
  ]
}

output "app_name_filter_is_useful" {
  value = alltrue(local.app_name_filter_result) && length(local.app_name_filter_result) > 0
}
`, testDataSourceSnapshots_base(), acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}
