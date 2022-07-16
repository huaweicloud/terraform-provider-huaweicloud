package meeting

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/meeting/v1/conferences"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/meeting"
)

func getConferenceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	token, err := meeting.NewMeetingToken(conf, state)
	if err != nil {
		return nil, err
	}

	opt := conferences.GetOpts{
		ConferenceId: state.Primary.ID,
		UserId:       state.Primary.Attributes["user_id"],
		Limit:        500,
		Token:        token,
	}
	return conferences.Get(meeting.NewMeetingV1Client(conf), opt)
}

func TestAccConference_basic(t *testing.T) {
	var (
		conference   conferences.Conference
		resourceName = "huaweicloud_meeting_conference.test"
		baiscName    = fmt.Sprintf("Test Conference (%s)", acctest.RandString(5))
		updateName   = fmt.Sprintf("(New) Test Conference (%s)", acctest.RandString(5))
		startTime    = time.Unix(time.Now().Unix()+10*60, 0).Format("2006-01-02 15:04") // 10 minutes later
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&conference,
		getConferenceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
			acceptance.TestAccPreCheckMeetingRoom(t)
			acceptance.TestAccPreCheckParticipants(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConference_basic(baiscName, startTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "topic", baiscName),
					resource.TestCheckResourceAttr(resourceName, "meeting_room_id", acceptance.HW_MEETING_ROOM_ID),
					resource.TestCheckResourceAttr(resourceName, "duration", "15"),
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "language", "zh-CN"),
					resource.TestCheckResourceAttr(resourceName, "timezone_id", "56"),
					resource.TestCheckResourceAttr(resourceName, "participant_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "participant.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.account_id", "chair_demo"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.role", "1"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.email", acceptance.HW_CHAIR_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "participant.1.role", "0"),
					resource.TestCheckResourceAttr(resourceName, "participant.1.email", acceptance.HW_GUEST_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_calendar", "true"),
				),
			},
			{
				Config: testAccConference_update(updateName, startTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "topic", updateName),
					resource.TestCheckResourceAttr(resourceName, "meeting_room_id", acceptance.HW_MEETING_ROOM_ID),
					resource.TestCheckResourceAttr(resourceName, "duration", "30"),
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "language", "zh-CN"),
					resource.TestCheckResourceAttr(resourceName, "timezone_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "participant_number", "10"),
					resource.TestCheckResourceAttr(resourceName, "participant.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.account_id", "chair_update"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.role", "1"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.email", acceptance.HW_CHAIR_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_notify", "false"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_calendar", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccConferenceImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"duration",
					"corp_id",
					"start_time",
				},
			},
		},
	})
}

func TestAccConference_cyclical(t *testing.T) {
	var (
		day          = 24 * 60 * 60
		week         = 7 * day
		conference   conferences.Conference
		resourceName = "huaweicloud_meeting_conference.test"
		rName        = fmt.Sprintf("Test Conference (%s)", acctest.RandString(5))
		startDate    = time.Unix(time.Now().Unix(), 0).Format("2006-01-02")               // today
		endDate      = time.Unix(time.Now().Unix()+int64(3*week), 0).Format("2006-01-02") // four weeks later
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&conference,
		getConferenceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
			acceptance.TestAccPreCheckMeetingRoom(t)
			acceptance.TestAccPreCheckParticipants(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConference_cyclical(rName, startDate, endDate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "topic", rName),
					resource.TestCheckResourceAttr(resourceName, "meeting_room_id", acceptance.HW_MEETING_ROOM_ID),
					resource.TestCheckResourceAttr(resourceName, "duration", "15"),
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "language", "zh-CN"),
					resource.TestCheckResourceAttr(resourceName, "timezone_id", "56"),
					resource.TestCheckResourceAttr(resourceName, "participant_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "participant.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.account_id", "chair_demo"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.role", "1"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.email", acceptance.HW_CHAIR_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "participant.1.role", "0"),
					resource.TestCheckResourceAttr(resourceName, "participant.1.email", acceptance.HW_GUEST_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.cycle", "Week"),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.pre_remind", "1"),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.start_date", startDate),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.end_date", endDate),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.points.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.points.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "cycle_params.0.points.1", "5"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_calendar", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccConferenceImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"duration",
					"corp_id",
				},
			},
		},
	})
}

