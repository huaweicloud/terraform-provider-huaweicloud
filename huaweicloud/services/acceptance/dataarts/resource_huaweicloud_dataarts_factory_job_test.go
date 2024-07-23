package dataarts

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

func getFactoryJobResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getJobHttpUrl = "v1/{project_id}/jobs/{job_name}"
		getJobProduct = "dataarts-dlf"
	)
	getJobClient, err := cfg.NewServiceClient(getJobProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}

	getJobPath := getJobClient.Endpoint + getJobHttpUrl
	getJobPath = strings.ReplaceAll(getJobPath, "{project_id}", getJobClient.ProjectID)
	getJobPath = strings.ReplaceAll(getJobPath, "{job_name}", state.Primary.ID)

	getJobOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    state.Primary.Attributes["workspace_id"],
		},
	}

	getJobResp, err := getJobClient.Request("GET", getJobPath, &getJobOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Job: %s", err)
	}

	getJobRespBody, err := utils.FlattenResponse(getJobResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Job: %s", err)
	}

	return getJobRespBody, nil
}

func TestAccFactoryJob_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_factory_job.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFactoryJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsCdmName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFactoryJob_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "nodes.0.name", "Rest_client_"+name),
					resource.TestCheckResourceAttr(rName, "nodes.0.type", "RESTAPI"),
					resource.TestCheckResourceAttr(rName, "nodes.0.location.0.x", "10"),
					resource.TestCheckResourceAttr(rName, "nodes.0.location.0.y", "11"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.0.name", "url"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.0.value", "https://www.huaweicloud.com/"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.1.name", "method"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.1.value", "GET"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.2.name", "retry"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.2.value", "false"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.3.name", "requestMode"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.3.value", "sync"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.4.name", "securityAuthentication"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.4.value", "NONE"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.5.name", "agentName"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.5.value", acceptance.HW_DATAARTS_CDM_NAME),
					resource.TestCheckResourceAttr(rName, "schedule.0.type", "EXECUTE_ONCE"),
				),
			},
			{
				Config: testFactoryJob_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "nodes.0.name", "Rest_client_"+name),
					resource.TestCheckResourceAttr(rName, "nodes.0.type", "RESTAPI"),
					resource.TestCheckResourceAttr(rName, "nodes.0.location.0.x", "103"),
					resource.TestCheckResourceAttr(rName, "nodes.0.location.0.y", "113"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.0.name", "url"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.0.value", "https://www.huaweicloud.com/console/"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.1.name", "method"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.1.value", "GET"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.2.name", "retry"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.2.value", "false"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.3.name", "requestMode"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.3.value", "sync"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.4.name", "securityAuthentication"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.4.value", "NONE"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.5.name", "agentName"),
					resource.TestCheckResourceAttr(rName, "nodes.0.properties.5.value", acceptance.HW_DATAARTS_CDM_NAME),
					resource.TestCheckResourceAttr(rName, "schedule.0.type", "EXECUTE_ONCE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testFactoryJobImportState(rName),
			},
		},
	})
}

func testFactoryJob_basic(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[1]s"
  workspace_id = "%[2]s"
  process_type = "REAL_TIME"

  nodes {
    name = "Rest_client_%[1]s"
    type = "RESTAPI"
    location {
      x = 10
      y = 11
    }

    properties {
      name  = "url"
      value = "https://www.huaweicloud.com/"
    }

    properties {
      name  = "method"
      value = "GET"
    }

    properties {
      name  = "retry"
      value = "false"
    }

    properties {
      name  = "requestMode"
      value = "sync"
    }

    properties {
      name  = "securityAuthentication"
      value = "NONE"
    }

    properties {
      name  = "agentName"
      value = "%[3]s"
    }

  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CDM_NAME)
}

