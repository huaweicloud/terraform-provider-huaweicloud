package dcs

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

func getCustomTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCustomTemplate: Query DCS custom template
	var (
		getCustomTemplateHttpUrl = "v2/{project_id}/config-templates/{template_id}?type=user"
		getCustomTemplateProduct = "dcs"
	)
	getCustomTemplateClient, err := cfg.NewServiceClient(getCustomTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	getCustomTemplatePath := getCustomTemplateClient.Endpoint + getCustomTemplateHttpUrl
	getCustomTemplatePath = strings.ReplaceAll(getCustomTemplatePath, "{project_id}", getCustomTemplateClient.ProjectID)
	getCustomTemplatePath = strings.ReplaceAll(getCustomTemplatePath, "{template_id}", state.Primary.ID)

	getCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCustomTemplateResp, err := getCustomTemplateClient.Request("GET", getCustomTemplatePath, &getCustomTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS custom template: %s", err)
	}

	getCustomTemplateRespBody, err := utils.FlattenResponse(getCustomTemplateResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS custom template: %s", err)
	}

	return getCustomTemplateRespBody, nil
}

func TestAccCustomTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_custom_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "params.#", "1"),
					resource.TestCheckResourceAttr(rName, "params.0.param_name", "timeout"),
					resource.TestCheckResourceAttr(rName, "params.0.param_value", "200"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "engine"),
					resource.TestCheckResourceAttrSet(rName, "engine_version"),
					resource.TestCheckResourceAttrSet(rName, "cache_mode"),
					resource.TestCheckResourceAttrSet(rName, "product_type"),
					resource.TestCheckResourceAttrSet(rName, "storage_type"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testCustomTemplate_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "params.#", "1"),
					resource.TestCheckResourceAttr(rName, "params.0.param_name", "maxmemory-policy"),
					resource.TestCheckResourceAttr(rName, "params.0.param_value", "allkeys-lru"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_id", "source_type", "params"},
			},
		},
	})
}

func testCustomTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dcs_custom_template" "test" {
  template_id = "16"
  source_type = "sys"
  name        = "%s"
  description = "test custom policy"

  params {
    param_name  = "timeout"
    param_value = "200"
  }
}
`, name)
}

func testCustomTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dcs_custom_template" "test" {
  template_id = "16"
  source_type = "sys"
  name        = "%s"
  description = "test custom policy update"

  params {
    param_name  = "maxmemory-policy"
    param_value = "allkeys-lru"
  }
}
`, name)
}
