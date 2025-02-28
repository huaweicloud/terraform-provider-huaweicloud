package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCheckBucket_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cts_check_bucket.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)
	baseConfig := testDataSourceCheckBucket_base(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCtsKmsId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCheckBucket_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.bucket_name", rName),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.check_bucket_response.0.response_code", "200"),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.check_bucket_response.0.success", "true"),
				),
			},
			{
				Config: testDataSourceCheckBucket_KmsId(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.bucket_name", rName),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.check_bucket_response.0.response_code", "200"),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.check_bucket_response.0.success", "true"),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.is_support_trace_files_encryption", "true"),
					resource.TestCheckResourceAttr(dataSource, "buckets.0.kms_id", acceptance.HW_CTS_KMS_ID),
				),
			},
		},
	})
}

func testDataSourceCheckBucket_basic(config string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cts_check_bucket" "test" {
  bucket_name     = huaweicloud_obs_bucket.bucket.bucket
  bucket_location = huaweicloud_obs_bucket.bucket.region
}
`, config)
}

func testDataSourceCheckBucket_KmsId(config string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cts_check_bucket" "test" {
  bucket_name                       = huaweicloud_obs_bucket.bucket.bucket
  bucket_location                   = huaweicloud_obs_bucket.bucket.region
  kms_id                            = "%[2]s"
  is_support_trace_files_encryption = true
}
`, config, acceptance.HW_CTS_KMS_ID)
}

func testDataSourceCheckBucket_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}
`, rName)
}
