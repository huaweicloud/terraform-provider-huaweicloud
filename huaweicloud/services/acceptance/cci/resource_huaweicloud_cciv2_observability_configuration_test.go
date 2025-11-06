package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getV2ObservabilityConfigurationResourceFunc(conf *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}

	httpUrl := "v1/observabilityconfiguration"
	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	enable := utils.PathSearch("event", getRespBody, nil)
	if enable == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return enable, nil
}

func TestAccV2ObservabilityConfiguration_basic(t *testing.T) {
	var ns namespaces.Namespace
	resourceName := "huaweicloud_cciv2_observability_configuration.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getV2ObservabilityConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV2ObservabilityConfiguration_basic(true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "event.enable", "true"),
				),
			},
			{
				Config: testAccV2ObservabilityConfiguration_basic(false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "event.enable", "false"),
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

func testAccV2ObservabilityConfiguration_basic(enable bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_cciv2_observability_configuration" "test" {
  event {
    enable = %v
  }
}
`, enable)
}
