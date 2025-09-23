package oms

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getMigrationSyncTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	queryTime := strconv.FormatInt(time.Now().UnixMilli(), 10)

	var (
		getSyncTaskHttpUrl = "v2/{project_id}/sync-tasks/{sync_task_id}?query_time=" + queryTime
		getSyncTaskProduct = "oms"
	)
	getSyncTaskClient, err := cfg.NewServiceClient(getSyncTaskProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OMS client: %s", err)
	}

	getSyncTaskPath := getSyncTaskClient.Endpoint + getSyncTaskHttpUrl
	getSyncTaskPath = strings.ReplaceAll(getSyncTaskPath, "{project_id}", getSyncTaskClient.ProjectID)
	getSyncTaskPath = strings.ReplaceAll(getSyncTaskPath, "{sync_task_id}", state.Primary.ID)

	getSyncTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getSyncTaskResp, err := getSyncTaskClient.Request("GET", getSyncTaskPath, &getSyncTaskOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving OMS migration sync task: %s", err)
	}

	return utils.FlattenResponse(getSyncTaskResp)
}

func TestAccMigrationSyncTask_basic(t *testing.T) {
	var syncTask interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	sourceBucketName := rName + "-source"
	destBucketName := rName + "-dest"
	resourceName := "huaweicloud_oms_migration_sync_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&syncTask,
		getMigrationSyncTaskResourceFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationSyncTask_basic(testMigrationTask_base(sourceBucketName, destBucketName), "stop"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "src_cloud_type", "HuaweiCloud"),
					resource.TestCheckResourceAttr(resourceName, "enable_kms", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "STOPPED"),
					resource.TestCheckResourceAttr(resourceName, "consistency_check", "crc64"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "last_start_at"),
				),
			},
			{
				Config: testMigrationSyncTask_basic(testMigrationTask_base(sourceBucketName, destBucketName), "start"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "src_cloud_type", "HuaweiCloud"),
					resource.TestCheckResourceAttr(resourceName, "enable_kms", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "SYNCHRONIZING"),
					resource.TestCheckResourceAttr(resourceName, "consistency_check", "crc64"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "last_start_at"),
				),
			},
		},
	})
}

func testMigrationSyncTask_basic(baseConfig, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_oms_migration_sync_task" "test" {
  src_cloud_type            = "HuaweiCloud"
  src_ak                    = "%[3]s"
  src_sk                    = "%[4]s"
  src_region                = "%[2]s"
  src_bucket                = huaweicloud_obs_bucket.source.bucket
  dst_ak                    = "%[3]s"
  dst_sk                    = "%[4]s"
  dst_bucket                = huaweicloud_obs_bucket.dest.bucket
  enable_kms                = true
  enable_metadata_migration = true
  enable_restore            = true
  consistency_check         = "crc64" 
  action                    = "%[5]s"
}
`, baseConfig, acceptance.HW_REGION_NAME, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY, action)
}
