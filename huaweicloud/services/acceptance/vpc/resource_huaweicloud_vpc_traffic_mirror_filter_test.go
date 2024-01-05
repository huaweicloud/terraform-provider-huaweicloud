package vpc

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

func getTrafficMirrorFilterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return "", fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	getTrafficMirrorFilterHttpUrl := "vpc/traffic-mirror-filters/" + state.Primary.ID
	getTrafficMirrorFilterPath := client.ResourceBaseURL() + getTrafficMirrorFilterHttpUrl
	getTrafficMirrorFilterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getTrafficMirrorFilterResp, err := client.Request("GET", getTrafficMirrorFilterPath, &getTrafficMirrorFilterOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving traffic mirror filter: %s", err)
	}

	return utils.FlattenResponse(getTrafficMirrorFilterResp)
}

func TestAccTrafficMirrorFilter_basic(t *testing.T) {
	var (
		trafficMirrorFilter interface{}
		name                = acceptance.RandomAccResourceNameWithDash()
		updatedName         = acceptance.RandomAccResourceNameWithDash()
		resourceName        = "huaweicloud_vpc_traffic_mirror_filter.test"
		desc                = "create VPC traffic mirror filter"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&trafficMirrorFilter,
			getTrafficMirrorFilterResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTrafficMirrorFilter_base(name, desc),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", desc),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccTrafficMirrorFilter_base(updatedName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func testAccTrafficMirrorFilter_base(name string, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_traffic_mirror_filter" "test" {
  name        = "%[1]s"
  description = "%[2]s"
}
`, name, description)
}
