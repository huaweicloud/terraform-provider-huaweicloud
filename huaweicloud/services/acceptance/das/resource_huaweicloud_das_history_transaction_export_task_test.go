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

func getHistoryTransactionExportTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	return das.GetHistoryTransactionExportTask(client, instanceId, state.Primary.ID)
}

func TestAccHistoryTransactionExportTask_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_das_history_transaction_export_task.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getHistoryTransactionExportTaskResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccHistoryTransactionExportTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "bucket_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "export_line_num"),
					resource.TestCheckResourceAttr(rName, "time_zone", "UTC+8"),
					resource.TestCheckResourceAttr(rName, "order_field", "collectTime"),
					resource.TestCheckResourceAttr(rName, "order_by", "asc"),
					resource.TestCheckResourceAttr(rName, "last_sec_min", "0"),
					resource.TestCheckResourceAttr(rName, "last_sec_max", "100"),
					resource.TestMatchResourceAttr(rName, "created_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccHistoryTransactionExportTask_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "bucket_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "export_line_num"),
					resource.TestCheckResourceAttr(rName, "time_zone", "GMT+8"),
					resource.TestCheckResourceAttr(rName, "order_field", "occurrenceTime"),
					resource.TestCheckResourceAttr(rName, "order_by", "desc"),
					resource.TestCheckResourceAttr(rName, "last_sec_min", "1"),
					resource.TestCheckResourceAttr(rName, "last_sec_max", "50"),
					resource.TestMatchResourceAttr(rName, "created_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHistoryTransactionExportTaskImportIdFunc(rName),
				ImportStateVerifyIgnore: []string{
					"start_time",
					"end_time",
					"bucket_name",
					"file_path",
					"time_zone",
					"order_field",
					"order_by",
					"last_sec_min",
					"last_sec_max",
					"download_url",
					"enable_force_new",
				},
			},
		},
	})
}

func testAccHistoryTransactionExportTaskImportIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccHistoryTransactionExportTask_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

locals {
  instance_ids = split(",", "%[2]s")
}
`, name, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccHistoryTransactionExportTask_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_history_transaction_export_task" "test" {
  instance_id  = local.instance_ids[0]
  bucket_name  = huaweicloud_obs_bucket.test.bucket
  start_time   = "2000-06-01T00:00:00+08:00"
  end_time     = "2099-06-02T00:00:00+08:00"
  time_zone    = "UTC+8"
  order_field  = "collectTime"
  order_by     = "asc"
  last_sec_min = 0
  last_sec_max = 100
}
`, testAccHistoryTransactionExportTask_base(name))
}

func testAccHistoryTransactionExportTask_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_history_transaction_export_task" "test" {
  instance_id  = local.instance_ids[0]
  bucket_name  = huaweicloud_obs_bucket.test.bucket
  start_time   = "2000-06-03T00:00:00+08:00"
  end_time     = "2099-06-04T00:00:00+08:00"
  time_zone    = "GMT+8"
  order_field  = "occurrenceTime"
  order_by     = "desc"
  last_sec_min = 1
  last_sec_max = 50

  enable_force_new = "true"
}
`, testAccHistoryTransactionExportTask_base(name))
}
