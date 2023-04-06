package smn

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

func getSmnMessageTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getMessageTemplate: Query SMN message template
	var (
		getMessageTemplateHttpUrl = "v2/{project_id}/notifications/message_template/{message_template_id}"
		getMessageTemplateProduct = "smn"
	)
	getMessageTemplateClient, err := cfg.NewServiceClient(getMessageTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN Client: %s", err)
	}

	getMessageTemplatePath := getMessageTemplateClient.Endpoint + getMessageTemplateHttpUrl
	getMessageTemplatePath = strings.ReplaceAll(getMessageTemplatePath, "{project_id}",
		getMessageTemplateClient.ProjectID)
	getMessageTemplatePath = strings.ReplaceAll(getMessageTemplatePath, "{message_template_id}",
		state.Primary.ID)

	getMessageTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getMessageTemplateResp, err := getMessageTemplateClient.Request("GET",
		getMessageTemplatePath, &getMessageTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SMN message template: %s", err)
	}
	return utils.FlattenResponse(getMessageTemplateResp)
}

func TestAccSmnMessageTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_smn_message_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSmnMessageTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSmnMessageTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protocol", "default"),
					resource.TestCheckResourceAttr(rName, "content",
						"this is a test content, contains {content1} and {content2}"),
					resource.TestCheckResourceAttr(rName, "tag_names.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_names.0", "content1"),
					resource.TestCheckResourceAttr(rName, "tag_names.1", "content2"),
				),
			},
			{
				Config: testSmnMessageTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protocol", "default"),
					resource.TestCheckResourceAttr(rName, "content",
						"this is a test update content, contains {content1}, {content2} and {content3}"),
					resource.TestCheckResourceAttr(rName, "tag_names.#", "3"),
					resource.TestCheckResourceAttr(rName, "tag_names.0", "content1"),
					resource.TestCheckResourceAttr(rName, "tag_names.1", "content2"),
					resource.TestCheckResourceAttr(rName, "tag_names.2", "content3"),
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

func testSmnMessageTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_message_template" "test" {
  name     = "%s"
  protocol = "default"
  content  = "this is a test content, contains {content1} and {content2}"
}
`, name)
}

func testSmnMessageTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_message_template" "test" {
  name     = "%s"
  protocol = "default"
  content  = "this is a test update content, contains {content1}, {content2} and {content3}"
}
`, name)
}
