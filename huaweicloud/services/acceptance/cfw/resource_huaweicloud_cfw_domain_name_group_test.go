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

func getDomainNameGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDomainNameGroup: Query the CFW domain name group detail
	var (
		getDomainNameGroupHttpUrl = "v1/{project_id}/domain-sets"
		getDomainNameGroupProduct = "cfw"
	)
	getDomainNameGroupClient, err := cfg.NewServiceClient(getDomainNameGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	getDomainNameGroupPath := getDomainNameGroupClient.Endpoint + getDomainNameGroupHttpUrl
	getDomainNameGroupPath = strings.ReplaceAll(getDomainNameGroupPath, "{project_id}", getDomainNameGroupClient.ProjectID)
	getDomainNameGroupqueryParams := buildGetDomainNameGroupQueryParams(state)
	getDomainNameGroupPath += getDomainNameGroupqueryParams

	getDomainNameGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getDomainNameGroupResp, err := getDomainNameGroupClient.Request("GET", getDomainNameGroupPath, &getDomainNameGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DomainNameGroup: %s", err)
	}

	getDomainNameGroupRespBody, err := utils.FlattenResponse(getDomainNameGroupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DomainNameGroup: %s", err)
	}

	jsonPath := fmt.Sprintf("data.records[?set_id=='%s']|[0]", state.Primary.ID)
	getDomainNameGroupRespBody = utils.PathSearch(jsonPath, getDomainNameGroupRespBody, nil)
	if getDomainNameGroupRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getDomainNameGroupRespBody, nil
}

func TestAccDomainNameGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_domain_name_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDomainNameGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDomainNameGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "fw_instance_id",
						"huaweicloud_cfw_firewall.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "object_id",
						"huaweicloud_cfw_firewall.test", "protect_objects.0.object_id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "domain_names.0.domain_name", "www.cfw-test1.com"),
					resource.TestCheckResourceAttr(rName, "domain_names.0.description", "test domain 1"),
					resource.TestCheckResourceAttr(rName, "domain_names.1.domain_name", "www.cfw-test2.com"),
					resource.TestCheckResourceAttr(rName, "domain_names.1.description", "test domain 2"),
					resource.TestCheckResourceAttrSet(rName, "config_status"),
					resource.TestCheckResourceAttrSet(rName, "ref_count"),
				),
			},
			{
				Config: testDomainNameGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "domain_names.0.domain_name", "www.cfw-test2.com"),
					resource.TestCheckResourceAttr(rName, "domain_names.0.description", "test domain 2"),
					resource.TestCheckResourceAttr(rName, "domain_names.1.domain_name", "www.cfw-test3.com"),
					resource.TestCheckResourceAttr(rName, "domain_names.1.description", "test domain 3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDomainNameGroupImportState(rName),
			},
		},
	})
}

func testDomainNameGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_domain_name_group" "test" {
  fw_instance_id = huaweicloud_cfw_firewall.test.id
  object_id      = huaweicloud_cfw_firewall.test.protect_objects[0].object_id
  name           = "%s"
  type           = 0
  description    = "created by terraform"

  domain_names {
    domain_name = "www.cfw-test1.com"
    description = "test domain 1"
  }

  domain_names {
    domain_name = "www.cfw-test2.com"
    description = "test domain 2"
  }
}
`, testFirewall_basic(name), name)
}

func testDomainNameGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_domain_name_group" "test" {
  fw_instance_id = huaweicloud_cfw_firewall.test.id
  object_id      = huaweicloud_cfw_firewall.test.protect_objects[0].object_id
  name           = "%s_update"
  type           = 0
  description    = ""

  domain_names {
    domain_name = "www.cfw-test2.com"
    description = "test domain 2"
  }

  domain_names {
    domain_name = "www.cfw-test3.com"
    description = "test domain 3"
  }
}
`, testFirewall_basic(name), name)
}

func buildGetDomainNameGroupQueryParams(state *terraform.ResourceState) string {
	res := fmt.Sprintf("?limit=1024&offset=0&fw_instance_id=%v&object_id=%v",
		state.Primary.Attributes["fw_instance_id"], state.Primary.Attributes["object_id"])

	if state.Primary.Attributes["name"] != "" {
		res += fmt.Sprintf("&key_word=%v", state.Primary.Attributes["name"])
	}
	return res
}

func testDomainNameGroupImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["fw_instance_id"] == "" {
			return "", fmt.Errorf("Attribute (fw_instance_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["object_id"] == "" {
			return "", fmt.Errorf("Attribute (object_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["fw_instance_id"] + "/" + rs.Primary.Attributes["object_id"] + "/" +
			rs.Primary.ID, nil
	}
}
