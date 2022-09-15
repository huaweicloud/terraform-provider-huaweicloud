package cmdb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
)

func getCmdbResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "cmdb", acceptance.HW_REGION_NAME)

	opts := entity.PageResourceListParam{
		CiId:     state.Primary.Attributes["env_id"],
		CiType:   "environment",
		Keywords: map[string]string{"RESOURCE_ID": state.Primary.ID},
	}

	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "cmdb", conf.Region, "v1/resource/"+state.Primary.Attributes["rf_resource_type"]+
			"/type/"+state.Primary.Attributes["type"]+"/ci-relationships").WithBody(opts)
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.ReadResourceResponse{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccAomResource_basic(t *testing.T) {
	var instance entity.ReadResourceResponse
	var cmdbResourceName = "ecs-cmdb3"
	var cmdbResourceId = "e90e1556-685f-400b-9f81-8c2eeda0c75a"
	name := "huaweicloud_aom_cmdb_resource_relationships.ecs"

	rc := acceptance.InitResourceCheck(
		name,
		&instance,
		getCmdbResourceFunc,
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
				Config: testAomResource_basic(cmdbResourceId, cmdbResourceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(name, "resource_id", cmdbResourceId),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
				ImportStateIdFunc:       testAccInstanceImportStateIdFunc(),
			},
		},
	})
}

func testAomResource_basic(resourceId string, resourceName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_cmdb_resource_relationships" "ecs" {
	rf_resource_type              = "ecs"
	type             = "cloudservers"
	env_id                     = "5f0e3558e54a45cfb8644391d30cd706"
	resource_id    = "%s"
	resource_name            = "%s"
	resource_region	="cn-north-7"
}`, resourceId, resourceName)
}

func testAccInstanceImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var rfResourceType, cmdbResourceType, envId, resourceId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_aom_cmdb_resource_relationships" {
				rfResourceType = rs.Primary.Attributes["rf_resource_type"]
				cmdbResourceType = rs.Primary.Attributes["type"]
				envId = rs.Primary.Attributes["env_id"]
				resourceId = rs.Primary.Attributes["resource_id"]
			}
		}
		if rfResourceType == "" || cmdbResourceType == "" || envId == "" || resourceId == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s/%s", rfResourceType, cmdbResourceType, envId, resourceId)
		}
		return fmt.Sprintf("%s/%s/%s/%s", rfResourceType, cmdbResourceType, envId, resourceId), nil
	}
}
