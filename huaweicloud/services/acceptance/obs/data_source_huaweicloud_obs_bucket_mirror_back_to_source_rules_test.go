package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccObsBucketMirrorBackToSourceRules_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_obs_bucket_mirror_back_to_source_rules.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBSBucketName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketMirrorBackToSourceRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "rules"),
				),
			},
		},
	})
}

func testAccObsBucketMirrorBackToSourceRules_basic() string {
	return fmt.Sprintf(`
# The bucket must contain the already created back to source rules.
data "huaweicloud_obs_bucket_mirror_back_to_source_rules" "test" {
  bucket = "%[1]s"
}
`, acceptance.HW_OBS_BUCKET_NAME)
}
