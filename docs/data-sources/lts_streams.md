---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_streams"
description: |-
  Use this data source to get the list of LTS log streams.
---

# huaweicloud_lts_streams

Use this data source to get the list of LTS log streams.

## Example Usage

```hcl
data "huaweicloud_lts_streams" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the log stream.

* `log_group_name` - (Optional, String) Specifies the name of the log group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `streams` - All log streams that match the filter parameters.

  The [streams](#streams_struct) structure is documented below.

<a name="streams_struct"></a>
The `streams` block supports:

* `id` - The ID of the log stream.

* `name` - The name of the log stream.

* `ttl_in_days` - The log expiration time (days).

* `tags` - The key/value pairs to associate with the log stream.

* `created_at` - The creation time of the log stream, in RFC3339 format.
