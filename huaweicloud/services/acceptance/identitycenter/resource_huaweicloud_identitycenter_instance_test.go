package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityCenterInstanceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listHttpUrl = "v1/instances"
		listProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	instanceId := state.Primary.ID

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center instance: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	instance := utils.PathSearch(fmt.Sprintf("instances|[?instance_id =='%s']|[0]", instanceId), listRespBody, nil)
	if instance == nil {
		return nil, fmt.Errorf("error get Identity Center instance")
	}
	return instance, nil
}

func TestAccIdentityCenterInstance_basic(t *testing.T) {
	var obj interface{}

	region := acceptance.HW_REGION_NAME
	rName := "huaweicloud_identitycenter_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterInstance_basic(region),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "instance_urn"),
					resource.TestCheckResourceAttr(rName, "alias", ""),
				),
			},
			{
				Config: testIdentityCenterInstance_basic_update(region),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "instance_urn"),
					resource.TestCheckResourceAttr(rName, "alias", "test"),
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

func testIdentityCenterInstance_basic(region string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_identitycenter_instance" "test" {
  depends_on = [huaweicloud_identitycenter_registered_region.test]
}
`, testRegisteredRegion_basic(region))
}

func testIdentityCenterInstance_basic_update(region string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_identitycenter_instance" "test" {
  depends_on = [huaweicloud_identitycenter_registered_region.test]
  alias      = "test"
}
`, testRegisteredRegion_basic(region))
}
