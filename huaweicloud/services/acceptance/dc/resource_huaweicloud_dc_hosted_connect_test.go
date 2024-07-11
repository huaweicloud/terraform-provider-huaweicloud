package dc

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

func getHostedConnectResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getHostedConnectHttpUrl = "v3/{project_id}/dcaas/hosted-connects/{id}"
		getHostedConnectProduct = "dc"
	)
	getHostedConnectClient, err := cfg.NewServiceClient(getHostedConnectProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DC client: %s", err)
	}

	getHostedConnectPath := getHostedConnectClient.Endpoint + getHostedConnectHttpUrl
	getHostedConnectPath = strings.ReplaceAll(getHostedConnectPath, "{project_id}", getHostedConnectClient.ProjectID)
	getHostedConnectPath = strings.ReplaceAll(getHostedConnectPath, "{id}", state.Primary.ID)

	getHostedConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getHostedConnectResp, err := getHostedConnectClient.Request("GET", getHostedConnectPath, &getHostedConnectOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving hosted connect: %s", err)
	}

	getHostedConnectRespBody, err := utils.FlattenResponse(getHostedConnectResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving hosted connect: %s", err)
	}

	getHostedConnectRespBody = utils.PathSearch("hosted_connect", getHostedConnectRespBody, nil)
	if getHostedConnectRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getHostedConnectRespBody, nil
}

// This resource needs other user's tenant_id and host_id, so skip the acceptance test
func TestAccHostedConnect_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dc_hosted_connect.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostedConnectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcHostedConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostedConnect_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "10"),
					resource.TestCheckResourceAttr(rName, "hosting_id", acceptance.HW_DC_HOSTTING_ID),
					resource.TestCheckResourceAttr(rName, "vlan", "441"),
					resource.TestCheckResourceAttr(rName, "resource_tenant_id", acceptance.HW_DC_RESOURCE_TENANT_ID),
					resource.TestCheckResourceAttr(rName, "status", "BUILD"),
				),
			},
			{
				Config: testHostedConnect_basic_update(name + "update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"update"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo update"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "12"),
					resource.TestCheckResourceAttr(rName, "hosting_id", acceptance.HW_DC_HOSTTING_ID),
					resource.TestCheckResourceAttr(rName, "vlan", "441"),
					resource.TestCheckResourceAttr(rName, "resource_tenant_id", acceptance.HW_DC_RESOURCE_TENANT_ID),
					resource.TestCheckResourceAttr(rName, "status", "BUILD"),
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

func testHostedConnect_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_hosted_connect" "test" {
  name               = "%s"
  description        = "This is a demo"
  resource_tenant_id = "%s"
  hosting_id         = "%s"
  vlan               = 441
  bandwidth          = 10
}
`, name, acceptance.HW_DC_RESOURCE_TENANT_ID, acceptance.HW_DC_HOSTTING_ID)
}

func testHostedConnect_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_hosted_connect" "test" {
  name               = "%s"
  description        = "This is a demo update"
  resource_tenant_id = "%s"
  hosting_id         = "%s"
  vlan               = 441
  bandwidth          = 12
}
`, name, acceptance.HW_DC_RESOURCE_TENANT_ID, acceptance.HW_DC_HOSTTING_ID)
}
