package oms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getMigrationTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getTaskHttpUrl = "v2/{project_id}/tasks/{task_id}"
		getTaskProduct = "oms"
	)
	getTaskClient, err := cfg.NewServiceClient(getTaskProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OMS client: %s", err)
	}

	getTaskPath := getTaskClient.Endpoint + getTaskHttpUrl
	getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", getTaskClient.ProjectID)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", state.Primary.ID)

	getTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTaskResp, err := getTaskClient.Request("GET", getTaskPath, &getTaskOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving OMS migration task: %s", err)
	}

	return utils.FlattenResponse(getTaskResp)
}

func TestAccMigrationTask_object(t *testing.T) {
	var task interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	sourceBucketName := rName + "-source"
	destBucketName := rName + "-dest"
	smnName := rName + "-smn"
	resourceName := "huaweicloud_oms_migration_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&task,
		getMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationTask_object(testMigrationTask_base(sourceBucketName, destBucketName), smnName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "object"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task"),
					resource.TestCheckResourceAttr(resourceName, "start_task", "false"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.max_bandwidth", "1"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.max_bandwidth", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "smn_config.0.topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
				),
			},
			{
				Config: testMigrationTask_update(testMigrationTask_base(sourceBucketName, destBucketName), smnName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "object"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task"),
					resource.TestCheckResourceAttr(resourceName, "start_task", "true"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.max_bandwidth", "2"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.max_bandwidth", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "smn_config.0.topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
				),
			},
		},
	})
}

func TestAccMigrationTask_prefix(t *testing.T) {
	var task interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	sourceBucketName := rName + "-source"
	destBucketName := rName + "-dest"
	resourceName := "huaweicloud_oms_migration_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&task,
		getMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationTask_prefix(testMigrationTask_base(sourceBucketName, destBucketName)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "prefix"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task"),
					resource.TestCheckResourceAttr(resourceName, "start_task", "true"),
				),
			},
		},
	})
}

func TestAccMigrationTask_list(t *testing.T) {
	var task interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	sourceBucketName := rName + "-source"
	destBucketName := rName + "-dest"
	listFileBucketName := rName + "-list"
	resourceName := "huaweicloud_oms_migration_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&task,
		getMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationTask_list(testMigrationTask_base(sourceBucketName, destBucketName), listFileBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "list"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task"),
					resource.TestCheckResourceAttr(resourceName, "start_task", "true"),
				),
			},
		},
	})
}

func testMigrationTask_base(sourceBucketName, destBucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "source" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket_object" "source" {
  bucket  = huaweicloud_obs_bucket.source.bucket
  key     = "test.txt"
  content = "test content"
}

resource "huaweicloud_obs_bucket" "dest" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}
`, sourceBucketName, destBucketName)
}

func testMigrationTask_object(baseConfig, smnName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "The display name of %[2]s"
}

resource "huaweicloud_oms_migration_task" "test" {
  source_object {
    data_source = "HuaweiCloud"
    region      = "%[3]s"
    bucket      = huaweicloud_obs_bucket.source.bucket
    access_key  = "%[4]s"
    secret_key  = "%[5]s"
    object      = [""]
  }

  destination_object {
    region     = "%[3]s"
    bucket     = huaweicloud_obs_bucket.dest.bucket
    access_key = "%[4]s"
    secret_key = "%[5]s"
  }

  start_task  = false
  type        = "object"
  description = "test task"

  enable_metadata_migration = true

  bandwidth_policy {
    max_bandwidth = 1
    start         = "15:00"
    end           = "16:00"
  }

  bandwidth_policy {
    max_bandwidth = 2
    start         = "16:00"
    end           = "17:00"
  }

  smn_config {
    topic_urn          = huaweicloud_smn_topic.test.topic_urn
    trigger_conditions = ["FAILURE", "SUCCESS"]
  }
}
`, baseConfig, smnName, acceptance.HW_REGION_NAME, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testMigrationTask_update(baseConfig, smnName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "The display name of %[2]s"
}

resource "huaweicloud_oms_migration_task" "test" {
  source_object {
    data_source = "HuaweiCloud"
    region      = "%[3]s"
    bucket      = huaweicloud_obs_bucket.source.bucket
    access_key  = "%[4]s"
    secret_key  = "%[5]s"
    object      = [""]
  }

  destination_object {
    region     = "%[3]s"
    bucket     = huaweicloud_obs_bucket.dest.bucket
    access_key = "%[4]s"
    secret_key = "%[5]s"
  }

  start_task  = true
  type        = "object"
  description = "test task"

  bandwidth_policy {
    max_bandwidth = 2
    start         = "15:00"
    end           = "16:00"
  }

  bandwidth_policy {
    max_bandwidth = 3
    start         = "16:00"
    end           = "17:00"
  }

  smn_config {
    topic_urn          = huaweicloud_smn_topic.test.topic_urn
    trigger_conditions = ["FAILURE", "SUCCESS"]
  }
}
`, baseConfig, smnName, acceptance.HW_REGION_NAME, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testMigrationTask_prefix(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_oms_migration_task" "test" {
  source_object {
    data_source = "HuaweiCloud"
    region      = "%[2]s"
    bucket      = huaweicloud_obs_bucket.source.bucket
    access_key  = "%[3]s"
    secret_key  = "%[4]s"
    object      = ["test"]
  }

  destination_object {
    region     = "%[2]s"
    bucket     = huaweicloud_obs_bucket.dest.bucket
    access_key = "%[3]s"
    secret_key = "%[4]s"
  }

  type        = "prefix"
  description = "test task"
}
`, baseConfig, acceptance.HW_REGION_NAME, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testMigrationTask_list(baseConfig, listFileBucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "list_file_bucket" {
  bucket = "%[2]s"
  acl    = "private"
}

resource "huaweicloud_obs_bucket_object" "list_file_object" {
  bucket        = huaweicloud_obs_bucket.list_file_bucket.bucket
  key           = "list_file.txt"
  content       = "test.txt"
}

resource "huaweicloud_oms_migration_task" "test" {
  source_object {
    data_source      = "HuaweiCloud"
    region           = "%[3]s"
    bucket           = huaweicloud_obs_bucket.source.bucket
    access_key       = "%[4]s"
    secret_key       = "%[5]s"
    list_file_bucket = huaweicloud_obs_bucket.list_file_bucket.bucket
    list_file_key    = huaweicloud_obs_bucket_object.list_file_object.key
  }

  destination_object {
    region     = "%[3]s"
    bucket     = huaweicloud_obs_bucket.dest.bucket
    access_key = "%[4]s"
    secret_key = "%[5]s"
  }

  type        = "list"  
  description = "test task"
}
`, baseConfig, listFileBucketName, acceptance.HW_REGION_NAME, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}
