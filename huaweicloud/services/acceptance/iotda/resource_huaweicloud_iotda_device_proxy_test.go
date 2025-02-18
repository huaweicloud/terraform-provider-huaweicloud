package iotda

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDeviceProxyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-proxies/{proxy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{proxy_id}", state.Primary.ID)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA device proxy: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccDeviceProxy_basic(t *testing.T) {
	var (
		deviceProxyObj interface{}
		name           = acceptance.RandomAccResourceName()
		rName          = "huaweicloud_iotda_device_proxy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&deviceProxyObj,
		getDeviceProxyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource only supports standard and enterprise version IoTDA instances.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceProxy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "devices.#", "2"),
					resource.TestCheckResourceAttr(rName, "effective_time_range.0.start_time", "20881010T121212Z"),
					resource.TestCheckResourceAttr(rName, "effective_time_range.0.end_time", "20881015T121212Z"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
				),
			},
			{
				Config: testDeviceProxy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "devices.#", "3"),
					resource.TestCheckResourceAttr(rName, "effective_time_range.0.start_time", "20991010T121212Z"),
					resource.TestCheckResourceAttr(rName, "effective_time_range.0.end_time", "20991015T121212Z"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDeviceProxy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device" "test" {
  count = 3

  node_id    = format("%[2]s_%%d", count.index)
  name       = format("%[2]s_%%d", count.index)
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
}

resource "huaweicloud_iotda_device_proxy" "test" {
  depends_on = [
    huaweicloud_iotda_device.test
  ]

  space_id = huaweicloud_iotda_space.test.id
  name     = "%[2]s"
  devices  = slice(huaweicloud_iotda_device.test[*].id, 0, 2)

  effective_time_range {
    start_time = "20881010T121212Z"
    end_time   = "20881015T121212Z"
  }
}
`, testProduct_basic(name), name)
}

func testDeviceProxy_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device" "test" {
  count = 3

  node_id    = format("%[2]s_%%d", count.index)
  name       = format("%[2]s_%%d", count.index)
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
}

resource "huaweicloud_iotda_device_proxy" "test" {
  depends_on = [
    huaweicloud_iotda_device.test
  ]

  space_id = huaweicloud_iotda_space.test.id
  name     = "%[2]s_update"
  devices  = slice(huaweicloud_iotda_device.test[*].id, 0, 3)

  effective_time_range {
    start_time = "20991010T121212Z"
    end_time   = "20991015T121212Z"
  }
}
`, testProduct_basic(name), name)
}
