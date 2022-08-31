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

func getAccessRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "lts", acceptance.HW_REGION_NAME)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "lts", conf.Region, "v2/"+conf.GetProjectID(conf.Region)+
			"/lts/aom-mapping/"+state.Primary.ID)
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := make([]entity.AomMappingRequestInfo, 0)
	err = json.Unmarshal(body, &rlt)
	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	if len(rlt) == 0 {
		return rlt, fmt.Errorf("resource is not exists")
	}
	return rlt, nil
}

func TestAccessRule_basic(t *testing.T) {
	var instance []entity.AomMappingRequestInfo
	resourceName := "huaweicloud_lts_access_rule.accessrule_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getAccessRuleResourceFunc,
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
				Config: tesAccessRule_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "rule_name", "cui-rule-1"),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "8caff49f-0317-11ed-ba1c-0255ac100b0d"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", "ictest"),
					resource.TestCheckResourceAttr(resourceName, "name_space", "cui-namespace"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"files"},
			},
		},
	})
}

func tesAccessRule_basic() string {
	return `
resource "huaweicloud_lts_access_rule" "accessrule_1" {
  rule_name     = "cui-rule-1"
  cluster_id    = "8caff49f-0317-11ed-ba1c-0255ac100b0d"
  cluster_name  = "ictest"
  name_space    = "cui-namespace"
  deployments   = ["__ALL_DEPLOYMENTS__"]
  files {
	file_name = "__ALL_FILES__"
	log_stream_info {
		target_log_group_id = "65808f94-a011-42f3-8acf-2fc99923ece1"
		target_log_group_name = "k8s-log-7430ae60-16fe-11ed-8b6a-0255ac100b0b"
		target_log_stream_id = "4dbcaf7a-cadc-4bf5-ae59-b0d7513d2283"
		target_log_stream_name = "stdout-7430ae60-16fe-11ed-8b6a-0255ac100b0b"
	}
  }
}`
}
