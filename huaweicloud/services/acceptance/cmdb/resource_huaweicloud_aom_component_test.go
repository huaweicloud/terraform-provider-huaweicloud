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

func getCompResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "cmdb", acceptance.HW_REGION_NAME)

	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "cmdb", conf.Region, "v1/components/"+state.Primary.Attributes["id"])
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.ComponentVo{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccAomComp_basic(t *testing.T) {
	var instance entity.BizAppVo
	var instanceName = acceptance.RandomAccResourceName()
	var modelId = "APPLICATION-ID"
	resourceName := "huaweicloud_aom_component.comp_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getCompResourceFunc,
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
				Config: tesAomComp_basic(modelId, instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "description", "component description"),
					resource.TestCheckResourceAttr(resourceName, "model_id", modelId),
					resource.TestCheckResourceAttr(resourceName, "model_type", "APPLICATION"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"model_id", "model_type"},
			},
			{
				Config: tesAomComp_updated(modelId, instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "component description1"),
				),
			},
		},
	})
}

func tesAomComp_basic(id, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_component" "comp_1" {
  description        = "component description"
  model_id           = "%s"
  model_type         = "APPLICATION"
  name               = "%s"
}`, id, name)
}

func tesAomComp_updated(id, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_component" "comp_1" {
  description        = "component description1"
  model_id           = "%s"
  model_type         = "APPLICATION"
  name               = "%s"
}`, id, name)
}
