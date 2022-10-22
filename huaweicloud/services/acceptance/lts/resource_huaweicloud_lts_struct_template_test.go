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
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lts_struct_template.template_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getLtsStructTemplateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesLtsStructTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"huaweicloud_lts_group.group_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id",
						"huaweicloud_lts_stream.stream_1", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "demo_log"),
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

func tesLtsStructTemplate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "group_1" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "huaweicloud_lts_stream" "stream_1" {
  group_id    = huaweicloud_lts_group.group_1.id
  stream_name = "%[1]s"
}

resource "huaweicloud_lts_struct_template" "template_1" {
  log_group_id  = huaweicloud_lts_group.group_1.id
  log_stream_id = huaweicloud_lts_stream.stream_1.id
  template_type = "custom"
}`, rName)
}
