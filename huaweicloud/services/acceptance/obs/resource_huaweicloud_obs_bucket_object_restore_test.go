package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBucketObjectRestore_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectRestore_basic(name),
			},
		},
	})
}

func testAccBucketObjectRestore_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  storage_class = "COLD"
  force_destroy = true
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket        = huaweicloud_obs_bucket.test.bucket
  key           = "%[1]s"
  content       = "some archived content"
  storage_class = "COLD"
}

resource "huaweicloud_obs_bucket_object" "test_with_expedited" {
  bucket        = huaweicloud_obs_bucket.test.bucket
  key           = "%[1]s-expedited"
  content       = "some archived content"
  storage_class = "COLD"
}
`, name)
}

func testAccBucketObjectRestore_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object_restore" "test" {
  bucket = huaweicloud_obs_bucket.test.bucket
  key    = huaweicloud_obs_bucket_object.test.key
  days   = 1
  tier   = "standard"
}

resource "huaweicloud_obs_bucket_object_restore" "test_with_expedited" {
  bucket = huaweicloud_obs_bucket.test.bucket
  key    = huaweicloud_obs_bucket_object.test_with_expedited.key
  days   = 7
  tier   = "expedited"
}
`, testAccBucketObjectRestore_base(name))
}
