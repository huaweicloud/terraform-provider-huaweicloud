package dli

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGlobalVariableResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGlobalVariable: Query the Global variable.
	var (
		getGlobalVariableHttpUrl = "v1.0/{project_id}/variables"
		getGlobalVariableProduct = "dli"
	)
	getGlobalVariableClient, err := cfg.NewServiceClient(getGlobalVariableProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getGlobalVariablePath := getGlobalVariableClient.Endpoint + getGlobalVariableHttpUrl
	getGlobalVariablePath = strings.ReplaceAll(getGlobalVariablePath, "{project_id}", getGlobalVariableClient.ProjectID)

	getGlobalVariableResp, err := pagination.ListAllItems(
		getGlobalVariableClient,
		"offset",
		getGlobalVariablePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving DLI global variable: %s", err)
	}

	getGlobalVariableRespJson, err := json.Marshal(getGlobalVariableResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DLI global variable: %s", err)
	}
	var getGlobalVariableRespBody interface{}
	err = json.Unmarshal(getGlobalVariableRespJson, &getGlobalVariableRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DLI global variable: %s", err)
	}

	jsonPath := fmt.Sprintf("global_vars[?var_name=='%s']|[0]", state.Primary.ID)
	globalVariable := utils.PathSearch(jsonPath, getGlobalVariableRespBody, nil)
	if globalVariable == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return globalVariable, nil
}

func TestAccGlobalVariable_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_global_variable.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGlobalVariableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGlobalVariable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "value", "abc"),
				),
			},
			{
				Config: testGlobalVariable_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "value", "abcd"),
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

func testGlobalVariable_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_global_variable" "test" {
  name = "%s"
  value = "abc"
}
`, name)
}

func testGlobalVariable_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_global_variable" "test" {
  name = "%s"
  value = "abcd"
}
`, name)
}
