package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudDisPartionV2DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDisStreamV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisPartitionV2_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisPartitionV2Exists(),
				),
			},
		},
	})
}

func testAccCheckDisPartitionV2Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.huaweicloud_dis_partition_v2.partition"]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_dis_partition_v2.partition exist, err=not found this resource")
		}

		if _, ok := rs.Primary.Attributes["partitions.0.id"]; !ok {
			return fmtp.Errorf("expect partitions to be set")
		}

		return nil
	}
}

func testAccDisPartitionV2_basic(random string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dis_partition_v2" "partition" {
  stream_name = "${huaweicloud_dis_stream_v2.stream.stream_name}"
}
`, testAccDisStreamV2_basic(random))
}
