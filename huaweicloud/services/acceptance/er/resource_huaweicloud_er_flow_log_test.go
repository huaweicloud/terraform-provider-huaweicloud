package er

import (
	"fmt"
	"regexp"
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
		flowLog interface{}
		rName   = "huaweicloud_er_flow_log.test"
		rc      = acceptance.InitResourceCheck(rName, &flowLog, getFlowResourceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testaccFlowLog_base(name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFlowLog_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_store_type", "LTS"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "resource_type", "attachment"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testFlowLog_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
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

func testAccFlowLogImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, flowLogId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER flow log is not found in the tfstate", rsName)
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

func testaccFlowLog_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)
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

func testFlowLog_basic_step1(baseConfig, name string) string {
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
  description    = "Created by script"
  enabled        = false
}
`, baseConfig, name)
}

func testFlowLog_basic_step2(baseConfig, name string) string {
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
  enabled        = true
}
`, baseConfig, name)
}
