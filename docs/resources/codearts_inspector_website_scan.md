---
subcategory: "CodeArts Inspector"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_inspector_website_scan"
description: ""
---

# huaweicloud_codearts_inspector_website_scan

Manages a CodeArts inspector website scan resource within HuaweiCloud.

## Example Usage

### With normal task type

```hcl
variable "url" {}
variable "timer" {}

resource "huaweicloud_codearts_inspector_website_scan" "test" {
  task_name = "normal-name"
  task_type = "normal"
  url       = var.url
  timer     = var.timer
  scan_mode = "deep"
  port_scan = true
}
```

### With monitor task type

```hcl
variable "url" {}
variable "trigger_time" {}

resource "huaweicloud_codearts_inspector_website_scan" "test" {
  task_name    = "monitor-name"
  task_type    = "monitor"
  url          = var.url
  trigger_time = var.trigger_time
  task_period  = "everyweek"
  scan_mode    = "normal"
  port_scan    = true
}
```

## Argument Reference

The following arguments are supported:

* `task_name` - (Required, String, ForceNew) Specifies the task name. Changing this parameter will create a new resource.
  The valid length is limited from `1` to `24`. Only Chinese characters, letters, digits, hyphens (-) and underscores (_)
  are allowed, and cannot start with a hyphen (-).

* `url` - (Required, String, ForceNew) Specifies the destination URL to scan. Changing this parameter will create a new
  resource. The maximum length is `256`. The format should be `http(s)://example.com` or
  `http(s)://{public IPv4 address}:{PORT}`. The value can only be the website address of CodeArts inspector website,
  and the website address must be authorized before it can be used.

* `task_type` - (Optional, String, ForceNew) Specifies the scan task type. Changing this parameter will create a new
  resource. Valid values are:
  + **normal**: Normal task type.
  + **monitor**: Monitor task type. The prerequisite for using **monitor** task type is to upgrade the vulnerability
  management service to the Professional editions or above.

  Defaults to **normal**.

* `timer` - (Optional, String, ForceNew) Specifies the scheduled trigger time of the normal task. Changing this parameter
  will create a new resource. This field is valid only when `task_type` is set to **normal**. The field format is
  **yyyy-mm-dd hh:mm:ss**. The trigger time needs to be after the current time. The normal task will start immediately
  when this field is not configured.

* `trigger_time` - (Optional, String, ForceNew) Specifies the scheduled trigger time of the monitor task. Changing this
  parameter will create a new resource. This field is required when `task_type` is set to **monitor**. The field format
  is **yyyy-mm-dd hh:mm:ss**.

* `task_period` - (Optional, String, ForceNew) Specifies the scheduled trigger period of the monitor task. Changing this
  parameter will create a new resource. This field is required when `task_type` is set to **monitor**. Valid values are:
  + **everyday**: Trigger monitor task every day.
  + **threedays**: Trigger monitor task every three days.
  + **everyweek**: Trigger monitor task every week.
  + **everymonth**: Trigger monitor task every month.

* `scan_mode` - (Optional, String, ForceNew) Specifies the task scan mode. Changing this parameter will create a new
  resource. Valid values are:
  + **fast**: Quick scan.
  + **normal**: Normal scan.
  + **deep**: Deep scan.
  
  Defaults to **normal**.

* `port_scan` - (Optional, Bool, ForceNew) Specifies whether to perform port scanning. Changing this parameter will
  create a new resource. Defaults to **false**. Basic, Professional, Advanced and Enterprise editions of vulnerability
  management services support configuring this parameter.

* `weak_pwd_scan` - (Optional, Bool, ForceNew) Specifies whether to scan for weak passwords. Changing this parameter will
  create a new resource. Defaults to **false**.

* `cve_check` - (Optional, Bool, ForceNew) Specifies whether to perform CVE vulnerability scanning. Changing this parameter
  will create a new resource. Defaults to **false**.

* `text_check` - (Optional, Bool, ForceNew) Specifies whether to conduct website content compliance text detection.
  Changing this parameter will create a new resource. Defaults to **false**.

-> Fields `weak_pwd_scan`, `cve_check` and `text_check` are only supported by the Professional, Advanced and Enterprise
editions of the vulnerability management service.

* `picture_check` - (Optional, Bool, ForceNew) Specifies whether to conduct website content compliance image detection.
  Changing this parameter will create a new resource. Defaults to **false**.

* `malicious_code` - (Optional, Bool, ForceNew) Specifies whether to perform malicious code scanning.
  Changing this parameter will create a new resource. Defaults to **false**.

* `malicious_link` - (Optional, Bool, ForceNew) Specifies whether to perform link health detection.
  Changing this parameter will create a new resource. Defaults to **false**.

-> Fields `picture_check`, `malicious_code` and `malicious_link` are only supported by the Enterprise editions of the
vulnerability management service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the task.

* `task_status` - The task status. Valid values are **running**, **success**, **waiting**, **ready** and **failure**.

* `schedule_status` - The monitor task status. Valid values are **running**, **waiting** and **finished**. This field is
  valid only when `task_type` is **monitor**.

* `progress` - The task progress.

* `reason` - The description of task status.

* `pack_num` - The total number of packages.

* `score` - The safety score.

* `safe_level` - The security level. Valid values are **safety**, **average** and **highrisk**.

* `high` - The number of high-risk vulnerabilities.

* `middle` - The number of medium-risk vulnerabilities.

* `low` - The number of low-severity vulnerabilities.

* `hint` - The number of hint-risk vulnerabilities.

## Import

The CodeArts inspector website scan can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_inspector_website_scan.test <id>
```
