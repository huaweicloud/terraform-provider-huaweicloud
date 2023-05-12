package cfw

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

func getServiceGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getServiceGroup: Query the CFW service group detail
	var (
		getServiceGroupHttpUrl = "v1/{project_id}/service-sets/{id}"
		getServiceGroupProduct = "cfw"
	)
	getServiceGroupClient, err := cfg.NewServiceClient(getServiceGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	getServiceGroupPath := getServiceGroupClient.Endpoint + getServiceGroupHttpUrl
	getServiceGroupPath = strings.ReplaceAll(getServiceGroupPath, "{project_id}", getServiceGroupClient.ProjectID)
	getServiceGroupPath = strings.ReplaceAll(getServiceGroupPath, "{id}", state.Primary.ID)

	getServiceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getServiceGroupResp, err := getServiceGroupClient.Request("GET", getServiceGroupPath, &getServiceGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ServiceGroup: %s", err)
	}
	return utils.FlattenResponse(getServiceGroupResp)
}

func TestAccServiceGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_service_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getServiceGroupResourceFunc,
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
				Config: testServiceGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
				),
			},
			{
				Config: testServiceGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"object_id"},
			},
		},
	})
}

func testServiceGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_service_group" "test" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%s"
  description = "terraform test"
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testServiceGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_service_group" "test" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%s-update"
  description = "terraform test update"
}
`, testAccDatasourceFirewalls_basic(), name)
}
