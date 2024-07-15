package dli

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

func getFlinkTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getFlinkTemplate: Query the Flink template.
	var (
		getFlinkTemplateHttpUrl = "v1.0/{project_id}/streaming/job-templates"
		getFlinkTemplateProduct = "dli"
	)
	getFlinkTemplateClient, err := cfg.NewServiceClient(getFlinkTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getFlinkTemplatePath := getFlinkTemplateClient.Endpoint + getFlinkTemplateHttpUrl
	getFlinkTemplatePath = strings.ReplaceAll(getFlinkTemplatePath, "{project_id}", getFlinkTemplateClient.ProjectID)

	if v, ok := state.Primary.Attributes["name"]; ok {
		getFlinkTemplatePath += fmt.Sprintf("?limit=100&name=%s", v)
	}

	getFlinkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getFlinkTemplateResp, err := getFlinkTemplateClient.Request("GET", getFlinkTemplatePath, &getFlinkTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving FlinkTemplate: %s", err)
	}

	getFlinkTemplateRespBody, err := utils.FlattenResponse(getFlinkTemplateResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Flink template: %s", err)
	}

	jsonPath := fmt.Sprintf("template_list.templates[?template_id==`%s`]|[0]", state.Primary.ID)
	template := utils.PathSearch(jsonPath, getFlinkTemplateRespBody, nil)
	if template == nil {
		return nil, fmt.Errorf("no data found")
	}
	return template, nil
}

func TestAccFlinkTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_flink_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFlinkTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFlinkTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "sql", "select * from source_table"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "type", "flink_sql_job"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config: testFlinkTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "sql", "select * from source_table2"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo2"),
					resource.TestCheckResourceAttr(rName, "type", "flink_sql_job"),
					// The tags is ForceNew.
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testFlinkTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_flink_template" "test" {
  name        = "%s"
  type        = "flink_sql_job"
  sql         = "select * from source_table"
  description = "This is a demo"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testFlinkTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_flink_template" "test" {
  name        = "%s"
  type        = "flink_sql_job"
  sql         = "select * from source_table2"
  description = "This is a demo2"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}
