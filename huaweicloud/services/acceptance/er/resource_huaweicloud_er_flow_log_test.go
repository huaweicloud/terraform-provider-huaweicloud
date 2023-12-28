package er

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getFlowResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	getFlowLogHttpUrl := "enterprise-router/{er_id}/flow-logs/{flow_log_id}"
	getFlowLogPath := client.ResourceBaseURL() + getFlowLogHttpUrl
	getFlowLogPath = strings.ReplaceAll(getFlowLogPath, "{er_id}", state.Primary.Attributes["instance_id"])
	getFlowLogPath = strings.ReplaceAll(getFlowLogPath, "{flow_log_id}", state.Primary.ID)

	getFlowLogOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getflowLogResp, err := client.Request("GET", getFlowLogPath, &getFlowLogOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving flow log: %s", err)
	}

	return utils.FlattenResponse(getflowLogResp)
}

func TestAccFlowLog_basic(t *testing.T) {
	var (
		flowLog  interface{}
		name     = acceptance.RandomAccResourceName()
		rName    = "huaweicloud_er_flow_log.test"
		bgpAsNum = acctest.RandIntRange(64512, 65534)
		rc       = acceptance.InitResourceCheck(
			rName,
			&flowLog,
			getFlowResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFlowLog_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_store_type", "LTS"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "resource_type", "attachment"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create ER flow log"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testFlowLog_basic_update(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFlowLogImportStateFunc(rName),
			},
		},
	})
}

func testAccFlowLogImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, flowLogId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		flowLogId = rs.Primary.ID
		if instanceId == "" || flowLogId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, flowLogId)
		}
		return fmt.Sprintf("%s/%s", instanceId, flowLogId), nil
	}
}

func testaccFlowLog_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name               = "%[2]s"
  asn                = %[3]d
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 7
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id            = huaweicloud_er_instance.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  name                   = "%[2]s"
  auto_create_vpc_routes = true

  tags = {
    foo = "bar"
  }
}
`, common.TestVpc(name), name, bgpAsNum)
}

func testFlowLog_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_flow_log" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  log_store_type = "LTS"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  resource_type  = "attachment"
  resource_id    = huaweicloud_er_vpc_attachment.test.id
  name           = "%[2]s"
  description    = "Create ER flow log"
  enabled        = false
}
`, testaccFlowLog_base(name, bgpAsNum), name)
}

func testFlowLog_basic_update(name string, bgpAsNum int) string {
	return fmt.Sprintf(`

%[1]s

resource "huaweicloud_er_flow_log" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  log_store_type = "LTS"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  resource_type  = "attachment"
  resource_id    = huaweicloud_er_vpc_attachment.test.id
  name           = "%[2]s-update"
  description    = ""
  enabled        = true
}
`, testaccFlowLog_base(name, bgpAsNum), name)
}
