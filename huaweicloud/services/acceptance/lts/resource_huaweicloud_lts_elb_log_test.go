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

func getElbResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "elb", acceptance.HW_REGION_NAME)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "elb", conf.Region, "v3/"+conf.GetProjectID(conf.Region)+
			"/elb/logtanks/"+state.Primary.ID)
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.CreateLogtankResponse{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccLtsElbComp_basic(t *testing.T) {
	var instance entity.CreateLogtankResponse
	resourceName := "huaweicloud_elb_log.elb_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getElbResourceFunc,
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
				Config: tesElbComp_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "log_topic_id", "1fab9f03-76d9-4d28-8c89-a90339d3e109"),
					resource.TestCheckResourceAttr(resourceName, "log_group_id", "f93f9fe0-de75-49c7-8eaf-a490b7ed202f"),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_id", "d2b7a742-fb38-49f2-b2d8-473a97a18434"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func tesElbComp_basic() string {
	return `
resource "huaweicloud_elb_log" "elb_1" {
  log_topic_id        = "1fab9f03-76d9-4d28-8c89-a90339d3e109"
  log_group_id      = "f93f9fe0-de75-49c7-8eaf-a490b7ed202f"
  loadbalancer_id    = "d2b7a742-fb38-49f2-b2d8-473a97a18434"
}`
}
