package dns

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

func getDNSCustomLineResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDNSCustomLine: Query DNS custom line
	var (
		getDNSCustomLineHttpUrl = "v2.1/customlines"
		getDNSCustomLineProduct = "dns"
	)
	getDNSCustomLineClient, err := cfg.NewServiceClient(getDNSCustomLineProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS Client: %s", err)
	}

	getDNSCustomLinePath := getDNSCustomLineClient.Endpoint + getDNSCustomLineHttpUrl
	getDNSCustomLinePath += fmt.Sprintf("?line_id=%s", state.Primary.ID)

	getDNSCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSCustomLineResp, err := getDNSCustomLineClient.Request("GET", getDNSCustomLinePath, &getDNSCustomLineOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS custom line: %s", err)
	}

	getDNSCustomLineRespBody, err := utils.FlattenResponse(getDNSCustomLineResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten DNS custom line response: %s", err)
	}
	return flattenCustomLineResponseBody(getDNSCustomLineRespBody, state.Primary.ID)
}

func flattenCustomLineResponseBody(resp interface{}, id string) (interface{}, error) {
	if resp == nil {
		return nil, fmt.Errorf("custom line response is empty")
	}
	curJson := utils.PathSearch("lines", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		lineId := utils.PathSearch("line_id", v, "")
		if id == lineId.(string) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("the target custom line (%s) not exist", id)
}

func TestAccDNSCustomLine_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dns_custom_line.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSCustomLineResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSCustomLine_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "ip_segments.0", "100.100.100.100-100.100.100.100"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testDNSCustomLine_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "ip_segments.#", "2"),
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

func testDNSCustomLine_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  name        = "%s"
  description = "test description"
  ip_segments = ["100.100.100.100-100.100.100.100"]
}
`, name)
}

func testDNSCustomLine_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  name        = "%s_update"
  ip_segments = ["100.100.100.102-100.100.100.102", "100.100.100.101-100.100.100.101"]
}
`, name)
}
