package er

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query the Enterprise router instance detail
	var (
		getInstanceHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		getInstanceProduct = "er"
	)
	getInstanceClient, err := config.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.Replace(getInstancePath, "{project_id}", getInstanceClient.ProjectID, -1)
	getInstancePath = strings.Replace(getInstancePath, "{id}", state.Primary.ID, -1)

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Instance: %s", err)
	}
	return utils.FlattenResponse(getInstanceResp)
}

func TestAccInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_er_instance.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "availability_zones.0", acceptance.HW_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testInstance_basic(name, bgpAsNum+1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum+1)),
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

func testInstance_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
resource "huaweicloud_er_instance" "test" {
  availability_zones = ["%[1]s"]

  name = "%[2]s"
  asn  = %[3]d
}
`, acceptance.HW_AVAILABILITY_ZONE, name, bgpAsNum)
}
