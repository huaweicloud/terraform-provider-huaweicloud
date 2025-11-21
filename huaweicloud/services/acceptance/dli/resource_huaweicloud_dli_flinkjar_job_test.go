package dli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func parseFlinkJarJobAgencyNames(agencies string) (agencyName, agencyUpdateName string) {
	if agencies == "" {
		return "", ""
	}

	agencyNames := strings.Split(agencies, ",")
	if len(agencyNames) < 2 {
		return "", ""
	}

	agencyName = agencyNames[0]
	agencyUpdateName = agencyNames[1]
	return
}

func TestAccFlinkJarJob_basic(t *testing.T) {
	var (
		name                         = acceptance.RandomAccResourceName()
		agencyName, agencyUpdateName = parseFlinkJarJobAgencyNames(acceptance.HW_DLI_FLINK_JAR_AGENCY_NAMES)

		obj          flinkjob.CreateJarJobOpts
		resourceName = "huaweicloud_dli_flinkjar_job.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDliFlinkJarJobResourceFunc)

		nameWithResourceConfigVersion = "huaweicloud_dli_flinkjar_job.test_resource_config_v2"
		rcWithResourceConfigVersion   = acceptance.InitResourceCheck(nameWithResourceConfigVersion, &obj,
			getDliFlinkJarJobResourceFunc)

		nameDefault = "huaweicloud_dli_flinkjar_job.default"
		rcDefault   = acceptance.InitResourceCheck(nameDefault, &obj, getDliFlinkJarJobResourceFunc)
	)

	// Currently, the authorization bucket API will return an error if the bucket is not authorized.
	// Before running the test, please provide the `HW_DLI_FLINK_JAR_OBS_BUCKET_NAME` has been authorized.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliJarPath(t)
			acceptance.TestAccPreCheckDliGenaralQueueName(t)
			acceptance.TestAccPreCheckDliFlinkVersion(t)
			acceptance.TestAccPreCheckDliFlinkJarObsBucketName(t)
			acceptance.TestAccPreCheckDliFlinkJarAgencyNames(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithResourceConfigVersion.CheckResourceDestroy(),
			rcDefault.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJarJob_basic_step1(name, agencyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", acceptance.HW_DLI_GENERAL_QUEUE_NAME),
					resource.TestCheckResourceAttr(resourceName, "obs_bucket", acceptance.HW_DLI_FLINK_JAR_OBS_BUCKET_NAME),
					resource.TestCheckResourceAttr(resourceName, "log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "description of dli job"),
					resource.TestCheckResourceAttr(resourceName, "feature", "basic"),
					resource.TestCheckResourceAttr(resourceName, "flink_version", acceptance.HW_DLI_FLINK_VERSION),
					resource.TestCheckResourceAttrSet(resourceName, "cu_num"),
					resource.TestCheckResourceAttr(resourceName, "manager_cu_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "parallel_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "tm_slot_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "smn_topic", name),
					resource.TestCheckResourceAttr(resourceName, "restart_when_exception", "true"),
					resource.TestCheckResourceAttr(resourceName, "resume_max_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "resume_checkpoint", "true"),
					resource.TestMatchResourceAttr(resourceName, "checkpoint_path",
						regexp.MustCompile(fmt.Sprintf("^%s/", acceptance.HW_DLI_FLINK_JAR_OBS_BUCKET_NAME))),
					resource.TestCheckResourceAttr(resourceName, "runtime_config.max_records_per_file", "10"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					rcWithResourceConfigVersion.CheckResourceExists(),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "checkpoint_enabled", "true"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "checkpoint_mode", "2"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "checkpoint_interval", "60"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "execution_agency_urn", agencyName),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config_version", "v2"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.max_slot", "4"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.parallel_number", "3"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.job_manager_resource_spec.0.cpu", "1"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.job_manager_resource_spec.0.memory", "8G"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.task_manager_resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.task_manager_resource_spec.0.memory", "6G"),
					resource.TestCheckResourceAttrSet(nameWithResourceConfigVersion, "checkpoint_path"),
				),
			},
			{
				Config: testAccFlinkJarJobResource_basic_step2(name, agencyUpdateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Update description of dli job."),
					resource.TestCheckResourceAttr(resourceName, "resume_checkpoint", "false"),
					resource.TestCheckResourceAttr(resourceName, "resume_max_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_config.max_records_per_file", "5"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					rcWithResourceConfigVersion.CheckResourceExists(),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "execution_agency_urn", agencyUpdateName),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.max_slot", "3"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.parallel_number", "2"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.job_manager_resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.job_manager_resource_spec.0.memory", "6G"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.task_manager_resource_spec.0.cpu", "1"),
					resource.TestCheckResourceAttr(nameWithResourceConfigVersion, "resource_config.0.task_manager_resource_spec.0.memory", "2G"),
					rcDefault.CheckResourceExists(),
					// Assert the fields that are Computed behavior.
					resource.TestCheckResourceAttr(nameDefault, "checkpoint_enabled", "false"),
					resource.TestCheckResourceAttr(nameDefault, "checkpoint_interval", "30"),
					resource.TestCheckResourceAttr(nameDefault, "checkpoint_mode", "1"),
					resource.TestCheckResourceAttr(nameDefault, "manager_cu_num", "1"),
					resource.TestCheckResourceAttr(nameDefault, "parallel_num", "1"),
					resource.TestCheckResourceAttr(nameDefault, "tm_cu_num", "1"),
					resource.TestCheckResourceAttr(nameDefault, "resume_max_num", "-1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config_version", "v1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.#", "1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.0.job_manager_resource_spec.#", "1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.0.job_manager_resource_spec.0.cpu", "1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.0.job_manager_resource_spec.0.memory", "4G"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.0.task_manager_resource_spec.#", "1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.0.task_manager_resource_spec.0.cpu", "1"),
					resource.TestCheckResourceAttr(nameDefault, "resource_config.0.task_manager_resource_spec.0.memory", "4G"),
				),
			},
			// The `status` value maybe change from `job_running` to `job_finish` after the job is created,
			// so we need to ignore it.
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
			{
				ResourceName:            nameWithResourceConfigVersion,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
			{
				ResourceName:            nameDefault,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func testAccFlinkJarJobResource_base(name string) string {
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

resource "huaweicloud_dli_package" "test" {
  group_name  = "%[4]s"
  type        = "jar"
  object_path = var.jarObsPath
}
`, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY, acceptance.HW_DLI_FLINK_JAR_OBS_PATH, name)
}

func testAccFlinkJarJob_basic_step1(name string, agencyName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_dli_flinkjar_job" "test" {
  name                   = "%[2]s"
  queue_name             = "%[4]s"
  entrypoint             = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args        = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/%[6]s/output"
  obs_bucket             = "%[6]s"
  log_enabled            = true
  description            = "description of dli job"
  feature                = "basic"
  flink_version          = "%[5]s"
  cu_num                 = 6
  manager_cu_num         = 2
  parallel_num           = 2
  tm_cu_num              = 2
  tm_slot_num            = 1
  smn_topic              = huaweicloud_smn_topic.test.name
  restart_when_exception = true
  resume_checkpoint      = true
  resume_max_num         = 2
  checkpoint_path        = "%[6]s/"
  execution_agency_urn   = "%[7]s"

  runtime_config = {
    "max_records_per_file" = "10"
  }

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_dli_flinkjar_job" "test_resource_config_v2" {
  name            = "%[2]s_resource_config_version"
  queue_name      = "%[4]s"
  entrypoint      = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/%[6]s/output"
  obs_bucket      = "%[6]s"
  log_enabled     = true
  feature         = "basic"
  flink_version   = "%[5]s"

  runtime_config = {
    "max_records_per_file" = "10"
  }

  checkpoint_enabled      = true
  checkpoint_mode         = "2"
  checkpoint_interval     = 60
  execution_agency_urn    = "%[7]s"
  resource_config_version = "v2"

  resource_config {
    max_slot        = 4
    parallel_number = 3

    job_manager_resource_spec {
      cpu    = 1
      memory = "8G"
    }

    task_manager_resource_spec {
      cpu    = 2
      memory = "6G"
    }
  }
}
`, testAccFlinkJarJobResource_base(name), name, acceptance.HW_REGION_NAME, acceptance.HW_DLI_GENERAL_QUEUE_NAME,
		acceptance.HW_DLI_FLINK_VERSION, acceptance.HW_DLI_FLINK_JAR_OBS_BUCKET_NAME, agencyName)
}

func testAccFlinkJarJobResource_basic_step2(name string, agencyUpdateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_dli_flinkjar_job" "test" {
  name                   = "%[2]s"
  queue_name             = "%[4]s"
  entrypoint             = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args        = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/%[6]s/output"
  obs_bucket             = "%[6]s"
  log_enabled            = true
  description            = "Update description of dli job."
  feature                = "basic"
  flink_version          = "%[5]s"
  cu_num                 = 6
  manager_cu_num         = 2
  parallel_num           = 2
  tm_cu_num              = 2
  tm_slot_num            = 1
  smn_topic              = huaweicloud_smn_topic.test.name
  restart_when_exception = true
  resume_checkpoint      = false
  resume_max_num         = 1
  execution_agency_urn   = "%[7]s"

  runtime_config = {
    "max_records_per_file" = "5"
  }

  tags = {
    owner = "terraform"
  }
}

resource "huaweicloud_dli_flinkjar_job" "test_resource_config_v2" {
  name            = "%[2]s_resource_config_version"
  queue_name      = "%[4]s"
  entrypoint      = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args = "--output.path obs://${var.ak}:${var.sk}@obs.%[3]s.myhuaweicloud.com/%[6]s/output"
  obs_bucket      = "%[6]s"
  log_enabled     = true
  feature         = "basic"
  flink_version   = "%[5]s"

  runtime_config = {
    "max_records_per_file" = "10"
  }

  checkpoint_enabled      = false
  execution_agency_urn    = "%[7]s"
  resource_config_version = "v2"

  resource_config {
    max_slot        = 3
    parallel_number = 2

    job_manager_resource_spec {
      cpu    = 2
      memory = "6G"
    }

    task_manager_resource_spec {
      cpu    = 1
      memory = "2G"
    }
  }
}

resource "huaweicloud_dli_flinkjar_job" "default" {
  name          = "%[2]s_default"
  queue_name    = "%[4]s"
  entrypoint    = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  feature       = "basic"
  flink_version = "%[5]s"
  cu_num        = 2

  # The code logic verifies that 'resource_config' is empty.
  resource_config {
    job_manager_resource_spec {}
  }
}
`, testAccFlinkJarJobResource_base(name), name, acceptance.HW_REGION_NAME, acceptance.HW_DLI_GENERAL_QUEUE_NAME,
		acceptance.HW_DLI_FLINK_VERSION, acceptance.HW_DLI_FLINK_JAR_OBS_BUCKET_NAME, agencyUpdateName)
}
