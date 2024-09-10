---
subcategory: "Meeting"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_meeting_conference"
description: ""
---

# huaweicloud_meeting_conference

Using this resource to book a conference within HuaweiCloud.

~> Using duplicate conference IDs (in the Console) will make terraform resource state (behavior) inconsistent.

## Example Usage

### Book an Ordinary Conference using third-party application authorization

```hcl
variable "app_id" {}
variable "app_key" {}
variable "conference_topic" {}
variable "meeting_room_id" {}
variable "account_id" {}

resource "huaweicloud_meeting_conference" "test" {
  app_id  = var.app_id
  app_key = var.app_key

  topic           = var.conference_topic
  meeting_room_id = var.meeting_room_id
  start_time      = formatdate("YYYY-MM-DD hh:mm", timeadd(timestamp(), "1h")) // Start in an hour
  duration        = 15
  media_types     = ["Voice", "Video", "Data"]

  language           = "zh-CN"
  timezone_id        = 56
  participant_number = 8

  participants {
    account_id = var.account_id
    role       = 1
    phone      = "+99123456789876432"
    email      = "email@example.com"
  }

  configuration {
    is_send_notify   = true
    is_send_calendar = true
  }
}
```

### Book an Ordinary Conference using account authorization

```hcl
variable "account_name" {}
variable "account_password" {}
variable "conference_topic" {}
variable "meeting_room_id" {}
variable "account_id" {}

resource "huaweicloud_meeting_conference" "test" {
  account_name     = var.account_name
  account_password = var.account_password

  topic           = var.conference_topic
  meeting_room_id = var.meeting_room_id
  start_time      = formatdate("YYYY-MM-DD hh:mm", timeadd(timestamp(), "1h")) // Start in an hour
  duration        = 15
  media_types     = ["Voice", "Video", "Data"]

  language           = "zh-CN"
  timezone_id        = 56
  participant_number = 8

  participants {
    account_id = var.account_id
    role       = 1
    phone      = "+99123456789876432"
    email      = "email@example.com"
  }

  configuration {
    is_send_notify   = true
    is_send_calendar = true
  }
}
```

### Book a Cyclical Conference

