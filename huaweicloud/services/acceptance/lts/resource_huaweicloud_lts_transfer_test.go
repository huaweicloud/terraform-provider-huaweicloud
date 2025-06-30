package lts

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTransferResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getTransfer: Query the log transfer task.
	var (
		getTransferHttpUrl = "v2/{project_id}/transfers"
		getTransferProduct = "lts"
	)
	getTransferClient, err := cfg.NewServiceClient(getTransferProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getTransferPath := getTransferClient.Endpoint + getTransferHttpUrl
	getTransferPath = strings.ReplaceAll(getTransferPath, "{project_id}", getTransferClient.ProjectID)

	getTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getTransferResp, err := getTransferClient.Request("GET", getTransferPath, &getTransferOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving LTS transfer: %s", err)
	}

	getTransferRespBody, err := utils.FlattenResponse(getTransferResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving LTS transfer: %s", err)
	}

	jsonPath := fmt.Sprintf("log_transfers[?log_transfer_id =='%s']|[0]", state.Primary.ID)
	getTransferRespBody = utils.PathSearch(jsonPath, getTransferRespBody, nil)
	if getTransferRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getTransferRespBody, nil
}

func TestAccTransfer_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_lts_transfer.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTransferResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTransfer_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period", "3"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "hour"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.obs_bucket_name",
						"huaweicloud_obs_bucket.output", "bucket"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name", "lts_transfer_obs_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id", "Etc/GMT"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
				),
			},
			{
				Config: testTransfer_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period", "2"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "min"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.obs_bucket_name",
						"huaweicloud_obs_bucket.output", "bucket"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name", "lts_transfer_obs_2_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_2_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC-02:00"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id", "Etc/GMT+2"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
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

func testAccTransfer_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}
`, name)
}

func testTransfer_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "output" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = huaweicloud_obs_bucket.output.bucket
      obs_dir_prefix_name = "lts_transfer_obs_"
      obs_prefix_name     = "obs_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }
  }
}
`, testAccTransfer_base(name), name)
}

func testTransfer_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "output" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "DISABLE"

    log_transfer_detail {
      obs_period          = 2
      obs_period_unit     = "min"
      obs_bucket_name     = huaweicloud_obs_bucket.output.bucket
      obs_dir_prefix_name = "lts_transfer_obs_2_"
      obs_prefix_name     = "obs_2_"
      obs_time_zone       = "UTC-02:00"
      obs_time_zone_id    = "Etc/GMT+2"
    }
  }
}
`, testAccTransfer_base(name), name)
}

func TestAccTransfer_dis(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_transfer.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTransferResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTransfer_dis(name, "ENABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "DIS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "realTime"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.dis_id",
						"huaweicloud_dis_stream.test", "stream_id"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.dis_name",
						"huaweicloud_dis_stream.test", "stream_name"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
				),
			},
			{
				Config: testTransfer_dis(name, "DISABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "DIS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "realTime"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "DISABLE"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.dis_id",
						"huaweicloud_dis_stream.test", "stream_id"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.dis_name",
						"huaweicloud_dis_stream.test", "stream_name"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
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

func testTransfer_dis(name, status string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dis_stream" "test" {
  stream_name     = "%[2]s"
  partition_count = 1
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = "DIS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "RAW"
    log_transfer_status = "%[3]s"

    log_transfer_detail {
      dis_id   = huaweicloud_dis_stream.test.stream_id
      dis_name = huaweicloud_dis_stream.test.stream_name
    }
  }
}
`, testAccTransfer_base(name), name, status)
}

