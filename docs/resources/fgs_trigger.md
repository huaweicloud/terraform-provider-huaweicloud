---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_trigger"
description: ""
---

# huaweicloud_fgs_trigger

Manages a trigger resource within HuaweiCloud FunctionGraph.  
It's recommend to use `huaweicloud_fgs_function_trigger`, which makes a great improvement of managing function triggers.

## Example Usage

### Create a Timing Trigger with rate schedule type

```hcl
variable "function_urn" {}
variable "trigger_name" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "TIMER"

  timer {
    name          = var.trigger_name
    schedule_type = "Rate"
    schedule      = "1d"
  }
}
```

### Create a Timing Trigger with cron schedule type

```hcl
variable "function_urn" {}
variable "trigger_name" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "TIMER"

  timer {
    name          = var.trigger_name
    schedule_type = "Cron"
    schedule      = "@every 1h30m"
  }
}
```

### Create an OBS trigger

```hcl
variable "function_urn" {}
variable "bucket_name" {}
variable "trigger_name" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "OBS"
  status       = "ACTIVE"

  obs {
    bucket_name             = var.bucket_name
    event_notification_name = var.trigger_name
    suffix                  = ".json"

    events = ["ObjectCreated"]
  }
}
```

### Create an SMN trigger

```hcl
variable "function_urn" {}
variable "topic_urn" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "SMN"
  status       = "ACTIVE"

  smn {
    topic_urn = var.topic_urn
  }
}
```

### Create a DIS trigger

```hcl
variable "function_urn" {}
variable "stream_name" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "DIS"
  status       = "ACTIVE"

  dis {
    stream_name       = var.stream_name
    starting_position = "TRIM_HORIZON"
    max_fetch_bytes   = 2097152
    pull_period       = 30000
    serial_enable     = true
  }
}
```

### Create a DMS Kafka trigger

```hcl
variable "function_urn" {}
variable "kafka_instance_id" {}
variable "kafka_topic_id" {}
variable "user_password" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "KAFKA"

  kafka {
    instance_id = var.kafka_instance_id
    user_name   = "user"
    password    = var.user_password     
    batch_size  = 100

    topic_ids = [
      var.kafka_topic_id
    ]
  }
}
```

### Create a Dedicated APIG trigger

```hcl
variable "function_urn" {}
variable "instance_id" {}
variable "group_id" {}
variable "api_name" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "DEDICATEDGATEWAY"
  status       = "ACTIVE"

  apig {
    instance_id = var.instance_id
    group_id    = var.group_id
    api_name    = var.api_name
    env_name    = "RELEASE"
  }
}
```

### Create a Shared APIG trigger

```hcl
variable "function_urn" {}
variable "group_id" {}
variable "api_name" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "APIG"
  status       = "ACTIVE"

  apig {
    group_id = var.group_id
    api_name = var.api_name
    env_name = "RELEASE"
  }
}
```

### Create a LTS trigger

