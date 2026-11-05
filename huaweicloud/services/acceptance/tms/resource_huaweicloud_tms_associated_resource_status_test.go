package tms

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

func getTmsAssociatedResourceStatusResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("tms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TMS client: %s", err)
	}

	getPath := client.Endpoint + "v2/tms/associated-resources/status"
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving TMS associated resource status: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", getRespBody, "")
	if status == "" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method: "GET",
				URL:    "/v2/tms/associated-resources/status",
			},
		}
	}

	return status, nil
}

func TestAccTmsAssociatedResourceStatus_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_tms_associated_resource_status.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTmsAssociatedResourceStatusResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// The delete operation of this resource is not supported. The Delete function only returns a warning
		// and removes the resource from the state, so no CheckDestroy verification is required here.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// The step covers the create behavior (PUT /v2/tms/associated-resources/status) with `enable`
				// and the read behavior (GET /v2/tms/associated-resources/status).
				Config: testAccTmsAssociatedResourceStatus_basic("enable"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "enable"),
				),
			},
			{
				// The step covers the update behavior (PUT /v2/tms/associated-resources/status),
				// switching the status from `enable` to `disable`.
				Config: testAccTmsAssociatedResourceStatus_basic("disable"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "disable"),
				),
			},
			{
				// The step covers the update behavior again, switching the status back from `disable` to `enable`.
				Config: testAccTmsAssociatedResourceStatus_basic("enable"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "enable"),
				),
			},
			{
				// The step covers the import behavior. The resource ID is a fixed value and the status
				// is refreshed by the read behavior after the import.
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTmsAssociatedResourceStatus_basic(status string) string {
	return fmt.Sprintf(`
resource "huaweicloud_tms_associated_resource_status" "test" {
  status = "%s"
}
`, status)
}
