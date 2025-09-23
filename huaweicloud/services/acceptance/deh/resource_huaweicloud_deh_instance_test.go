package deh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDehInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v1.0/{project_id}/dedicated-hosts/{dedicated_host_id}"
		product = "deh"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DEH client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dedicated_host_id}", state.Primary.ID)

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getBackupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DEH instance: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccDehInstance_basic(t *testing.T) {
	var proxy instances.Proxy
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_deh_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&proxy,
		getDehInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDehInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "cn-southwest-242a"),
					resource.TestCheckResourceAttr(resourceName, "host_type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "auto_placement", "on"),
					resource.TestCheckResourceAttr(resourceName, "metadata.ha_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "available_vcpus"),
					resource.TestCheckResourceAttrSet(resourceName, "available_memory"),
					resource.TestCheckResourceAttrSet(resourceName, "allocated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_total"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_uuids.#"),
					resource.TestCheckResourceAttrSet(resourceName, "sys_tags.%"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.#"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.host_type"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.host_type_name"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.vcpus"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.cores"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.sockets"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.memory"),
					resource.TestCheckResourceAttrSet(resourceName, "host_properties.0.available_instance_capacities.#"),
					resource.TestCheckResourceAttrSet(resourceName,
						"host_properties.0.available_instance_capacities.0.flavor"),
				),
			},
			{
				Config: testDehInstance_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "auto_placement", "off"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner_update", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func testDehInstance_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_deh_instance" "test" {
  availability_zone = "cn-southwest-242a"
  name              = "%[1]s"
  host_type         = "s3"
  auto_placement    = "on"

  metadata = {
    "ha_enabled" = "true"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"

  tags = {
    "key"   = "value"
    "owner" = "terraform"
  }
}
`, rName)
}

func testDehInstance_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_deh_instance" "test" {
  availability_zone = "cn-southwest-242a"
  name              = "%[1]s"
  host_type         = "s3"
  auto_placement    = "off"

  metadata = {
    "ha_enabled" = "true"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"

  tags = {
    "key_update"   = "value_update"
    "owner_update" = "terraform_update"
  }
}
`, rName)
}
