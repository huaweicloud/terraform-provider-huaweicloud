package ga

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getHealthCheckResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/health-checks/{id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GA health check: %s", err)
	}
	return utils.FlattenResponse(resp)
}

func TestAccHealthCheck_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_ga_health_check.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHealthCheckResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHealthCheck_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "interval", "10"),
					resource.TestCheckResourceAttr(rName, "max_retries", "5"),
					resource.TestCheckResourceAttr(rName, "port", "8001"),
					resource.TestCheckResourceAttr(rName, "timeout", "10"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "endpoint_group_id",
						"huaweicloud_ga_endpoint_group.test", "id"),
				),
			},
			{
				Config: testHealthCheck_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "interval", "20"),
					resource.TestCheckResourceAttr(rName, "max_retries", "10"),
					resource.TestCheckResourceAttr(rName, "port", "8002"),
					resource.TestCheckResourceAttr(rName, "timeout", "20"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
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

func testHealthCheck_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ga_health_check" "test" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  enabled           = true
  interval          = 10
  max_retries       = 5
  port              = 8001
  timeout           = 10
}
`, testEndpointGroup_basic(name))
}

func testHealthCheck_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ga_health_check" "test" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  enabled           = true
  interval          = 20
  max_retries       = 10
  port              = 8002
  timeout           = 20
}
`, testEndpointGroup_basic(name))
}