```hcl
variable "log_group_id" {}
variable "log_topic_id" {}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = var.function_urn
  type         = "LTS"
  status       = "ACTIVE"

  lts {
    log_group_id = var.log_group_id
    log_topic_id = var.log_topic_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the trigger resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new trigger resource.

* `function_urn` - (Required, String, ForceNew) Specifies the Uniform Resource Name (URN) of the function.
  Changing this will create a new trigger resource.

* `type` - (Required, String, ForceNew) Specifies the type of the function.
  The valid values currently only support **TIMER**, **OBS**, **SMN**, **DIS**, **KAFKA**, **APIG**, **LTS**, and
  **DEDICATEDGATEWAY**. Changing this will create a new trigger resource.

* `status` - (Optional, String) Specifies whether trigger is enabled. The valid values are **ACTIVE** and **DISABLED**.
  About DMS kafka trigger, the default value is **ACTIVE**.

  -> **NOTE:** Currently, SMN triggers do not support `status`, and OBS triggers do not support updating `status`.

* `timer` - (Optional, List, ForceNew) Specifies the configuration of the timing trigger.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_timer) structure is documented below.

* `obs` - (Optional, List, ForceNew) Specifies the configuration of the OBS trigger.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_obs) structure is documented below.

* `smn` - (Optional, List, ForceNew) Specifies the configuration of the SMN trigger.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_smn) structure is documented below.

* `dis` - (Optional, List, ForceNew) Specifies the configuration of the DIS trigger.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_dis) structure is documented below.

  -> **NOTE:** Specify an agency with DIS access permissions for the function version before you can create a DIS
  trigger.

* `kafka` - (Optional, List, ForceNew) Specifies the configuration of the DMS trigger for Kafka.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_kafka) structure is documented below.

  -> **NOTE:** VPC access must be enabled for the function before you create a Kafka trigger.
  The port `9092` must be opened for security group ingress rules.

* `apig` - (Optional, List, ForceNew) Specifies the configuration of the shared APIG and dedicated APIG trigger.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_apig) structure is documented below.

* `lts` - (Optional, List, ForceNew) Specifies the configuration of the LTS trigger.
  Changing this will create a new trigger resource.
  The [object](#fgs_trigger_lts) structure is documented below.

<a name="fgs_trigger_timer"></a>
The `timer` block supports:

* `name` - (Required, String, ForceNew) Specifies the trigger name, which can contains of `1` to `64` characters.
  The name must start with a letter, only letters, digits, hyphens (-) and underscores (_) are allowed.
  Changing this will create a new trigger resource.

* `schedule_type` - (Required, String, ForceNew) Specifies the type of the time schedule.
  The valid values are **Rate** and **Cron**.
  Changing this will create a new trigger resource.

* `schedule` - (Required, String, ForceNew) Specifies the time schedule.
  For the rate type, schedule is composed of time and time unit.
  The time unit supports minutes (m), hours (h) and days (d).
  For the corn expression, please refer to the HuaweiCloud
  [document](https://support.huaweicloud.com/en-us/usermanual-functiongraph/functiongraph_01_0908.html).
  Changing this will create a new trigger resource.

* `additional_information` - (Optional, String, ForceNew) Specifies the event used by the timer to trigger the function.
  Changing this will create a new trigger resource.

<a name="fgs_trigger_obs"></a>
The `obs` block supports:

* `bucket_name` - (Required, String, ForceNew) Specifies the OBS bucket name.
  Changing this will create a new trigger resource.

* `events` - (Required, List, ForceNew) Specifies the events that can trigger functions.
  Changing this will create a new trigger resource.
  The valid values are as follows:
  + **ObjectCreated**, **Put**, **Post**, **Copy** and **CompleteMultipartUpload**.
  + **ObjectRemoved**, **Delete** and **DeleteMarkerCreated**.

  -> **NOTE:** If **ObjectCreated** is configured, **Put**, **Post**, **Copy** and **CompleteMultipartUpload** cannot
  be configured. If **ObjectRemoved** is configured, **Delete** and **DeleteMarkerCreated** cannot be configured.

* `event_notification_name` - (Required, String, ForceNew) Specifies the event notification name.
  Changing this will create a new trigger resource.

* `prefix` - (Optional, String, ForceNew) Specifies the prefix to limit notifications to objects beginning with this keyword.
  Changing this will create a new trigger resource.

* `suffix` - (Optional, String, ForceNew) Specifies the suffix to limit notifications to objects ending with this keyword.
  Changing this will create a new trigger resource.

<a name="fgs_trigger_smn"></a>
The `smn` block supports:

* `topic_urn` - (Required, String, ForceNew) Specifies the Uniform Resource Name (URN) for SMN topic.
  Changing this will create a new trigger resource.

<a name="fgs_trigger_dis"></a>
The `dis` block supports:

* `stream_name` - (Required, String, ForceNew) Specifies the name of the DIS stream resource.
  Changing this will create a new trigger resource.

* `starting_position` - (Required, String, ForceNew) Specifies the type of starting position for DIS queue.
  The valid values are as follows:
  + **TRIM_HORIZON**: Starts reading from the earliest data stored in the partitions.
  + **LATEST**: Starts reading from the latest data stored in the partitions.
  Changing this will create a new trigger resource.

* `max_fetch_bytes` - (Required, Int, ForceNew) Specifies the maximum volume of data that can be obtained for a single
  request, in Byte. Only the records with a size smaller than this value can be obtained.
  The valid value is range from `1,024` to `4,194,304`.
  Changing this will create a new trigger resource.

* `pull_period` - (Required, Int, ForceNew) Specifies the interval at which data is pulled from the specified stream.
  The valid value is range from `2` to `60,000`.
  Changing this will create a new trigger resource.

* `serial_enable` - (Required, Bool, ForceNew) Specifies the determines whether to pull data only after the data pulled
  in the last period has been processed.
  Changing this will create a new trigger resource.

<a name="fgs_trigger_kafka"></a>
The `kafka` block supports:

* `instance_id` - (Required, String, ForceNew) Specifies the DMS instance ID for kafka.
  Changing this will create a new trigger resource.

* `user_name` - (Optional, String, ForceNew) Specifies the username for logging in to the Kafka Manager.
  Changing this will create a new trigger resource.

* `password` - (Optional, String, ForceNew) Specifies the password for logging in to the Kafka Manager.
  Changing this will create a new trigger resource.

* `topic_ids` - (Required, List, ForceNew) Specifies one or more topic IDs of DMS kafka instance.
  Changing this will create a new trigger resource.

* `batch_size` - (Optional, Int, ForceNew) Specifies the The number of messages consumed from the topic each time.
  The valid value is range from `1` to `1,000`. Defaults to `100`.
  Changing this will create a new trigger resource.

<a name="fgs_trigger_apig"></a>
The `apig` block supports:

* `group_id` - (Required, String, ForceNew) Specifies the ID of the APIG group to which the API belongs.
  Changing this will create a new trigger resource.

* `env_name` - (Required, String, ForceNew) Specifies the API environment name.
  Changing this will create a new trigger resource.

* `api_name` - (Required, String, ForceNew) Specifies the API name. Changing this will create a new trigger resource.

* `instance_id` - (Optional, String, ForceNew) Specifies the ID of the APIG dedicated instance to which the API belongs.
  Required if the `type` is `DEDICATEDGATEWAY`. Changing this will create a new trigger resource.

* `security_authentication` - (Optional, String, ForceNew) Specifies the security authentication mode. The valid values
  are **NONE**, **APP** and **IAM**, default to **IAM**. Changing this will create a new trigger resource.

* `request_protocol` - (Optional, String, ForceNew) Specifies the request protocol of the API. The valid value are
  **HTTP** and **HTTPS**. Default to **HTTPS**. Changing this will create a new trigger resource.

* `timeout` - (Optional, Int, ForceNew) Specifies the timeout for request sending. The valid value is range form
  `1` to `60,000`, default to `5,000`. Changing this will create a new trigger resource.

<a name="fgs_trigger_lts"></a>
The `lts` block supports:

* `log_group_id` - (Required, String, ForceNew) Specifies the log group ID.
  Changing this will create a new trigger resource.

* `log_topic_id` - (Required, String, ForceNew) Specifies the log stream ID.
  Changing this will create a new trigger resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 2 minutes.
