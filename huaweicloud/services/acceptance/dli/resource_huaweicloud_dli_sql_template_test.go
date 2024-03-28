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

func getSQLTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSQLTemplate: Query the SQLTemplate.
	var (
		getSQLTemplateHttpUrl = "v1.0/{project_id}/sqls"
		getSQLTemplateProduct = "dli"
	)
	getSQLTemplateClient, err := cfg.NewServiceClient(getSQLTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getSQLTemplatePath := getSQLTemplateClient.Endpoint + getSQLTemplateHttpUrl
	getSQLTemplatePath = strings.ReplaceAll(getSQLTemplatePath, "{project_id}", getSQLTemplateClient.ProjectID)

	if v, ok := state.Primary.Attributes["name"]; ok {
		getSQLTemplatePath += fmt.Sprintf("?keyword=%s", v)
	}

	getSQLTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSQLTemplateResp, err := getSQLTemplateClient.Request("GET", getSQLTemplatePath, &getSQLTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SQLTemplate: %s", err)
	}
	getSQLTemplateRespBody, err := utils.FlattenResponse(getSQLTemplateResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SQLTemplate: %s", err)
	}

	jsonPath := fmt.Sprintf("sqls[?sql_id=='%s']|[0]", state.Primary.ID)
	sqlTemplate := utils.PathSearch(jsonPath, getSQLTemplateRespBody, nil)
	if sqlTemplate == nil {
		return nil, fmt.Errorf("no data found")
	}
	return sqlTemplate, nil
}

func TestAccSQLTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_sql_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "sql", "select * from t1"),
					resource.TestCheckResourceAttr(rName, "group", "demo"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
				),
			},
			{
				Config: testSQLTemplate_basic_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "sql", "select * from t2"),
					resource.TestCheckResourceAttr(rName, "group", "demo_2"),
					resource.TestCheckResourceAttr(rName, "description", "This is an updated demo"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
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

func testSQLTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_sql_template" "test" {
  name        = "%s"
  sql         = "select * from t1"
  group       = "demo"
  description ="This is a demo"
}
`, name)
}

func testSQLTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_sql_template" "test" {
  name        = "%s"
  sql         = "select * from t2"
  group       = "demo_2"
  description = "This is an updated demo"
}
`, name)
}
