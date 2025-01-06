---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_access_logs"
description: |-
  Use this data source to get the list of the access logs.
---

# huaweicloud_ga_access_logs

Use this data source to get the list of the access logs.

## Example Usage

```hcl
data "huaweicloud_ga_access_logs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `log_id` - (Optional, String) Specifies the ID of the access log.

* `status` - (Optional, String) Specifies the status of the access log.
  The valid values are as follows:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `resource_type` - (Optional, String) Specifies the type of the resource to which the access log belongs.
  Currently, only **LISTENER** is supported.

* `resource_ids` - (Optional, List) Specifies the ID list of the resource to which the access log belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of the access logs.

  The [logs](#logs_struct) structure is documented below.

<a name="logs_struct"></a>
The `logs` block supports:

* `id` - The ID of the access log.

* `status` - The status of the access log.

* `resource_type` - The type of the resource to which the access log belongs.

* `resource_id` - The ID of the resource to which the access log belongs.

* `log_group_id` - The ID of the log group to which the access log belongs.

* `log_stream_id` - The ID of the log stream to which the access log belongs.

* `created_at` - The creation time of the access log, in RFC3339 format.

* `updated_at` - The latest update time of the access log, in RFC3339 format.