```hcl
variable "app_id" {}
variable "app_key" {}
variable "conference_topic" {}
variable "meeting_room_id" {}
variable "account_id" {}

resource "huaweicloud_meeting_conference" "test" {
  app_id  = var.app_id
  app_key = var.app_key

  topic           = var.conference_topic
  meeting_room_id = var.meeting_room_id
  duration        = 15
  media_types     = ["Voice", "Video", "Data"]

  language           = "zh-CN"
  timezone_id        = 56
  participant_number = 8

  participants {
    account_id = var.account_id
    role       = 1
    phone      = "+99123456789876432"
    email      = "email@example.com"
  }

  cycle_params {
    cycle      = "Week"
    pre_remind = 1
    start_date = formatdate("YYYY-MM-DD", timestamp())
    end_date   = formatdate("YYYY-MM-DD", timeadd(timestamp(), "336h")) // The meeting period is a week.
    interval   = 2
    points     = [1, 5]
  }

  configuration {
    is_send_notify   = true
    is_send_calendar = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Optional, String, ForceNew) Specifies the (HUAWEI Cloud meeting) user account name to which the
  meeting initiator belongs. Changing this parameter will create a new resource.

* `account_password` - (Optional, String, ForceNew) Specifies the user password.
  Required if `account_name` is set. Changing this parameter will create a new resource.

* `app_id` - (Optional, String, ForceNew) Specifies the ID of the Third-party application.
  Changing this parameter will create a new resource.

  -> You can apply for an application and obtain the App ID and App Key in the console.

* `app_key` - (Optional, String, ForceNew) Specifies the Key information of the Third-party APP.
  Required if `app_id` is set. Changing this parameter will create a new resource.

* `corp_id` - (Optional, String, ForceNew) Specifies the corporation ID.
  Required if the application is used in multiple enterprises. Only available if `app_id` is set.
  Changing this parameter will create a new resource.

* `user_id` - (Optional, String, ForceNew) Specifies the user ID of the meeting initiator.
  Only available if `app_id` is set. If omitted, the user ID of default administrator will be used.
  Changing this parameter will create a new resource.

-> Exactly one of account authorization and application authorization you must select.

* `topic` - (Required, String) Specifies the conference topic. The topic can contain `1` to `128` characters.

* `meeting_room_id` - (Required, String) Specifies the cloud meeting room ID.

* `duration` - (Required, Int) Specifies the duration of the conference, in minutes.
  The valid value is range from `15` to `1,440`, defaults to `30`.

  -> After the conference starts, only support extend duration, shorten duration is not supported.
    And only the duration can be updated after the meeting starts.

* `start_time` - (Optional, String) Specifies the conference start time (UTC time).
  The time format is `YYYY-MM-DD hh:mm`, e.g. `2006-01-02 15:04`.
  There is no need to set if you book a cyclical conference.

  -> If you want to start a conference at `08:00` (UTC+8), you need to specify the time with `00:00`.
  And the start time cannot be earlier than now.

* `media_types` - (Optional, List) Specifies the conference media type list.
  It consists of one or more enumerations, and the valid values are as follows:
  + **Voice**: Voice.
  + **Video**: SD video.
  + **HDVideo**: High-definition video (mutually exclusive with Video, if Video and HDVideo are selected at the same
    time, the system will select Video by default).
  + **Data**: Multimedia (If omitted, the system configuration will determines whether to automatically add **Data**).

* `is_auto_record` - (Optional, Int) Specifies whether the conference automatically starts recording, it only takes
  effect when the recording type is:
  + **1**: Automatically start recording.
  + **0**: Do not start recording automatically.

  The default value is `0` (not to start automatically).

* `encrypt_mode` - (Optional, Int) Specifies the conference media encryption mode.
  + **0**: Adaptive encryption.
  + **1**: Force encryption.
  + **2**: Do not encrypt.

  The default value is populated by enterprise-level configuration.

* `language` - (Optional, String) Specifies the default language of the conference, the default value is defined by the
  conference cloud service. For languages supported by the system, it is passed according to the RFC3066 specification.
  The valid values are as follows:
  + **zh-CN**: Simplified Chinese.
  + **en-US**: US English.

* `timezone_id` - (Optional, Int) Specifies the time zone information of the conference time in the conference
  notification. For time zone information, refer to the [time zone mapping relationship](#time_zone_mapping).

  -> For example: the timeZoneID `26`, time in the conference notification sent through HUAWEI CLOUD conference will be
  marked as "2021/11/11 Thursday 00:00 - 02:00 (GMT) Greenwich Standard When: Dublin, Edinburgh, Lisbon, London".
  For an aperiodic conference, if the conference notification is sent through a third-party system, this field does not
  need to be filled in.

* `record_type` - (Optional, Int) Specifies the recording type.
  + **0**: Disabled.
  + **1**: Live broadcast.
  + **2**: Record and broadcast.
  + **3**: Live + Recording.

  The default value is `0` (disabled).

* `live_address` - (Optional, String) Specifies the mainstream live broadcast address, with a maximum of 255 characters.
  Only available if `record_type` is `2` or `3`.

* `aux_address` - (Optional, String) Specifies the auxiliary streaming address, the maximum length is 255 characters.
  Only available if `record_type` is `2` or `3`.

* `is_record_aux_stream` - (Optional, Int) Specifies whether to record auxiliary stream.
  + **0**: Do not record.
  + **1**: Record.

  Only available if `record_type` is `2` or `3`, and the default value is `0`.

* `record_auth_type` - (Optional, Int) Specifies the recording authentication method.
  + **0**: Viewable/downloadable via link.
  + **1**: Enterprise users can watch/download.
  + **2**: Attendees can watch/download.

  Only available if `record_type` is `2` or `3`.

* `participant_number` - (Optional, Int) Specifies the number of parties in the conference, the maximum number of
  participants in the conference. Defaults to `0` (Unlimited).

* `participant` - (Optional, List) Specifies the attendee list.
  The [object](#conference_participant) structure is documented below.

* `cycle_params` - (Optional, List) Specifies the configurations of the cyclical conference.
  The [object](#conference_cycle_params) structure is documented below.

* `configuration` - (Optional, List) Specifies the other conference configurations.
  The [object](#conference_configuration) structure is documented below.

<a name="conference_participant"></a>
The `participant` block supports:

* `user_id` - (Optional, String) Specifies the user ID of the participant.

* `account_id` - (Optional, String) Specifies the account ID of the participant.

* `name` - (Optional, String) Specifies the attendee name or nickname.  
  The valid length is limited from `1` to `96`.

* `role` - (Optional, Int) Specifies the role in the conference. The valid values are as follows:
  + **0**: Normal attendee.
  + **1**: The conference chair.

* `type` - (Optional, String) Specifies the call-in type. The valid values are as follows:
  + **normal**: The soft terminal.
  + **terminal**: The conference room or hard terminal.
  + **outside**: The outside participant.
  + **mobile**: The user's landline phone.
  + **ideahub**: The ideahub.

* `is_mute` - (Optional, Int) Specifies whether the user needs to be automatically muted when joining the conference
  (only effective when invited in the conference). The valid values are as follows:
  + **0**: No mute.
  + **1**: Mute.

  The default value is `0`.

* `is_auto_invite` - (Optional, Int) Specifies whether to automatically invite this participant when the conference
  starts. The valid values are as follows:
  + **0**: Do not automatically invite.
  + **1**: Automatic invitation.

  The default value is populated by enterprise-level configuration.

* `phone` - (Optional, String) Specifies the SIP or TEL number, maximum of 127 characters.

* `email` - (Optional, String) Specifies the email address.

* `sms` - (Optional, String) Specifies the mobile number for SMS notification, maximum of 32 characters.

  -> At least one of `phone`, `email` and `sms` must be set.

<a name="conference_cycle_params"></a>
The `cycle_params` block supports:

* `cycle` - (Required, String) Specifies the period type. The valid values are as follows:
  + **Day**
  + **Week**
  + **Month**

* `pre_remind` - (Required, Int) Specifies the number of days for advance conference notice.
  The valid value is range from `0` to `30`, defaults to `1`.

* `start_date` - (Required, String) Specifies the start date of the recurring conference.
  The format is `YYYY-MM-DD`.

* `end_date` - (Required, String) Specifies the end date of the recurring conference.
  The format is `YYYY-MM-DD`.

* `interval` - (Optional, Int) Specifies the cycle interval.
  For different `cycle` types, the value range of interval are as follows:
  + **Day**: Means that it will be held every few days, and the valid value is range from `1` to `15`.
  + **Week**: Means that it is held every few weeks, and the valid value is range from `1` to `5`.
  + **Month**: Means every few months, the value range is `1` to `3`.

* `points` - (Optional, List) Specifies the conference point in the cycle. Only valid by **Week** and **Month**.
  For different `cycle` types, the value range of elements are as follows:
  + **Week**: The valid value is range from `0` to `6`. The `0` means Sunday, `6` means Saturday.
  + **Month**: The valid range for the elements is `1` to `31`. If the value does not exist in the current month, the
  value means the end of the month.

<a name="conference_configuration"></a>
The `configuration` block supports:

* `is_send_notify` - (Optional, Bool) Specifies whether to send conference email notification.

* `is_send_sms` - (Optional, Bool) Specifies whether to send conference SMS notification.

* `is_send_calendar` - (Optional, Bool) Specifies whether to send conference calendar notifications.

* `is_auto_mute` - (Optional, Bool) Specifies whether the soft terminal is automatically muted when the guest joins the
  conference.

* `is_hard_terminal_auto_mute` - (Optional, Bool) Specifies whether the guest joins the conference, whether the hard
  terminal is automatically muted.

* `is_guest_free_password` - (Optional, Bool) Specifies whether the guest is password-free (only valid for random
  conferences).

* `callin_restriction` - (Optional, Int) Specifies the range to allow incoming calls.
  + **0**: All users.
  + **2**: Users within the enterprise.
  + **3**: The invited user.

* `allow_guest_start` - (Optional, Bool) Specifies whether to allow guests to start conferences (only valid for random
  ID conferences).

* `guest_password` - (Optional, String) Specifies the guest password (pure number which is `4` to `16` digits long).

* `prolong_time` - (Optional, Int) Specifies the Automatically extend duration, the valid value is range from `0` to
  `60`.

* `waiting_room_enabled` - (Optional, Bool) Specifies whether to open the waiting room (only valid for RTC enterprises).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The conference ID.

* `conference_uuid` - The conference UUID.

* `conference_type` - The conference type, the valid values are as follows:
  + **FUTURE**
  + **IMMEDIATELY**
  + **CYCLE**

* `access_number` - The access number of the conference.

* `status` - The conference status, the valid values are as follows:
  + **Schedule**: the conference is in schedule.
  + **Created**: The conference is in progress.

* `chair_join_uri` - The host meeting link address.

* `guest_join_uri` - The common attendee meeting link address.

* `audience_join_uri` - The audience meeting link address.

* `subconferences` - The list of periodic sub-conferences.
  The [object](#conference_subconferences) structure is documented below.

* `join_password` - The meeting password.
  The [join_password](#join_password) structure is documented below.

<a name="conference_subconferences"></a>
The `subconferences` block supports:

* `id` - The sub-conference ID.

* `media_types` - The media type list.

* `start_time` - The sub-conference start time.

* `end_time` - The sub-conference end time.

* `is_auto_record` - Whether to automatically start recording.

* `record_auth_type` - The recording authentication method.

* `subconfiguration` - The other configuration information of periodic subconferences.
  The [object](#conference_subconfiguration) structure is documented below.

<a name="join_password"></a>
The `join_password` block supports:

* `host` - The password of the meeting host.

* `guest` - The password of the common participant.

<a name="conference_subconfiguration"></a>
The `subconfiguration` block supports:

* `callin_restriction` - The range to allow incoming calls.

* `audience_callin_restriction` - The range that the webinar audience is allowed to call in.
  The valid values are as follows:
  + **0**: All users.
  + **2**: Users within the enterprise.

* `allow_guest_start` - Whether to allow guests to start conferences (only valid for random ID conferences).

* `waiting_room_enabled` - Whether to open the waiting room (only valid for RTC enterprises).

* `show_audience_policies` - The webinar Audience Display Strategy.
  The [object](#conference_show_audience_policy) structure is documented below.

<a name="conference_show_audience_policy"></a>
The `show_audience_policy` block supports:

* `mode` - Audience display strategy: The server is used to calculate the number of audiences and send it to the client
  to control the audience display.
  + **0**: Do not display.
  + **1**: Multiply display the number of participants, based on the real-time number of participants or the cumulative
  number of participants, the multiplication setting can be performed.

* `multiple` - Specifies the multiplier. The valid values is range from `0` to `10`, it can be set to 1 decimal place.

* `base_audience_count` - Specifies the basic number of people, the valid values is range from `0` to `10,000`.

-> After `base_audience_count` and `multiple` are both setting, the number of people displayed is: **NX+Y**.
  When NX calculates a non-integer, it will be rounded down.

## Import

Conferences (only scheduled conference and progressing conference) can be imported using their `id` and authorization
parameters, separated by slashes, e.g.

Import a conference and authenticated by account.

```bash
$ terraform import huaweicloud_meeting_conference.test <id>/<account_name>/<account_password>
```

Import a conference and authenticated by `APP ID`/`APP Key`.

```bash
$ terraform import huaweicloud_meeting_conference.test <id>/<app_id>/<app_key>/<corp_id>/<user_id>
```

The slashes cannot be missing even corporation ID and user ID are empty.

Note that importing is not supported for expired conferences and the start time of the meeting is not imported along
with it. You can ignore this change as below.

```hcl
resource "huaweicloud_meeting_conference" "test" {
    ...

  lifecycle {
    ignore_changes = [
      start_time,
    ]
  }
}
```

## Appendix

<a name="time_zone_mapping"></a>
The time zone mapping relationship supports:

| Timezone ID | Timezone |
| ---- | ---- |
| 1 | (GMT-12:00) Eniwetok, Kwajalein |
| 2 | (GMT-11:00) Midway Island, Samoa |
| 3 | (GMT-10:00) Hawaii |
| 4 | (GMT-09:00) Alaska |
| 5 | (GMT-08:00) Pacific Time(US&Canada);Tijuana |
| 6 | (GMT-07:00) Arizona |
| 7 | (GMT-07:00) Mountain Time(US&Canada) |
| 8 | (GMT-06:00) Central America |
| 9 | (GMT-06:00) Central Time(US&Canada) |
| 10 | (GMT-06:00) Mexico City |
| 11 | (GMT-06:00) Saskatchewan |
| 12 | (GMT-05:00) Bogota, Lima, Quito |
| 13 | (GMT-05:00) Eastern Time(US&Canada) |
| 14 | (GMT-05:00) Indiana(East) |
| 15 | (GMT-04:00) Atlantic time(Canada) |
| 16 | (GMT-04:00) Caracas, La Paz |
| 17 | (GMT-04:00) Santiago |
| 18 | (GMT-03:30) Newfoundland |
| 19 | (GMT-03:00) Brasilia |
| 20 | (GMT-03:00) Buenos Aires, Georgetown |
| 21 | (GMT-03:00) Greenland |
| 22 | (GMT-02:00) Mid-Atlantic |
| 23 | (GMT-01:00) Azores |
| 24 | (GMT-01:00) Cape Verde Is. |
| 25 | (GMT) Casablanca, Monrovia |
| 26 | (GMT) Greenwich Mean Time:Dublin, Edinburgh, Lisbon, London |
| 27 | (GMT+01:00) Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna |
| 28 | (GMT+01:00) Belgrade, Bratislava, Budapest, Ljubljana, Prague |
| 29 | (GMT+01:00) Brussels, Copenhagen, Madrid, Paris |
| 30 | (GMT+01:00) Sarajevo, Skopje, Sofija, Warsaw, Zagreb |
| 31 | (GMT+01:00) West Central Africa |
| 32 | (GMT+02:00) Athens, Istanbul, Vilnius |
| 33 | (GMT+02:00) Bucharest |
| 34 | (GMT+02:00) Cairo |
| 35 | (GMT+02:00) Harare, Pretoria |
| 36 | (GMT+02:00) Helsinki, Riga, Tallinn |
| 37 | (GMT+02:00) Jerusalem |
| 38 | (GMT+03:00) Baghdad, Minsk |
| 39 | (GMT+03:00) Kuwait, Riyadh |
| 40 | (GMT+03:00) Moscow, St. Petersburg, Volgograd |
| 41 | (GMT+03:00) Nairobi |
| 42 | (GMT+03:30) Tehran |
| 43 | (GMT+04:00) Abu Dhabi, Muscat |
| 44 | (GMT+04:00) Baku, Tbilisi, Yerevan |
| 45 | (GMT+04:30) Kabul |
| 46 | (GMT+05:00) Ekaterinburg |
| 47 | (GMT+05:00) Islamabad, Karachi, Tashkent |
| 48 | (GMT+05:30) Calcutta, Chennai, Mumbai, New Delhi |
| 49 | (GMT+05:45) Kathmandu |
| 50 | (GMT+06:00) Almaty, Novosibirsk |
| 51 | (GMT+06:00) Astana, Dhaka |
| 52 | (GMT+06:00) Sri Jayawardenepura |
| 53 | (GMT+06:30) Rangoon |
| 54 | (GMT+07:00) Bangkok, Hanoi, Jakarta |
| 55 | (GMT+07:00) Krasnoyarsk |
| 56 | (GMT+08:00) Beijing, Chongqing, Hong Kong, Urumqi, Taipei |
| 57 | (GMT+08:00) Irkutsk, Ulaan Bataar |
| 58 | (GMT+08:00) Kuala Lumpur, Singapore |
| 59 | (GMT+08:00) Perth |
| 60 | (GMT+09:00) Osaka, Sapporo, Tokyo |
| 61 | (GMT+09:00) Seoul |
| 62 | (GMT+09:00) Yakutsk |
| 63 | (GMT+09:30) Adelaide |
| 64 | (GMT+09:30) Darwin |
| 65 | (GMT+10:00) Brisbane |
| 66 | (GMT+10:00) Canberra, Melbourne, Sydney |
| 67 | (GMT+10:00) Guam, Port Moresby |
| 68 | (GMT+10:00) Hobart |
| 69 | (GMT+10:00) Vladivostok |
| 70 | (GMT+11:00) Magadan, Solomon Is., New Caledonia |
| 71 | (GMT+12:00) Auckland, Welington |
| 72 | (GMT+12:00) Fiji |
| 73 | (GMT+13:00) Nuku'alofa |
| 74 | (GMT+09:00) Irkutsk |
| 75 | (GMT) Casablanca |
| 76 | (GMT+04:00) Baku |
| 77 | (GMT+12:00) Kamchatka, Marshall Is. |
