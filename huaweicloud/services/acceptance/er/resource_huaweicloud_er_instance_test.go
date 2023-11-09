package er

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query the Enterprise router instance detail
	var (
		getInstanceHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		getInstanceProduct = "er"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{id}", state.Primary.ID)

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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic_step1(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_er_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testInstance_basic_step2(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
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

func testInstance_basic_step1(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = %[2]d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, bgpAsNum)
}

func testInstance_basic_step2(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = %[2]d

  tags = {
    foo    = "baar"
    newkey = "value"
  }
}
`, name, bgpAsNum)
}
