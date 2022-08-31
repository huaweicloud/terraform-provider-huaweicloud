package cmdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
)

func getLtsStructTemplateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "lts", acceptance.HW_REGION_NAME)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "lts", conf.Region, "v2/"+conf.GetProjectID(conf.Region)+
			"/lts/struct/template?logGroupId="+state.Primary.Attributes["log_group_id"]+"&logStreamId="+state.Primary.Attributes["log_stream_id"])
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}
	body = body[1 : len(body)-1]
	body2 := strings.Replace(string(body), `\\\`, "**", -1)
	body3 := strings.Replace(body2, `\`, "", -1)
	body4 := strings.Replace(body3, "**", `\`, -1)
	rlt := &entity.ShowStructTemplateResponse{}
	err = json.Unmarshal([]byte(body4), rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccLtsStructTemplate_basic(t *testing.T) {
	var instance entity.ShowStructTemplateResponse
	resourceName := "huaweicloud_lts_struct_template.struct_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getLtsStructTemplateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckInternal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesLtsStructTemplate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "log_group_id", "e5e5a6de-354d-45af-a7da-0fb91e9a3796"),
					resource.TestCheckResourceAttr(resourceName, "log_stream_id", "45bbeee7-2144-4d40-9c80-ba452298b6b8"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_type"},
				ImportStateIdFunc:       testAccLtsStructImportStateIdFunc(),
			},
		},
	})
}

func testAccLtsStructImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var logGroupId, logStreamId, id string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_lts_struct_template" {
				logGroupId = rs.Primary.Attributes["log_group_id"]
				logStreamId = rs.Primary.Attributes["log_stream_id"]
				id = rs.Primary.ID
			}
		}
		if logGroupId == "" || logStreamId == "" || id == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", id, logGroupId, logStreamId)
		}
		return fmt.Sprintf("%s/%s/%s", id, logGroupId, logStreamId), nil
	}
}

func tesLtsStructTemplate_basic() string {
	return `
resource "huaweicloud_lts_struct_template" "struct_1" {
  log_group_id        = "e5e5a6de-354d-45af-a7da-0fb91e9a3796"
  log_stream_id      = "45bbeee7-2144-4d40-9c80-ba452298b6b8"
  template_type    = "custom"
}`
}
