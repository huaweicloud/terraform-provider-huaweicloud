package tms

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	tags "github.com/chnsz/golangsdk/openstack/tms/v1/resourcetags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/tms"
)

func getResourceTagsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.TmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TMS v2 client: %s", err)
	}

	var (
		projectId       = state.Primary.Attributes["project_id"]
		resourcesLen, _ = strconv.Atoi(state.Primary.Attributes["resources.#"])
		tagsLen, _      = strconv.Atoi(state.Primary.Attributes["tags.%"])
		tagsConfigured  = false
	)

	for i := 0; i < resourcesLen; i++ {
		resourceId := state.Primary.Attributes[fmt.Sprintf("resources.%d.resource_id", i)]
		opts := tags.QueryOpts{
			ResourceId:   resourceId,
			ResourceType: state.Primary.Attributes[fmt.Sprintf("resources.%d.resource_type", i)],
			ProjectId:    projectId,
		}
		resp, err := tags.Get(client, opts)
		if err != nil {
			return nil, fmt.Errorf("error query resource (%s) tags: %s", resourceId, err)
		}
		actualTags := tms.FlattenTagsToMap(resp)
		if len(actualTags) != tagsLen {
			return nil, fmt.Errorf("some tags were not set successfully")
		}
		if len(actualTags) > 0 {
			tagsConfigured = true
		}
	}
	if !tagsConfigured {
		return nil, golangsdk.ErrDefault404{}
	}
	return tagsConfigured, nil
}

func TestAccResourceTags_basic(t *testing.T) {
	var (
		tagsConfigured bool

		rName    = "huaweicloud_tms_resource_tags.test"
		basicCfg = testAccResourceTags_base()
		rc       = acceptance.InitResourceCheck(rName, &tagsConfigured, getResourceTagsFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTags_basic_step1(basicCfg),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "resources.0.resource_type", "DNS_public_zone"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.resource_id", "huaweicloud_dns_zone.test", "id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccResourceTags_basic_step2(basicCfg),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "resources.0.resource_type", "instance"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.resource_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(rName, "tags.creator", "terraform"),
				),
			},
		},
	})
}

func testAccResourceTags_base() string {
	var (
		name         = acceptance.RandomAccResourceName()
		nameWithDash = acceptance.RandomAccResourceNameWithDash()
		bgpAsNum     = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dns_zone" "test" {
  name      = "%[2]s.com."
  email     = "jdoe@example.com"
  ttl       = 3000
  zone_type = "public"

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  
  name = "%[1]s"
  asn  = "%[3]d"

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
`, name, nameWithDash, bgpAsNum)
}

func testAccResourceTags_basic_step1(basicCfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_tms_resource_tags" "test" {
  project_id = "%[2]s"

  resources {
    resource_type = "DNS_public_zone"
    resource_id   = huaweicloud_dns_zone.test.id
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, basicCfg, acceptance.HW_PROJECT_ID)
}

func testAccResourceTags_basic_step2(basicCfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_tms_resource_tags" "test" {
  project_id = "%[2]s"

  resources {
    resource_type = "instance"
    resource_id   = huaweicloud_er_instance.test.id
  }

  tags = {
    foo     = "baar"
    creator = "terraform"
  }
}
`, basicCfg, acceptance.HW_PROJECT_ID)
}
