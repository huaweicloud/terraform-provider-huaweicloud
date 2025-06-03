package evs

import (
	"errors"
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

func getSnapshotMetadataResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region         = acceptance.HW_REGION_NAME
		product        = "evs"
		metadataResult = make(map[string]interface{})
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + "v3/{project_id}/snapshots/{snapshot_id}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	metadataMap := make(map[string]string)
	for k, v := range state.Primary.Attributes {
		if strings.HasPrefix(k, "metadata.") && k != "metadata.%" {
			key := strings.TrimPrefix(k, "metadata.")
			metadataMap[key] = v
		}
	}

	for key := range metadataMap {
		getPath := fmt.Sprintf("%s/%s", requestPath, key)
		getResp, err := client.Request("GET", getPath, &requestOpt)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				continue
			}

			return nil, fmt.Errorf("error retrieving EVS snapshot metadata: %s", err)
		}

		respBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		metaResp := utils.PathSearch("meta", respBody, make(map[string]interface{})).(map[string]interface{})
		for k, v := range metaResp {
			metadataResult[k] = v
		}
	}

	// Only when all key queries in `metadata` return `404`, it is considered that the resource does not exist.
	if len(metadataMap) > 0 && len(metadataResult) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return metadataResult, nil
}

func TestAccSnapshotMetadata_basic(t *testing.T) {
	var (
		snapshotMetadata interface{}
		rName            = acceptance.RandomAccResourceName()
		resourceName     = "huaweicloud_evs_snapshot_metadata.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&snapshotMetadata,
		getSnapshotMetadataResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotMetadata_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "snapshot_id", "huaweicloud_evsv3_snapshot.test", "id"),
				),
			},
			{
				Config: testAccSnapshotMetadata_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "metadata.test_key", "test_value"),
				),
			},
			{
				Config: testAccSnapshotMetadata_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "metadata.#", "0"),
				),
			},
			{
				Config: testAccSnapshotMetadata_update3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
				),
			},
		},
	})
}

func testAccSnapshotMetadata_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}

resource "huaweicloud_evsv3_snapshot" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  name      = "%[1]s"
}
`, rName)
}

func testAccSnapshotMetadata_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = huaweicloud_evsv3_snapshot.test.id

  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, testAccSnapshotMetadata_base(rName))
}

func testAccSnapshotMetadata_update1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = huaweicloud_evsv3_snapshot.test.id

  metadata = {
    foo      = "bar_update"
    test_key = "test_value"
    key      = "value_update"
  }
}
`, testAccSnapshotMetadata_base(rName))
}

func testAccSnapshotMetadata_update2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = huaweicloud_evsv3_snapshot.test.id

  metadata = {}
}
`, testAccSnapshotMetadata_base(rName))
}

func testAccSnapshotMetadata_update3(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = huaweicloud_evsv3_snapshot.test.id

  metadata = {
    foo = "bar"
  }
}
`, testAccSnapshotMetadata_base(rName))
}
