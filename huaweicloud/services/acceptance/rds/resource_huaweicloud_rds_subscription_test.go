package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRdsSubscription(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/subscriptions"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	subscription := utils.PathSearch(fmt.Sprintf("subscriptions[?id=='%s']|[0]", state.Primary.ID), listRespBody, nil)
	if subscription == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return subscription, nil
}

func TestAccRdsSubscription_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_rds_subscription.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRdsSubscription,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsSubscriberInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsSubscription_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "subscription_database", "test_sub_db"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "push"),
					resource.TestCheckResourceAttr(rName, "initialize_at", "immediate"),
					resource.TestCheckResourceAttrPair(rName, "local_subscription.0.publication_id",
						"data.huaweicloud_rds_publications.test", "publications.0.id"),
					resource.TestCheckResourceAttrPair(rName, "local_subscription.0.publication_name",
						"data.huaweicloud_rds_publications.test", "publications.0.publication_name"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.job_schedule_type", "one_time"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.one_time_occurrence.0.active_start_date",
						"2026-05-06"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.one_time_occurrence.0.active_start_time",
						"07:20:30"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "is_cloud"),
				),
			},
			{
				Config: testAccRdsSubscription_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "subscription_database", "test_sub_db"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "push"),
					resource.TestCheckResourceAttr(rName, "initialize_at", "immediate"),
					resource.TestCheckResourceAttrPair(rName, "local_subscription.0.publication_id",
						"data.huaweicloud_rds_publications.test", "publications.0.id"),
					resource.TestCheckResourceAttrPair(rName, "local_subscription.0.publication_name",
						"data.huaweicloud_rds_publications.test", "publications.0.publication_name"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.frequency.0.freq_type", "monthly_week"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.frequency.0.freq_interval", "5"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.frequency.0.freq_interval_monthly",
						"weekday"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.frequency.0.freq_relative_interval_monthly",
						"second"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.daily_frequency.0.freq_subday_type",
						"once"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.daily_frequency.0.active_start_time",
						"12:00:00"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.duration.0.active_start_date", "2020-07-15"),
					resource.TestCheckResourceAttr(rName, "job_schedule.0.duration.0.active_end_date", "2030-10-20"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRdsSubscriptionImportStateIdFunc(rName),
				ImportStateVerifyIgnore: []string{
					"initialize_at",
				},
			},
		},
	})
}

func testAccRdsSubscription_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publications" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_subscription" "test" {
  instance_id           = "%[2]s"
  subscription_database = "test_sub_db"
  subscription_type     = "push"
  initialize_at         = "immediate"

  local_subscription {
    publication_id   = data.huaweicloud_rds_publications.test.publications[0].id
    publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
  }

  job_schedule{
    job_schedule_type = "one_time"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }
  }
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_SUBSCRIBER_INSTANCE_ID)
}

func testAccRdsSubscription_update() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publications" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_subscription" "test" {
  instance_id           = "%[2]s"
  subscription_database = "test_sub_db"
  subscription_type     = "push"
  initialize_at         = "immediate"

  local_subscription {
    publication_id   = data.huaweicloud_rds_publications.test.publications[0].id
    publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
  }

  job_schedule{
    job_schedule_type = "recurring"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }

    frequency {
      freq_type                      = "monthly_week"
      freq_interval                  = 5
      freq_interval_monthly          = "weekday"
      freq_relative_interval_monthly = "second"
    }

    daily_frequency {
      freq_subday_type     = "once"
      active_start_time    = "12:00:00"
    }

    duration {
      active_start_date = "2020-07-15"
      active_end_date   = "2030-10-20"
    }
  }
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_SUBSCRIBER_INSTANCE_ID)
}

func testAccRdsSubscriptionImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
