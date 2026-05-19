package geminidb

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

func getParameterTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/{config_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GeminiDB parameter template: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccGeminiDBParameterTemplate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_geminidb_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBParameterTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "datastore_version_name"),
					resource.TestCheckResourceAttrSet(resourceName, "datastore_name"),
				),
			},
			{
				Config: testAccGeminiDBParameterTemplate_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test configuration update"),
					resource.TestCheckResourceAttr(resourceName, "values.request_timeout_in_ms", "10000"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"datastore", "values"},
			},
		},
	})
}

func TestAccGeminiDBParameterTemplate_basedOnInstance(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_geminidb_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBParameterTemplate_basedOnInstance(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "configuration based on instance"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func testAccGeminiDBParameterTemplate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_parameter_template" "test" {
  name = "%s"

  datastore {
    type    = "cassandra"
    version = "3.11"
    mode    = "CloudNativeCluster"
  }

  values = {
    request_timeout_in_ms = "20000"
  }
}
`, rName)
}

func testAccGeminiDBParameterTemplate_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_parameter_template" "test" {
  name        = "%s"
  description = "test configuration update"

  datastore {
    type    = "cassandra"
    version = "3.11"
    mode    = "CloudNativeCluster"
  }

  values = {
    request_timeout_in_ms = "10000"
  }
}
`, rName)
}

func testAccGeminiDBParameterTemplate_basedOnInstance(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_parameter_template" "test" {
  name        = "%[2]s"
  description = "configuration based on instance"
  instance_id = huaweicloud_geminidb_instance.test.id
}
`, testAccGeminiDbInstance_basic(rName), rName)
}
