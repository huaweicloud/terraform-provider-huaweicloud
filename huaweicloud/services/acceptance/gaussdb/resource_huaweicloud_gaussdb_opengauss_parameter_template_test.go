package gaussdb

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

func getOpenGaussParameterTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB OpenGauss parameter template: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB OpenGauss parameter template: %s", err)
	}

	return getRespBody, nil
}

func TestAccOpenGaussParameterTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussParameterTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test terraform description"),
					resource.TestCheckResourceAttr(rName, "engine_version", "8.201"),
					resource.TestCheckResourceAttr(rName, "instance_mode", "independent"),
					resource.TestCheckResourceAttr(rName, "parameters.#", "1"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "audit_system_object"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "100"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.need_restart"),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.data_type"),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.description"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_configuration_id", "parameters"},
			},
		},
	})
}

func TestAccOpenGaussParameterTemplate_copy(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussParameterTemplate_copy(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test terraform description"),
					resource.TestCheckResourceAttrPair(rName, "engine_version",
						"huaweicloud_gaussdb_opengauss_parameter_template.source", "engine_version"),
					resource.TestCheckResourceAttrPair(rName, "instance_mode",
						"huaweicloud_gaussdb_opengauss_parameter_template.source", "instance_mode"),
					resource.TestCheckResourceAttrPair(rName, "source_configuration_id",
						"huaweicloud_gaussdb_opengauss_parameter_template.source", "id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_configuration_id", "engine_version", "instance_mode"},
			},
		},
	})
}

func testOpenGaussParameterTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  name           = "%[1]s"
  description    = "test terraform description"
  engine_version = "8.201"
  instance_mode  = "independent"

  parameters {
    name  = "audit_system_object"
    value = "100"
  }
}
`, name)
}

func testOpenGaussParameterTemplate_copy(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_parameter_template" "source" {
  name           = "%[1]s_source"
  description    = "test terraform description"
  engine_version = "8.201"
  instance_mode  = "independent"

  parameters {
    name  = "audit_system_object"
    value = "100"
  }
}

resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  name                    = "%[1]s"
  description             = "test terraform description"
  source_configuration_id = huaweicloud_gaussdb_opengauss_parameter_template.source.id
}
`, name)
}
