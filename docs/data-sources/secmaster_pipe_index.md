---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_pipe_index"
description: |-
  Use this data source to get the list of SecMaster pipe index.
---

# huaweicloud_secmaster_pipe_index

Use this data source to get the list of SecMaster pipe index.

## Example Usage

```hcl
data "huaweicloud_secmaster_pipe_index" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `pipe_id` - (Required, String) Specifies the pipe ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `mapping` - The mapping information in JSON format of the pipe index.

* `status` - The status of the pipe index. Valid values are **open** and **closed**.

* `timestamp_field` - The timestamp field of the pipe index.
