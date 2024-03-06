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

func getDataStandardResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDataStandard: query DataArts Architecture data standard
	var (
		getDataStandardHttpUrl = "v2/{project_id}/design/standards"
		getDataStandardProduct = "dataarts"
	)
	getDataStandardClient, err := cfg.NewServiceClient(getDataStandardProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getDataStandardBasePath := getDataStandardClient.Endpoint + getDataStandardHttpUrl
	getDataStandardBasePath = strings.ReplaceAll(getDataStandardBasePath, "{project_id}", getDataStandardClient.ProjectID)

	var currentTotal int
	var dataStandard interface{}
	directoryID := state.Primary.Attributes["directory_id"]
	getDataStandardPath := getDataStandardBasePath + buildGetDataStandardQueryParams(directoryID, currentTotal)
	for {
		getDataStandardOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"workspace": state.Primary.Attributes["workspace_id"],
			},
		}

		getDataStandardResp, err := getDataStandardClient.Request("GET", getDataStandardPath, &getDataStandardOpt)

		if err != nil {
			return nil, fmt.Errorf("error retrieving DataArts Architecture data standard")
		}

		getDataStandardRespBody, err := utils.FlattenResponse(getDataStandardResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("data.value.records", getDataStandardRespBody, make([]interface{}, 0))
		dataStandard = utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", state.Primary.ID), records, nil)
		if dataStandard != nil {
			break
		}
		total := utils.PathSearch("data.value.total", getDataStandardRespBody, float64(0)).(float64)
		currentTotal += len(records.([]interface{}))
		if currentTotal == int(total) {
			break
		}
		getDataStandardPath = getDataStandardBasePath + buildGetDataStandardQueryParams(directoryID, currentTotal)
	}

	if dataStandard == nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture data standard")
	}

	return dataStandard, nil
}

func buildGetDataStandardQueryParams(directoryID string, offset int) string {
	return fmt.Sprintf("?directory_id=%v&limit=100&offset=%v", directoryID, offset)
}

func TestAccDataStandard_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_architecture_data_standard.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataStandardResourceFunc,
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
				Config: testDataStandard_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "directory_id",
						"huaweicloud_dataarts_architecture_directory.test", "id"),
					resource.TestCheckResourceAttr(rName, "values.#", "3"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testDataStandard_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "directory_id",
						"huaweicloud_dataarts_architecture_directory.test", "id"),
					resource.TestCheckResourceAttr(rName, "values.#", "4"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataStandardImportState(rName),
			},
		},
	})
}

func testDataStandard_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_data_standard" "test" {
  workspace_id = "%[2]s"
  directory_id = huaweicloud_dataarts_architecture_directory.test.id

  values {
    fd_name  = "nameCh"
    fd_value = "%[3]s"
  }

  values {
    fd_name  = "nameEn"
    fd_value = "%[3]s"
  }

  values {
    fd_name  = "description"
    fd_value = "this is a terraform"
  }
}
`, testAccArchitectureDirectory_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataStandard_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_data_standard" "test" {
  workspace_id = "%[2]s"
  directory_id = huaweicloud_dataarts_architecture_directory.test.id

  values {
    fd_name  = "nameCh"
    fd_value = "%[3]s"
  }

  values {
    fd_name  = "nameEn"
    fd_value = "%[3]s"
  }

  values {
    fd_name  = "dataType"
    fd_value = "STRING"
  }

  values {
    fd_name  = "description"
    fd_value = ""
  }
}
`, testAccArchitectureDirectory_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataStandardImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		dataStandard, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, dataStandard)
		}

		var workspaceID string
		if workspaceID = dataStandard.Primary.Attributes["workspace_id"]; workspaceID == "" {
			return "", fmt.Errorf("attribute (workspace_id) of Resource (%s) not found: %s", name, dataStandard)
		}
		return fmt.Sprintf("%s/%s", workspaceID, dataStandard.Primary.ID), nil
	}
}
