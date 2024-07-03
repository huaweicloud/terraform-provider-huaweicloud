package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEndpointWhiteListFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := endpoints.ListOpts{
		InstanceId: state.Primary.ID,
	}
	resp, err := endpoints.ListPermissions(client, opts)

	var (
		domainId          = cfg.DomainID
		initialPermission = "iam:domain::" + domainId
	)
	var whiteLists []string
	for _, endpointPermission := range resp {
		if endpointPermission.Permission != initialPermission {
			whiteLists = append(whiteLists, endpointPermission.Permission)
		}
	}
	if err == nil && len(whiteLists) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, err
}

func TestAccEndpointWhiteList_basic(t *testing.T) {
	var (
		permissions []endpoints.EndpointPermission

		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_apig_endpoint_whitelist.test"

		rc = acceptance.InitResourceCheck(
			rName,
			&permissions,
			getEndpointWhiteListFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointWhiteList_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "whitelists.#", "2"),
				),
			},
			{
				Config: testAccEndpointWhiteList_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "whitelists.#", "3"),
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

func testAccEndpointWhiteList_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_endpoint_whitelist" "test" {
  instance_id = "%[1]s"
  whitelists  = [
    "iam:domain::1cc2018e40394f7c9692f1713e76234d",
    "iam:domain::2cc2018e40394f7c9692f1713e76234d",
  ]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccEndpointWhiteList_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_endpoint_whitelist" "test" {
  instance_id = "%[1]s"
  whitelists  = [
    "iam:domain::3cc2018e40394f7c9692f1713e76234d",
    "iam:domain::4cc2018e40394f7c9692f1713e76234d",
    "iam:domain::5cc2018e40394f7c9692f1713e76234d",
  ]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
