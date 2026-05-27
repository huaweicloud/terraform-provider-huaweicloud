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

func getBinlogParseTaskExportResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}

	userId := state.Primary.Attributes["user_id"]
	return das.GetBinlogParseTaskExport(client, userId, state.Primary.ID)
}

func TestAccBinlogParseTaskExport_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_das_binlog_parse_task_export.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getBinlogParseTaskExportResourceFunc)

		name = acceptance.RandomAccResourceName()
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
				Config: testAccBinlogParseTaskExport_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "task_id"),
					resource.TestCheckResourceAttrSet(rName, "bucket_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "export_line_num"),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBinlogParseTaskExportImportIdFunc(rName),
				ImportStateVerifyIgnore: []string{
					"task_id",
					"download_url",
					"status",
					"filter_condition",
					"filter_condition.0",
					"filter_condition.0.types",
					"filter_condition.0.start_time",
					"filter_condition.0.end_time",
					"filter_condition.0.parse_double_insert",
				},
			},
		},
	})
}

func testAccBinlogParseTaskExportImportIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["user_id"] == "" || rs.Primary.Attributes["bucket_name"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", rs.Primary.Attributes["user_id"],
				rs.Primary.Attributes["bucket_name"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["user_id"],
			rs.Primary.Attributes["bucket_name"], rs.Primary.ID), nil
	}
}

func testAccBinlogParseTaskExport_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "tf-test-bucket"
  acl           = "private"
  force_destroy = true
}

locals {
  instance_ids = split(",", "%[2]s")
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

resource "huaweicloud_das_binlog_parse_task" "test" {
  user_id     = local.user_id
  binlog_type = "latest"
  file_name   = try(data.huaweicloud_das_binlogs.all.binlogs.0.file_name, "")
}
`, name, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccBinlogParseTaskExport_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_binlog_parse_task_export" "test" {
  user_id     = local.user_id
  task_id     = huaweicloud_das_binlog_parse_task.test.id
  bucket_name = huaweicloud_obs_bucket.test.bucket

  filter_condition {
    types               = ["insert", "update", "delete", "ddl"]
    start_time          = "2000-06-01T00:00:00+08:00"
    end_time            = "2099-06-02T00:00:00+08:00"
    parse_double_insert = true
  }
}
`, testAccBinlogParseTaskExport_base(name))
}