func testFactoryJob_basic_update(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[1]s"
  workspace_id = "%[2]s"
  process_type = "REAL_TIME"

  nodes {
    name = "Rest_client_%[1]s"
    type = "RESTAPI"
    location {
      x = 103
      y = 113
    }

    properties {
      name  = "url"
      value = "https://www.huaweicloud.com/console/"
    }

    properties {
      name  = "method"
      value = "GET"
    }

    properties {
      name  = "retry"
      value = "false"
    }

    properties {
      name  = "requestMode"
      value = "sync"
    }

    properties {
      name  = "securityAuthentication"
      value = "NONE"
    }

    properties {
      name  = "agentName"
      value = "%[3]s"
    }

  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CDM_NAME)
}

func testFactoryJobImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["workspace_id"] + "/" + rs.Primary.ID, nil
	}
}

func TestAccFactoryJob_batch_singleTask_job(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_factory_job.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFactoryJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFactoryJob_batch_singleTask_job_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(rName, "nodes.0.type", "DLISQL"),
					resource.TestCheckResourceAttr(rName, "schedule.0.type", "CRON"),
					resource.TestCheckResourceAttr(rName, "schedule.0.cron.0.expression", "0 0 0 * * ?"),
					resource.TestCheckResourceAttr(rName, "schedule.0.cron.0.depend_jobs.0.jobs.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "schedule.0.cron.0.depend_jobs.0.jobs.0",
						"huaweicloud_dataarts_factory_job.test2", "id"),
				),
			},
			{
				Config: testAccFactoryJob_batch_singleTask_job_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "schedule.0.cron.0.depend_jobs.0.jobs.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testFactoryJobImportState(rName),
			},
		},
	})
}

func testAccFactoryJob_batch_singleTask_job_base(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_dli_database" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dataarts_factory_job" "test2" {
  name         = "%[1]s_2"
  workspace_id = "%[2]s"
  process_type = "BATCH"

  nodes {
    name = "%[1]s_2"
    type = "DLISQL"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "database"
      value = huaweicloud_dli_database.test.name
    }

    properties {
      name  = "queueName"
      value = "default"
    }

    properties {
      name  = "statementOrScript"
      value = "SCRIPT"
    }

    properties {
      name  = "overloadRetryInterval"
      value = "300 sec"
    }

    properties {
      name  = "maximumOverloadRetries"
      value = "0"
    }
  }

  schedule {
    type = "CRON"

    cron {
      start_time = "2024-07-23T13:08:37+08"
      expression = "0 0 0 * * ?"
    }
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccFactoryJob_batch_singleTask_job_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[2]s_1"
  workspace_id = "%[3]s"
  process_type = "BATCH"
  
  nodes {
    name = "%[2]s_1"
    type = "DLISQL"

    location {
      x = 10
      y = 11
    }
  
    properties {
      name  = "database"
      value = huaweicloud_dli_database.test.name
    }
  
    properties {
      name  = "queueName"
      value = "default"
    }
  
    properties {
      name  = "statementOrScript"
      value = "SCRIPT"
    }

    properties {
      name  = "overloadRetryInterval"
      value = "300 sec"
    }

    properties {
      name  = "maximumOverloadRetries"
      value = "0"
    }
  }
  
  schedule {
    type = "CRON"

    cron {
      start_time = "2024-07-23T13:08:37+08"
      expression = "0 0 0 * * ?"

      depend_jobs {
        jobs = [huaweicloud_dataarts_factory_job.test2.id]
      }
    }
  }
}
`, testAccFactoryJob_batch_singleTask_job_base(name), name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccFactoryJob_batch_singleTask_job_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[2]s_1"
  workspace_id = "%[3]s"
  process_type = "BATCH"

  nodes {
    name = "%[2]s_1"
    type = "DLISQL"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "database"
      value = huaweicloud_dli_database.test.name
    }

    properties {
      name  = "queueName"
      value = "default"
    }

    properties {
      name  = "statementOrScript"
      value = "SCRIPT"
    }

    properties {
      name  = "overloadRetryInterval"
      value = "300 sec"
    }

    properties {
      name  = "maximumOverloadRetries"
      value = "0"
    }
  }

  schedule {
    type = "CRON"

    cron {
      start_time = "2024-07-23T13:08:37+08"
      expression = "0 0 0 * * ?"

      depend_jobs {
        jobs = []
      }
    }
  }
}
`, testAccFactoryJob_batch_singleTask_job_base(name), name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
