package eip

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

func getGlobalEipSegmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geip", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GEIP client: %s", err)
	}

	getHttpUrl := "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getPath = strings.ReplaceAll(getPath, "{global_eip_segment_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccGlobalEipSegment_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_global_eip_segment.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGlobalEipSegmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGlobalEipSegment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(rName, "geip_pool_name",
						"${data.huaweicloud_global_eip_pools.all.geip_pools.0.name}"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "access_site",
						"${data.huaweicloud_global_eip_pools.all.geip_pools.0.access_site}"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					resource.TestCheckResourceAttr(rName, "tags.#", "2"),
					resource.TestCheckResourceAttr(rName, "tags.0.key", "key1"),
					resource.TestCheckResourceAttr(rName, "tags.0.value", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.1.key", "key2"),
					resource.TestCheckResourceAttr(rName, "tags.1.value", "value2"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "isp"),
					resource.TestCheckResourceAttrSet(rName, "ip_version"),
					resource.TestCheckResourceAttrSet(rName, "cidr"),
					resource.TestCheckResourceAttrSet(rName, "freezen"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "is_pre_paid"),
					resource.TestCheckResourceAttrSet(rName, "is_charged"),
				),
			},
			{
				Config: testGlobalEipSegment_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttr(rName, "tags.#", "1"),
					resource.TestCheckResourceAttr(rName, "tags.0.key", "key_updated"),
					resource.TestCheckResourceAttr(rName, "tags.0.value", "value_updated"),
				),
			},
			{
				Config: testGlobalEipSegment_update2(""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", ""),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "tags.#", "0"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"mask"},
			},
		},
	})
}

func testGlobalEipSegment_base() string {
	return `
# Currently, only this global eip pool can create a global IP segment
data "huaweicloud_global_eip_pools" "all" {
  access_site = "cn-south-guangzhou"
  name        = "bgp_segment_default"
}
`
}

func testGlobalEipSegment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_global_eip_segment" "test" {
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  mask                  = 29
  name                  = "%[2]s"
  description           = "description test"
  enterprise_project_id = "%[3]s"

  tags {
    key   = "key1"
    value = "value1"
  }
  tags {
    key   = "key2"
    value = "value2"
  }
}
`, testGlobalEipSegment_base(), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testGlobalEipSegment_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_global_eip_segment" "test" {
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  mask                  = 29
  name                  = "%[2]s"
  description           = "description test"
  enterprise_project_id = "%[3]s"

  tags {
    key   = "key_updated"
    value = "value_updated"
  }
}
`, testGlobalEipSegment_base(), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testGlobalEipSegment_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_global_eip_segment" "test" {
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  mask                  = 29
  name                  = "%[2]s"
  description           = ""
  enterprise_project_id = "%[3]s"
}
`, testGlobalEipSegment_base(), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
