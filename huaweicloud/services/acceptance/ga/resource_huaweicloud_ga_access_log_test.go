package ga

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ga"
)

func getAccessLogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("ga", region)
	if err != nil {
		return nil, fmt.Errorf("error creating GA client: %s", err)
	}

	return ga.GetAccessLogInfo(client, state.Primary.ID)
}

func TestAccAccessLog_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ga_access_log.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccessLogResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccessLog_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resource_type", "LISTENER"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_ga_listener.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test1", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccessLog_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test2", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccessLog_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test1" {
  group_name  = "%[1]s-group-1"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test1" {
  group_id    = huaweicloud_lts_group.test1.id
  stream_name = "%[1]s-stream-1"
}

resource "huaweicloud_lts_group" "test2" {
  group_name  = "%[1]s-group-2"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test2" {
  group_id    = huaweicloud_lts_group.test2.id
  stream_name = "%[1]s-stream-2"
}
`, name)
}

func testAccessLog_basic(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_ga_access_log" "test" {
  resource_type = "LISTENER"
  resource_id   = huaweicloud_ga_listener.test.id
  log_group_id  = huaweicloud_lts_group.test1.id
  log_stream_id = huaweicloud_lts_stream.test1.id
}
`, testAccessLog_base(name), testListener_basic(name))
}

func testAccessLog_update(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_ga_access_log" "test" {
  resource_type = "LISTENER"
  resource_id   = huaweicloud_ga_listener.test.id
  log_group_id  = huaweicloud_lts_group.test2.id
  log_stream_id = huaweicloud_lts_stream.test2.id
}
`, testAccessLog_base(name), testListener_basic(name))
}
