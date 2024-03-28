package dli

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDliFlinkJarJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}
	jobId, _ := strconv.Atoi(state.Primary.ID)
	return flinkjob.Get(client, jobId)
}

func TestAccResourceDliFlinkJarJob_basic(t *testing.T) {
	var obj flinkjob.CreateJarJobOpts
	resourceName := "huaweicloud_dli_flinkjar_job.test"
	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliFlinkJarJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliJarPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJarJobResource_basic(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "job_running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccFlinkJarJobResource_update(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccFlinkJarJobResource_base(bucketName, name string) string {
	return fmt.Sprintf(`
variable "ak" {
  type        = string
  description = "value"
  default     = "%[1]s"
}

variable "sk" {
  type        = string
  description = "value"
  default     = "%[2]s"
}

variable "jarObsPath" {
  type        = string
  description = "value"
  default     = "%[3]s"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[4]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_dli_package" "test" {
  group_name  = "%[5]s"
  type        = "jar"
  object_path = var.jarObsPath
}

resource "huaweicloud_dli_queue" "test" {
  name       = "%[5]s"
  cu_count   = 16
  queue_type = "general"

  tags = {
    foo = "bar"
  }
}
`, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY, acceptance.HW_DLI_FLINK_JAR_OBS_PATH, bucketName, name)
}

func testAccFlinkJarJobResource_basic(name, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinkjar_job" "test" {
  name            = "%[2]s"
  queue_name      = "common_test"
  entrypoint      = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/${huaweicloud_obs_bucket.test.bucket}/output"
  obs_bucket      = huaweicloud_obs_bucket.test.bucket
  log_enabled     = true

  tags = {
    foo = "bar"
  }
}
`, testAccFlinkJarJobResource_base(bucketName, name), name, acceptance.HW_REGION_NAME)
}

func testAccFlinkJarJobResource_update(name, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinkjar_job" "test" {
  name            = "%[2]s"
  queue_name      = huaweicloud_dli_queue.test.name
  entrypoint      = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/${huaweicloud_obs_bucket.test.bucket}/output"
  obs_bucket      = huaweicloud_obs_bucket.test.bucket
  log_enabled     = true

  tags = {
    owner = "terraform"
  }
}
`, testAccFlinkJarJobResource_base(bucketName, name), name, acceptance.HW_REGION_NAME)
}

func TestAccResourceDliFlinkJarJob_all(t *testing.T) {
	var obj flinkjob.CreateJarJobOpts
	resourceName := "huaweicloud_dli_flinkjar_job.test"
	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliFlinkJarJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliJarPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJarJobResource_all(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "job_running"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", name),
					resource.TestCheckResourceAttrPair(resourceName, "obs_bucket", "huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "description of dli job"),
					resource.TestCheckResourceAttr(resourceName, "main_class", "com.huaweicloud.dli.obs.Main"),
					resource.TestCheckResourceAttr(resourceName, "feature", "basic"),
					resource.TestCheckResourceAttr(resourceName, "flink_version", "1.10"),
					resource.TestCheckResourceAttr(resourceName, "cu_num", "6"),
					resource.TestCheckResourceAttr(resourceName, "manager_cu_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "parallel_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "tm_slot_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "smn_topic", name),
					resource.TestCheckResourceAttr(resourceName, "restart_when_exception", "true"),
					resource.TestCheckResourceAttr(resourceName, "resume_max_num", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFlinkJarJobResource_all(name, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_dli_flinkjar_job" "test" {
  name                   = "%[2]s"
  queue_name             = huaweicloud_dli_queue.test.name
  entrypoint             = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args        = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/${huaweicloud_obs_bucket.test.bucket}/output"
  obs_bucket             = huaweicloud_obs_bucket.test.bucket
  log_enabled            = true
  description            = "description of dli job"
  main_class             = "com.huaweicloud.dli.obs.Main"
  dependency_files       = ["jar_package/user.csv"]
  feature                = "basic"
  flink_version          = "1.10"
  cu_num                 = 6
  manager_cu_num         = 2
  parallel_num           = 2
  tm_cu_num              = 2
  tm_slot_num            = 1
  smn_topic              = huaweicloud_smn_topic.test.name
  restart_when_exception = "true"
  resume_checkpoint      = true
  resume_max_num         = 2
  checkpoint_path        = "${huaweicloud_obs_bucket.test.bucket}/checkpoint/"
}
`, testAccFlinkJarJobResource_base(bucketName, name), name, acceptance.HW_REGION_NAME)
}
