---
subcategory: "CodeArts Inspector""
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_inspector_website_scan_tasks"
description: |-
  Use this data source to get the list of CodeArts inspector website scan tasks.
---

# huaweicloud_codearts_inspector_website_scan_tasks

Use this data source to get the list of CodeArts inspector website scan tasks.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_codearts_inspector_website_scan_tasks" "test" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String) Specifies the domain ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the tasks list.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `task_id` - Indicates the task ID.

* `task_name` - Indicates the task name.

* `task_type` - Indicates the task type.
  Value can be as follows:
  + **normal**: Normal task type.
  + **monitor**: Monitor task type.

* `url` - Indicates the destination URL to scan.

* `task_status` - Indicates the task status.
  Value can be **running**, **success**, **waiting**, **ready** and **failure**.

* `reason` - Indicates the description of task status.

* `created_at` - Indicates the create time of the task.

* `start_time` - Indicates the start time of the task.

* `end_time` - Indicates the end time of the task.

* `schedule_status` - Indicates the monitor task status.
  Value can be **running**, **waiting** and **finished**.

* `safe_level` - Indicates the security level.
  Value can be **safety**, **average** and **highrisk**.

* `progress` - Indicates the task progress.

* `pack_num` - Indicates the total number of packages.

* `score` - Indicates the safety score.

* `domain_name` - Indicates the domain name.

* `high` - Indicates the number of high-risk vulnerabilities.

* `middle` - Indicates the number of medium-risk vulnerabilities.

* `low` - Indicates the number of low-severity vulnerabilities.

* `hint` - Indicates the number of hint-risk vulnerabilities.

* `task_period` - Indicates the scheduled trigger period of the monitor task.
  Value can be as follows:
  + **everyday**: Trigger monitor task every day.
  + **threedays**: Trigger monitor task every three days.
  + **everyweek**: Trigger monitor task every week.
  + **everymonth**: Trigger monitor task every month.

* `timer` - Indicates the scheduled trigger time of the normal task.

* `trigger_time` - Indicates the scheduled trigger time of the monitor task.

* `malicious_link` - Indicates whether to perform link health detection.

* `scan_mode` - Indicates the task scan mode.
  Value can be as follows:
  + **fast**: Quick scan.
  + **normal**: Normal scan.
  + **deep**: Deep scan.

* `port_scan` - Indicates whether to perform port scanning.

* `weak_pwd_scan` - Indicates whether to scan for weak passwords.

* `cve_check` - Indicates whether to perform CVE vulnerability scanning.

* `text_check` - Indicates whether to conduct website content compliance text detection.

* `picture_check` - Indicates whether to conduct website content compliance image detection.

* `malicious_code` - Indicates whether to perform malicious code scanning.
