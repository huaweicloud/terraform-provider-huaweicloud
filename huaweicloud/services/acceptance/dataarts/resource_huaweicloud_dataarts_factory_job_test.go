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
