package cfw

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDNSResolutionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v1/{project_id}/dns/servers"
		product = "cfw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += "?offset=0&limit=100"
	path += fmt.Sprintf("&fw_instance_id=%s", state.Primary.ID)

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		path,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, fmt.Errorf("error retrieving CFW DNS resolution configuration: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshaling CFW DNS resolution configuration: %s", err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling CFW DNS resolution configuration: %s", err)
	}

	dnsServers := utils.PathSearch("data[?is_applied==`1`]", respBody, make([]interface{}, 0)).([]interface{})
	if len(dnsServers) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccDNSResolution_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_cfw_dns_resolution.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSResolutionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSResolution_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "default_dns_servers.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "health_check_domain_name"),
				),
			},
			{
				Config: testAccDNSResolution_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "default_dns_servers.#", "2"),
					resource.TestCheckResourceAttr(rName, "custom_dns_servers.#", "1"),
					resource.TestCheckResourceAttr(rName, "health_check_domain_name", "www.baidu.com"),
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

func testAccDNSResolution_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_dns_resolution" "test" {
  fw_instance_id      = "%s"
  default_dns_servers = ["8.8.8.8"]
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testAccDNSResolution_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_dns_resolution" "test" {
  fw_instance_id           = "%s"
  default_dns_servers      = ["8.8.8.8","114.114.114.114"]
  custom_dns_servers       = ["199.85.126.10"]
  health_check_domain_name = "www.baidu.com"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
