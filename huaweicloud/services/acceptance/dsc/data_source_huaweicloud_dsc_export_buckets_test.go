package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscExportBuckets_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dsc_export_buckets.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscExportBuckets_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.0.bucket_location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.0.bucket_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.0.bind_task"),
					resource.TestCheckResourceAttrSet(dataSourceName, "buckets.0.enable_audit_status"),
				),
			},
		},
	})
}

const testAccDataSourceDscExportBuckets_basic = `
data "huaweicloud_dsc_export_buckets" "test" {
}
`
