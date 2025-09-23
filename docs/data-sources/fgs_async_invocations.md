---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_async_invocations"
description: |-
  Use this data source to query async invocations of the specified function within HuaweiCloud.
---

# huaweicloud_fgs_async_invocations

Use this data source to query async invocations of the specified function within HuaweiCloud.

## Example Usage

### Query all async invocations of the specified function

```hcl
variable "function_urn" {}

data "huaweicloud_fgs_async_invocations" "test" {
  function_urn = var.function_urn
}
```

### Query all success async invocations in the specified time window

```hcl
variable "function_urn" {}

data "huaweicloud_fgs_async_invocations" "test" {
  function_urn     = var.function_urn
  status           = "SUCCESS"
  query_begin_time = "2025-01-01T00:00:00Z"
  query_end_time   = "2025-01-31T23:59:59Z"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the async invocations are located.  
  If omitted, the provider-level region will be used.

* `function_urn` - (Required, String) Specifies the function URN to which the async invocations belong.

* `request_id` - (Optional, String) Specifies the specified request ID of async invocation to be queried.

* `status` - (Optional, String) Specifies the status of async invocations to be queried.  
  The valid values are as follows:
  + **WAIT**
  + **RUNNING**
  + **SUCCESS**
  + **FAIL**
  + **DISCARD**

* `query_begin_time` - (Optional, String) Specifies the begin time to query async invocations, in RFC3339 format (UTC
  time).  
  For example, `1999-01-01T08:00:00Z`.

* `query_end_time` - (Optional, String) Specifies the end time to query async invocations, in RFC3339 format (UTC
  time).  
  For example, `1999-01-01T08:00:00Z`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `invocations` - The list of async invocations that matched filter parameters.  
  The [invocations](#fgs_async_invocations_attr) structure is documented below.

<a name="fgs_async_invocations_attr"></a>
The `invocations` block supports:

* `request_id` - The request ID of the async invocation.

* `status` - The status of the async invocation.

* `error_code` - The error code of the async invocation.

* `error_message` - The error message of the async invocation.

* `start_time` - The start time of the async invocation, in RFC3339 format (UTC time).

* `end_time` - The end time of the async invocation, in RFC3339 format (UTC time).
