package aom

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

func getCloudServiceAccessResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{prom_instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{provider}", state.Primary.Attributes["service"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM cloud service access: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening AOM cloud service access: %s", err)
	}

	// API will return 200 and nil if `instance_id` and `servcie` is invalid
	if getRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccCloudServiceAccess_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_aom_cloud_service_access.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCloudServiceAccessResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCloudServiceAccess_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_aom_prom_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "service", "OBS"),
					resource.TestCheckResourceAttr(resourceName, "tag_sync", "manual"),
				),
			},
			{
				Config: testCloudServiceAccess_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_aom_prom_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "service", "OBS"),
					resource.TestCheckResourceAttr(resourceName, "tag_sync", "auto"),
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

func testCloudServiceAccessBase(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "test" {
  prom_name             = "%[1]s"
  prom_type             = "CLOUD_SERVICE"
  enterprise_project_id = "0"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true

  tags = {
    owner = "terraform"
    key   = "value"
  }
}
`, name)
}

func testCloudServiceAccess_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_cloud_service_access" "test" {
  instance_id = huaweicloud_aom_prom_instance.test.id
  service     = "OBS"
  tag_sync    = "manual"

  tags {
    sync   = false
    key    = "key"
    values = [huaweicloud_obs_bucket.test.tags["key"]]
  }

  tags {
    sync   = true
    key    = "owner"
    values = [huaweicloud_obs_bucket.test.tags["owner"]]
  }
}
`, testCloudServiceAccessBase(name))
}

func testCloudServiceAccess_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_cloud_service_access" "test" {
  instance_id = huaweicloud_aom_prom_instance.test.id
  service     = "OBS"
  tag_sync    = "auto"

  tags {
    sync   = true
    key    = "key"
    values = [huaweicloud_obs_bucket.test.tags["key"]]
  }

  tags {
    sync   = false
    key    = "owner"
    values = [huaweicloud_obs_bucket.test.tags["owner"]]
  }
}
`, testCloudServiceAccessBase(name))
}
