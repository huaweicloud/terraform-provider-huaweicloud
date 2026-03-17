package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getScheduleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region   = acceptance.HW_REGION_NAME
		product  = "cfw"
		objectId = state.Primary.Attributes["object_id"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	return cfw.GetScheduleById(client, objectId, state.Primary.ID)
}

func TestAccSchedule_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cfw_schedule.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getScheduleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSchedule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "periodic.#", "3"),
					resource.TestCheckResourceAttr(rName, "periodic.0.type", "0"),
					resource.TestCheckResourceAttr(rName, "periodic.0.start_time", "00:00:00"),
					resource.TestCheckResourceAttr(rName, "periodic.0.end_time", "23:59:59"),
					resource.TestCheckResourceAttr(rName, "periodic.1.type", "1"),
					resource.TestCheckResourceAttr(rName, "periodic.1.start_time", "00:00:00"),
					resource.TestCheckResourceAttr(rName, "periodic.1.end_time", "23:59:59"),
					resource.TestCheckResourceAttr(rName, "periodic.1.week_mask.#", "1"),
					resource.TestCheckResourceAttr(rName, "periodic.2.type", "2"),
					resource.TestCheckResourceAttr(rName, "periodic.2.start_time", "00:00:00"),
					resource.TestCheckResourceAttr(rName, "periodic.2.end_time", "23:59:59"),
					resource.TestCheckResourceAttr(rName, "periodic.2.start_week", "1"),
					resource.TestCheckResourceAttr(rName, "periodic.2.end_week", "1"),
					resource.TestCheckResourceAttr(rName, "absolute.#", "1"),
					resource.TestCheckResourceAttr(rName, "absolute.0.start_time", "1773730600349"),
					resource.TestCheckResourceAttr(rName, "absolute.0.end_time", "1774076220000"),
					resource.TestCheckResourceAttrSet(rName, "ref_count"),
				),
			},
			{
				Config: testSchedule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "periodic.#", "1"),
					resource.TestCheckResourceAttr(rName, "periodic.0.type", "0"),
					resource.TestCheckResourceAttr(rName, "periodic.0.start_time", "01:00:00"),
					resource.TestCheckResourceAttr(rName, "periodic.0.end_time", "23:55:59"),
					resource.TestCheckResourceAttr(rName, "absolute.#", "1"),
					resource.TestCheckResourceAttr(rName, "absolute.0.start_time", "1774730600349"),
					resource.TestCheckResourceAttr(rName, "absolute.0.end_time", "1775076220000"),
					resource.TestCheckResourceAttrSet(rName, "ref_count"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testScheduleImportState(rName),
			},
		},
	})
}

func testSchedule_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_schedule" "test" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%s"
  description = "test description"

  absolute {
    end_time   = 1774076220000
    start_time = 1773730600349
  }

  periodic {
    type       = 0
    start_time = "00:00:00"
    end_time   = "23:59:59"
  }

  periodic {
    type       = 1
    start_time = "00:00:00"
    end_time   = "23:59:59"
    week_mask = [
      1,
    ]
  }

  periodic {
    type       = 2
    start_time = "00:00:00"
    end_time   = "23:59:59"
    start_week = 1
    end_week   = 1
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testSchedule_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_schedule" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name      = "%s-update"

  absolute {
    end_time   = 1775076220000
    start_time = 1774730600349
  }

  periodic {
    type       = 0
    start_time = "01:00:00"
    end_time   = "23:55:59"
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testScheduleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		objectId := rs.Primary.Attributes["object_id"]
		if objectId == "" {
			return "", fmt.Errorf("attribute (object_id) of Resource (%s) not found", name)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", objectId, rs.Primary.ID), nil
	}
}
