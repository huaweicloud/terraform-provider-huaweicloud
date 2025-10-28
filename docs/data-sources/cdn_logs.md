---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_logs"
description: |-
  Use this datasource to get the list of CDN logs within HuaweiCloud.
---

# huaweicloud_cdn_logs

Use this datasource to get the list of CDN logs within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_cdn_logs" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name to which the logs belong.

* `start_time` - (Optional, Int) Specifies the start time for querying logs.
  The value is the millisecond timestamp of the hour.
  If this parameter is left empty, **00:00:00** of the current day is used by default.

* `end_time` - (Optional, Int) Specifies the end time for querying logs (excluding the end time point).
  The value is the millisecond timestamp of the hour. The maximum time span between the start time and
  end time is `30` days.
  If this parameter is left empty, the start time plus one day is used by default.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource
  belongs.  
  This parameter is valid only when the enterprise project function is enabled.
  The value **all** indicates all enterprise projects.
  This parameter is mandatory when you use an IAM user to call this API.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of logs that matched filter parameters.
  The [logs](#cdn_logs_struct) structure is documented below.

<a name="cdn_logs_struct"></a>
The `logs` block supports:

* `domain_name` - The domain name to which the log belongs.

* `name` - The name of the log file.

* `size` - The size of the log file, in KB.

* `link` - The log file download link.

* `start_time` - The start time for querying log. The value is a timestamp in milliseconds.

* `end_time` - The end time for querying log. The value is a timestamp in milliseconds.