func TestAccTransfer_agency(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_lts_transfer.test"
	rNameObs := "huaweicloud_lts_transfer.test2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTransferResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPrecheckDomainName(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTransfer_agency(name, "ENABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "DIS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "realTime"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.dis_id",
						"huaweicloud_dis_stream.test", "stream_id"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.dis_name",
						"huaweicloud_dis_stream.test", "stream_name"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_agency_transfer.0.agency_domain_id",
						acceptance.HW_DOMAIN_ID),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_agency_transfer.0.agency_domain_name",
						acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_agency_transfer.0.agency_name",
						"huaweicloud_identity_agency.test", "name"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_agency_transfer.0.agency_project_id",
						acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),

					resource.TestCheckResourceAttrPair(rNameObs, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rNameObs, "log_streams.0.log_stream_id",
						"huaweicloud_lts_stream.stream_2", "id"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_agency_transfer.0.agency_domain_id",
						acceptance.HW_DOMAIN_ID),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_agency_transfer.0.agency_domain_name",
						acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttrPair(rNameObs, "log_transfer_info.0.log_agency_transfer.0.agency_name",
						"huaweicloud_identity_agency.test", "name"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_agency_transfer.0.agency_project_id",
						acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttrSet(rNameObs, "log_group_name"),
				),
			},
			{
				Config: testTransfer_agency(name, "DISABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "DISABLE"),
					resource.TestCheckResourceAttr(rNameObs, "log_transfer_info.0.log_transfer_status", "DISABLE"),
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

func testTransfer_agency(name, status string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream" "stream_2" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s_obs"
}

resource "huaweicloud_dis_stream" "test" {
  stream_name     = "%[2]s"
  partition_count = 1
}

resource "huaweicloud_obs_bucket" "output" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_identity_agency" "test" {
  name                  = "%[2]s"
  description           = "This is a test agency"
  delegated_domain_name = "%[4]s"

  domain_roles = [
    "LTS FullAccess",
  ]
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = "DIS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "RAW"
    log_transfer_status = "%[3]s"

    log_transfer_detail {
      dis_id   = huaweicloud_dis_stream.test.stream_id
      dis_name = huaweicloud_dis_stream.test.stream_name
    }

    log_agency_transfer {
      agency_domain_id   = "%[5]s"
      agency_domain_name = huaweicloud_identity_agency.test.delegated_domain_name
      agency_name        = huaweicloud_identity_agency.test.name
      agency_project_id  = "%[6]s"
    }
  }
}

resource "huaweicloud_lts_transfer" "test2" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.stream_2.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "%[3]s"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = huaweicloud_obs_bucket.output.bucket
      obs_dir_prefix_name = "lts_transfer_obs_"
      obs_prefix_name     = "obs_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }

    log_agency_transfer {
      agency_domain_id   = "%[5]s"
      agency_domain_name = huaweicloud_identity_agency.test.delegated_domain_name
      agency_name        = huaweicloud_identity_agency.test.name
      agency_project_id  = "%[6]s"
    }
  }
}
`, testAccTransfer_base(name), name, status, acceptance.HW_DOMAIN_NAME, acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID)
}

// Before running this test, please ensure the kafka instance is registered in the LTS service.
func TestAccTransfer_dms(t *testing.T) {
	var (
		transfer interface{}
		name     = acceptance.RandomAccResourceName()
		rName    = "huaweicloud_lts_transfer.test"
		rc       = acceptance.InitResourceCheck(rName, &transfer, getTransferResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsDmsTransfer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTransfer_dms_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_name", "huaweicloud_lts_stream.test", "stream_name"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "DMS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "realTime"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "JSON"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.kafka_id",
						acceptance.HW_LTS_REGISTERED_KAFKA_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.kafka_topic",
						"huaweicloud_dms_kafka_topic.test", "name"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.lts_tags.#", "2"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.stream_tags.0", "all"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.struct_fields.0", "all"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.invalid_field_value", "not_matched_field"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccTransfer_dms_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.lts_tags.#", "1"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.lts_tags.0", "all"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.stream_tags.#", "0"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.struct_fields.#", "0"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.invalid_field_value", ""),
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

func testAccTransfer_dms_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[2]s"
  name        = "%[3]s"
  partitions  = 3
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }

  log_transfer_info {
    log_transfer_type   = "DMS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "JSON"
    log_transfer_status = "DISABLE"

    log_transfer_detail {
      kafka_id            = "%[2]s"
      kafka_topic         = huaweicloud_dms_kafka_topic.test.name
      lts_tags            = ["hostName", "collectTime"]
      stream_tags         = ["all"]
      struct_fields       = ["all"]
      invalid_field_value = "not_matched_field"
    }
  }
}
`, testAccTransfer_base(name), acceptance.HW_LTS_REGISTERED_KAFKA_INSTANCE_ID, name)
}

func testAccTransfer_dms_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[2]s"
  name        = "%[3]s"
  partitions  = 3
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }

  log_transfer_info {
    log_transfer_type   = "DMS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "JSON"
    log_transfer_status = "DISABLE"

    log_transfer_detail {
      kafka_id    = "%[2]s"
      kafka_topic = huaweicloud_dms_kafka_topic.test.name
      lts_tags    = ["all"]
    }
  }
}
`, testAccTransfer_base(name), acceptance.HW_LTS_REGISTERED_KAFKA_INSTANCE_ID, name)
}
