package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getLtsStreamResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.LtsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	resourceID := state.Primary.ID
	groupID := state.Primary.Attributes["group_id"]
	streams, err := logstreams.List(client, groupID).Extract()
	if err != nil {
		return nil, err
	}

	for _, item := range streams.LogStreams {
		if item.ID == resourceID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("the log stream %s does not exist", resourceID)
}

func TestAccLtsStream_basic(t *testing.T) {
	var stream logstreams.LogStream
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_lts_stream.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&stream,
		getLtsStreamResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsStream_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter_count", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStreamImportState(resourceName),
			},
		},
	})
}

func testLtsStreamImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		streamID := rs.Primary.ID
		groupID := rs.Primary.Attributes["group_id"]

		return fmt.Sprintf("%s/%s", groupID, streamID), nil
	}
}

func testAccLtsStream_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}
`, rName)
}
