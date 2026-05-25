package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/das"
)

func getBinlogParseTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}

	connectionId := state.Primary.Attributes["user_id"]
	return das.GetBinlogParseTask(client, connectionId, state.Primary.ID)
}

func TestAccBinlogParseTask_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_das_binlog_parse_task.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getBinlogParseTaskResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBinlogParseTask_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "binlog_type", "latest"),
					resource.TestCheckResourceAttrSet(rName, "file_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccBinlogParseTaskImportIdFunc(rName),
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccBinlogParseTaskImportIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["user_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["user_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["user_id"], rs.Primary.ID), nil
	}
}

func testAccBinlogParseTask_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}

data "huaweicloud_das_database_users" "test" {
  instance_id = local.instance_ids[0]
}

locals {
  user_id = try(data.huaweicloud_das_database_users.test.users.0.id, "")
}

data "huaweicloud_das_binlogs" "all" {
  user_id     = local.user_id
  binlog_type = "latest"
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccBinlogParseTask_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_binlog_parse_task" "test" {
  user_id     = local.user_id
  binlog_type = "latest"
  file_name   = try(data.huaweicloud_das_binlogs.all.binlogs.0.file_name, "")
}
`, testAccBinlogParseTask_base())
}
