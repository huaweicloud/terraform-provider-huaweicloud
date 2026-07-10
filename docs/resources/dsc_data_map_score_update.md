---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_data_map_score_update"
description: |-
  Use this resource to update the data map score within HuaweiCloud.
---

# huaweicloud_dsc_data_map_score_update

Use this resource to update the data map score within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not revert the score update,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_dsc_data_map_score_update" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to update the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `request_id` - The unique ID of this request.

* `msg` - The operation result message.

* `status` - The operation status identifier.
