package apm

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

func getAppResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "apm", acceptance.HW_REGION_NAME)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "apm", conf.Region, "v1/apm2/access-keys")
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.GetAkSkListVO{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	for _, v := range rlt.AccessAkSkModels {
		if v.Ak == state.Primary.ID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
}

func TestAccApmAkSk_basic(t *testing.T) {
	var instance entity.AccessAkSkVO
	resourceName := "huaweicloud_apm_aksk.aksk_basic"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getAppResourceFunc,
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
				Config: testApmAkSk_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "create aksk by terraform"),
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

func testApmAkSk_basic() string {
	return `resource "huaweicloud_apm_aksk" "aksk_basic" {
  description = "create aksk by terraform"
}`
}
