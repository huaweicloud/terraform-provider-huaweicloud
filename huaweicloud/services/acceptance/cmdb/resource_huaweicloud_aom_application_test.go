package cmdb

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"testing"
)

var c *httpclient_go.HttpClientGo

func init() {
	config := &config.Config{AccessKey: acceptance.HW_ACCESS_KEY, SecretKey: acceptance.HW_SECRET_KEY}
	c, _ = httpclient_go.NewHttpClientGo(config)
}

func getAppResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	url := "https://aom.cn-north-7.myhuaweicloud.com/v1/applications/" + state.Primary.Attributes["id"]

	c.WithMethod(httpclient_go.MethodGet).WithUrl(url)
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.BizAppVo{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccAomApp_basic(t *testing.T) {
	var instance entity.BizAppVo
	var instanceName = acceptance.RandomAccResourceName()
	var instanceNameUpdate = instanceName + "_update"
	resourceName := "huaweicloud_aom_application.app_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getAppResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesAomApp_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "description", "application description"),
					resource.TestCheckResourceAttr(resourceName, "display_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "eps_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "CONSOLE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tesAomApp_updated(instanceName, instanceNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", instanceNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "application description"),
				),
			},
		},
	})
}

func tesAomApp_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_application" "app_1" {
  description        = "application description"
  display_name       = "%s"
  name               = "%s"
  eps_id             = "0"
  register_type      = "CONSOLE"
}`, instanceName, instanceName)
}

func tesAomApp_updated(instanceName, instanceNameUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_application" "app_1" {
  description        = "application description"
  display_name       = "%s"
  name               = "%s"
  eps_id             = "0"
  register_type      = "CONSOLE"
}`, instanceName, instanceNameUpdate)
}