func TestAccConference_pwdAuth(t *testing.T) {
	var (
		conference   conferences.Conference
		resourceName = "huaweicloud_meeting_conference.test"
		rName        = fmt.Sprintf("Test Conference (%s)", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&conference,
		getConferenceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckPwdAuth(t)
			acceptance.TestAccPreCheckMeetingRoom(t)
			acceptance.TestAccPreCheckParticipants(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConference_pwdAuth(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "account_name", acceptance.HW_MEETING_ACCOUNT_NAME),
					resource.TestCheckResourceAttr(resourceName, "account_password", acceptance.HW_MEETING_ACCOUNT_PASSWORD),
					resource.TestCheckResourceAttr(resourceName, "topic", rName),
					resource.TestCheckResourceAttr(resourceName, "meeting_room_id", acceptance.HW_MEETING_ROOM_ID),
					resource.TestCheckResourceAttr(resourceName, "duration", "15"),
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "language", "zh-CN"),
					resource.TestCheckResourceAttr(resourceName, "timezone_id", "56"),
					resource.TestCheckResourceAttr(resourceName, "participant_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "participant.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.account_id", "chair_demo"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.role", "1"),
					resource.TestCheckResourceAttr(resourceName, "participant.0.email", acceptance.HW_CHAIR_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "participant.1.role", "0"),
					resource.TestCheckResourceAttr(resourceName, "participant.1.email", acceptance.HW_GUEST_EMAIL),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.is_send_calendar", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccConferenceImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"duration",
					"start_time",
				},
			},
		},
	})
}

func testAccConferenceImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var (
			accountName, password         string
			appId, appKey, corpId, userId string
			conferenceId                  string
		)
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_meeting_conference" {
				accountName = rs.Primary.Attributes["account_name"]
				password = rs.Primary.Attributes["account_password"]
				appId = rs.Primary.Attributes["app_id"]
				appKey = rs.Primary.Attributes["app_key"]
				corpId = rs.Primary.Attributes["corp_id"]
				userId = rs.Primary.Attributes["user_id"]
				conferenceId = rs.Primary.ID
			}
		}
		if conferenceId != "" && accountName != "" && password != "" {
			return fmt.Sprintf("%s/%s/%s", conferenceId, accountName, password), nil
		}
		if conferenceId != "" && appId != "" && appKey != "" {
			return fmt.Sprintf("%s/%s/%s/%s/%s", conferenceId, appId, appKey, corpId, userId), nil
		}
		return "", fmt.Errorf("resource not found: %s", conferenceId)
	}
}

func testAccConference_basic(rName, startTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_conference" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  topic              = "%[7]s"
  start_time         = "%[8]s"
  meeting_room_id    = "%[4]s"
  duration           = 15
  media_types        = ["Voice", "Video", "Data"]
  language           = "zh-CN"
  timezone_id        = 56
  participant_number = 5

  participant {
    account_id = "chair_demo"
    role       = 1
    email      = "%[5]s"
  }
  participant {
    role  = 0
    email = "%[6]s"
  }

  configuration {
    is_send_notify   = true
    is_send_calendar = true
  }
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID,
		acceptance.HW_MEETING_ROOM_ID, acceptance.HW_CHAIR_EMAIL, acceptance.HW_GUEST_EMAIL, rName, startTime)
}

func testAccConference_update(rName, startTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_conference" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  topic              = "%[6]s"
  start_time         = "%[7]s"
  meeting_room_id    = "%[4]s"
  duration           = 30
  media_types        = ["Voice", "Video", "Data"]
  language           = "zh-CN"
  timezone_id        = 3
  participant_number = 10

  participant {
    account_id = "chair_update"
    role       = 1
    email      = "%[5]s"
  }

  configuration {
    is_send_notify   = false
    is_send_calendar = false
  }
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID,
		acceptance.HW_MEETING_ROOM_ID, acceptance.HW_CHAIR_EMAIL, rName, startTime)
}

func testAccConference_cyclical(rName, startDate, endDate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_conference" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  topic              = "%[7]s"
  meeting_room_id    = "%[4]s"
  duration           = 15
  media_types        = ["Voice", "Video", "Data"]
  language           = "zh-CN"
  timezone_id        = 56
  participant_number = 5

  participant {
    account_id = "chair_demo"
    role       = 1
    email      = "%[5]s"
  }
  participant {
    role  = 0
    email = "%[6]s"
  }

  cycle_params {
    cycle      = "Week"
    pre_remind = 1
    start_date = "%[8]s"
    end_date   = "%[9]s"
    interval   = 2
    points     = [1, 5]
  }

  configuration {
    is_send_notify   = true
    is_send_calendar = true
  }
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID,
		acceptance.HW_MEETING_ROOM_ID, acceptance.HW_CHAIR_EMAIL, acceptance.HW_GUEST_EMAIL, rName, startDate, endDate)
}

// Start the conference now.
func testAccConference_pwdAuth(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_conference" "test" {
  account_name     = "%[1]s"
  account_password = "%[2]s"

  topic              = "%[6]s"
  meeting_room_id    = "%[3]s"
  duration           = 15
  media_types        = ["Voice", "Video", "Data"]
  language           = "zh-CN"
  timezone_id        = 56
  participant_number = 5

  participant {
    account_id = "chair_demo"
    role       = 1
    email      = "%[4]s"
  }
  participant {
    role  = 0
    email = "%[5]s"
  }

  configuration {
    is_send_notify   = true
    is_send_calendar = true
  }
}
`, acceptance.HW_MEETING_ACCOUNT_NAME, acceptance.HW_MEETING_ACCOUNT_PASSWORD, acceptance.HW_MEETING_ROOM_ID,
		acceptance.HW_CHAIR_EMAIL, acceptance.HW_GUEST_EMAIL, rName)
}
