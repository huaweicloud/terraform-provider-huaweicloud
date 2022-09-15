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

func getEnvResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "cmdb", acceptance.HW_REGION_NAME)

	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "cmdb", conf.Region, "v1/environments/"+state.Primary.Attributes["id"])
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.EnvVo{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccAomEnv_basic(t *testing.T) {
	var instance entity.BizAppVo
	var instanceName = acceptance.RandomAccResourceName()
	var componentId = "COMPONENT-ID"
	resourceName := "huaweicloud_aom_environment.env1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getEnvResourceFunc,
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
				Config: tesAomEnv_basic(componentId, instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "env_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "description", "environment description"),
					resource.TestCheckResourceAttr(resourceName, "component_id", componentId),
					resource.TestCheckResourceAttr(resourceName, "env_type", "DEVos_type"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "LINUX"),
					resource.TestCheckResourceAttr(resourceName, "region", "cn-north-7"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tesEnvApp_updated(componentId, instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "environment description1"),
				),
			},
		},
	})
}

func tesAomEnv_basic(id, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_environment" "env1" {
  component_id       = "%s"
  description        = "environment description"
  env_name           = "%s"
  env_type           = "DEV"
  os_type            = "LINUX"
  region             = "cn-north-7"
  register_type      = "API"
}`, id, name)
}

func tesEnvApp_updated(id, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_environment" "env1" {
  component_id       = "%s"
  description        = "environment description1"
  env_name           = "%s"
  env_type           = "DEV"
  os_type            = "LINUX"
  region             = "cn-north-7"
  register_type      = "API"
}`, id, name)
}
