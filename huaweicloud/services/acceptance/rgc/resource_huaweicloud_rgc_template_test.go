package rgc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getTemplate: Query template RGC template via RFS API
	var (
		region              = acceptance.HW_REGION_NAME
		listTemplateHttpUrl = "v1/{project_id}/templates"
		listTemplateProduct = "rfs"
	)
	listTemplateClient, err := cfg.NewServiceClient(listTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	projectId := listTemplateClient.ProjectID

	var templateMeta interface{}
	var marker string

	for {
		listTemplatePath := listTemplateClient.Endpoint + listTemplateHttpUrl + buildTemplateQueryParams(marker)
		listTemplatePath = strings.ReplaceAll(listTemplatePath, "{project_id}", projectId)

		randUUID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("unable to generate ID: %s", err)
		}
		listTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Client-Request-Id": randUUID,
				"Content-Type":      "application/json",
				"X-Language":        "en-us",
			},
		}

		listTemplateResp, err := listTemplateClient.Request("GET", listTemplatePath, &listTemplateOpt)
		if err != nil {
			return nil, err
		}

		listTemplateRespBody, err := utils.FlattenResponse(listTemplateResp)
		if err != nil {
			return nil, err
		}

		templateMeta = utils.PathSearch(fmt.Sprintf("templates[?template_id=='%s']|[0]", state.Primary.ID), listTemplateRespBody, nil)

		if templateMeta != nil {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", listTemplateRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	if templateMeta == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return templateMeta, nil
}

func buildTemplateQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func TestAccTemplate_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_rgc_template.test"
	templateName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
			acceptance.TestAccPreCheckRGCTemplate(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTemplate_customized(templateName, acceptance.HW_RGC_TEMPLATE_BODY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "region"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "latest_version_id"),
				),
			},
		},
	})
}

func TestAccTemplate_predefined(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_rgc_template.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTemplate_predefined_vpc(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "region"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "latest_version_id"),
				),
			},
		},
	})
}

func testTemplate_customized(name string, body string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_template" "test" {
  template_name = "%[1]s"
  template_type = "customized"
  template_description = ""
  template_body = "%[2]s"
}
`, name, body)
}

func testTemplate_predefined_vpc() string {
	return `
resource "huaweicloud_rgc_template" "test" {
  template_name = "VPC"
  template_type = "predefined"
}
`
}
