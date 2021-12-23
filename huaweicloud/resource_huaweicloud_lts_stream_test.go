package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccLogTankStreamV2_basic(t *testing.T) {
	var stream logstreams.LogStream
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogTankStreamV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankStreamV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankStreamV2Exists(
						"huaweicloud_lts_stream.testacc_stream", &stream),
					resource.TestCheckResourceAttr("huaweicloud_lts_stream.testacc_stream", "stream_name", rName),
					resource.TestCheckResourceAttr("huaweicloud_lts_stream.testacc_stream", "filter_count", "0"),
				),
			},
		},
	})
}

func testAccCheckLogTankStreamV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	ltsclient, err := config.LtsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lts_stream" {
			continue
		}

		group_id := rs.Primary.Attributes["group_id"]
		_, err = logstreams.List(ltsclient, group_id).Extract()
		if err == nil {
			return fmtp.Errorf("Log group (%s) still exists.", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckLogTankStreamV2Exists(n string, stream *logstreams.LogStream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		ltsclient, err := config.LtsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
		}

		group_id := rs.Primary.Attributes["group_id"]
		streams, err := logstreams.List(ltsclient, group_id).Extract()
		if err != nil {
			return fmtp.Errorf("Log stream get list err: %s", err.Error())
		}
		for _, logstream := range streams.LogStreams {
			if logstream.ID == rs.Primary.ID {
				*stream = logstream
				return nil
			}
		}

		return fmtp.Errorf("Error HuaweiCloud log stream %s: No Found", rs.Primary.ID)
	}
}

func testAccLogTankStreamV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "testacc_group" {
  group_name  = "%s"
  ttl_in_days = 1
}
resource "huaweicloud_lts_stream" "testacc_stream" {
  group_id    = huaweicloud_lts_group.testacc_group.id
  stream_name = "%s"
}
`, rName, rName)
}
