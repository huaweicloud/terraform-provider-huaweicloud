package evs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getVolumeMetadataResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		mErr    *multierror.Error
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v3/{project_id}/volumes/{volume_id}/metadata/{key}"
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}
	metadataInput := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	metadataResult := make(map[string]interface{})

	for key := range metadataInput {
		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{volume_id}", state.Primary.ID)
		requestPath = strings.ReplaceAll(requestPath, "{key}", key)

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", requestPath, &requestOpt)
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error retrieving metadata for key %s: %s", key, err))
			continue
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error formatting metadata for key %s: %s", key, err))
			continue
		}
		path := fmt.Sprintf("meta.%s", key)
		metadataResult[key] = utils.PathSearch(path, respBody, nil)
	}

	return metadataResult, mErr
}

func TestAccEVSVolumeMetadata_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_evs_volume_metadata.test"
		volumeName   = acceptance.RandomAccResourceNameWithDash()
	)
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getVolumeMetadataResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeMetadata_basic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key3", "value3"),
				),
			},
			{
				Config: testAccVolumeMetadata_update(volumeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "new_value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key3", "value3"),
				),
			},
			{
				Config: testAccVolumeMetadata_update_null(volumeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "0"),
				),
			},
			{
				Config: testAccVolumeMetadata_update_for_destroy(volumeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key3", "value3"),
				),
			},
		},
	})
}

func testAccVolumeMetadata_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "Created by volume metadata test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}
`, name)
}

func testAccVolumeMetadata_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume_metadata" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  metadata = {
    key1 = "value1"
    key2 = "value2"
    key3 = "value3"
  }
}
`, testAccVolumeMetadata_base(name))
}

func testAccVolumeMetadata_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume_metadata" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  metadata = {
    key1 = "new_value1"
    key2 = "value2"
    key3 = "value3"
  }
}
`, testAccVolumeMetadata_base(name))
}

func testAccVolumeMetadata_update_null(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume_metadata" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  metadata = {}
}
`, testAccVolumeMetadata_base(name))
}

func testAccVolumeMetadata_update_for_destroy(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume_metadata" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  metadata = {
    key1 = "value1"
    key2 = "value2"
    key3 = "value3"
  }
}
`, testAccVolumeMetadata_base(name))
}
