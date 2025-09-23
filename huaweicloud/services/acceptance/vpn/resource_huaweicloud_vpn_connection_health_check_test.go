package vpn

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

func getConnectionHealthCheckResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getConnectionHealthCheck: Query the VPN ConnectionHealthCheck detail
	var (
		getConnectionHealthCheckHttpUrl = "v5/{project_id}/connection-monitors/{id}"
		getConnectionHealthCheckProduct = "vpn"
	)
	getConnectionHealthCheckClient, err := cfg.NewServiceClient(getConnectionHealthCheckProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN Client: %s", err)
	}

	getConnectionHealthCheckPath := getConnectionHealthCheckClient.Endpoint + getConnectionHealthCheckHttpUrl
	getConnectionHealthCheckPath = strings.ReplaceAll(getConnectionHealthCheckPath, "{project_id}", getConnectionHealthCheckClient.ProjectID)
	getConnectionHealthCheckPath = strings.ReplaceAll(getConnectionHealthCheckPath, "{id}", state.Primary.ID)

	getConnectionHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getConnectionHealthCheckResp, err := getConnectionHealthCheckClient.Request("GET", getConnectionHealthCheckPath, &getConnectionHealthCheckOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ConnectionHealthCheck: %s", err)
	}
	return utils.FlattenResponse(getConnectionHealthCheckResp)
}

func TestAccConnectionHealthCheck_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_connection_health_check.test"
	ipAddress := "172.16.1.4"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getConnectionHealthCheckResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConnectionHealthCheck_basic(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "connection_id",
						"huaweicloud_vpn_connection.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "destination_ip"),
					resource.TestCheckResourceAttrSet(rName, "source_ip"),
					resource.TestCheckResourceAttrSet(rName, "status"),
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

func testConnectionHealthCheck_basic(name, ipAddress string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_connection_health_check" "test" {
  connection_id = huaweicloud_vpn_connection.test.id
}
`, testConnection_basic(name, ipAddress))
}
