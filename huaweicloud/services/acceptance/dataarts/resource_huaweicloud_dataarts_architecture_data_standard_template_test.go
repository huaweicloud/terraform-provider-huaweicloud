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

func getDataStandardTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDataStandardTemplate: query DataArts Architecture data standard template
	var (
		getDataStandardTemplateHttpUrl = "v2/{project_id}/design/standards/templates"
		getDataStandardTemplateProduct = "dataarts"
	)
	getDataStandardTemplateClient, err := cfg.NewServiceClient(getDataStandardTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}

	getDataStandardTemplatePath := getDataStandardTemplateClient.Endpoint + getDataStandardTemplateHttpUrl
	getDataStandardTemplatePath = strings.ReplaceAll(getDataStandardTemplatePath, "{project_id}",
		getDataStandardTemplateClient.ProjectID)

	getDataStandardTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.ID},
	}

	getDataStandardTemplateResp, err := getDataStandardTemplateClient.Request("GET", getDataStandardTemplatePath,
		&getDataStandardTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture data standard template: %s", err)
	}

	getDataStandardTemplateRespBody, err := utils.FlattenResponse(getDataStandardTemplateResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture data standard template: %s", err)
	}

	hasTemplate := utils.PathSearch("data.value.hasTemplate", getDataStandardTemplateRespBody, false).(bool)
	if !hasTemplate {
		return nil, golangsdk.ErrDefault404{}
	}

	return getDataStandardTemplateRespBody, nil
}

func TestAccDataStandardTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_architecture_data_standard_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataStandardTemplateResourceFunc,
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
				Config: testDataStandardTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "optional_fields.#", "3"),
					resource.TestCheckResourceAttr(rName, "custom_fields.#", "1"),
					resource.TestCheckResourceAttr(rName, "custom_fields.0.fd_name", name),
					resource.TestCheckResourceAttr(rName, "custom_fields.0.optional_values", "111;222"),
					resource.TestCheckResourceAttr(rName, "custom_fields.0.required", "true"),
					resource.TestCheckResourceAttr(rName, "custom_fields.0.searchable", "true"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.id"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.created_by"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.updated_by"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "custom_fields.0.id"),
					resource.TestCheckResourceAttrSet(rName, "custom_fields.0.created_by"),
					resource.TestCheckResourceAttrSet(rName, "custom_fields.0.updated_by"),
					resource.TestCheckResourceAttrSet(rName, "custom_fields.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "custom_fields.0.updated_at"),
				),
			},
			{
				Config: testDataStandardTemplate_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "optional_fields.#", "5"),
					resource.TestCheckResourceAttr(rName, "custom_fields.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDataStandardTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_data_standard_template" "test" {
  workspace_id = "%[1]s"

  optional_fields {
    fd_name    = "dataLength"
    required   = false
    searchable = false
  }

  optional_fields {
    fd_name    = "hasAllowValueList"
    required   = false
    searchable = true
  }

  optional_fields {
    fd_name    = "allowList"
    required   = true
    searchable = false
  }

  custom_fields {
    fd_name         = "%[2]s"
    optional_values = "111;222"
    required        = true
    searchable      = true
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataStandardTemplate_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_data_standard_template" "test" {
  workspace_id = "%[1]s"

  optional_fields {
    fd_name    = "dataLength"
    required   = false
    searchable = false
  }

  optional_fields {
    fd_name    = "hasAllowValueList"
    required   = false
    searchable = true
  }

  optional_fields {
    fd_name    = "allowList"
    required   = true
    searchable = false
  }

  optional_fields {
    fd_name    = "dqcRule"
    required   = true
    searchable = false
  }

  optional_fields {
    fd_name    = "ruleOwner"
    required   = false
    searchable = true
  }

  custom_fields {
    fd_name         = "%[2]s_1"
    optional_values = "aaa;bbb"
    required        = true
    searchable      = false
  }

  custom_fields {
    fd_name         = "%[2]s_2"
    optional_values = "111;222"
    required        = false
    searchable      = true
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
