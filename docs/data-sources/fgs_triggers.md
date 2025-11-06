---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_triggers"
description: |-
  Use this data source to query trigger list under all functions within HuaweiCloud.
---

# huaweicloud_fgs_triggers

Use this data source to query trigger list under all functions within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_triggers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the triggers are located.  
  If omitted, the provider-level region will be used.

* `trigger_type` - (Optional, String) Specifies the type of the trigger.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `triggers` - The list of triggers that match the filter parameters.  
  The [triggers](#fgs_triggers_attr) structure is documented below.

<a name="fgs_triggers_attr"></a>
The `triggers` block supports:

* `trigger_id` - The ID of the trigger.

* `trigger_type_code` - The type of the trigger.

* `trigger_status` - The status of the trigger.
  + **ACTIVE**
  + **DISABLED**

* `event_data` - The detailed configuration of the trigger, in JSON format.

* `func_urn` - The URN of the function to which the trigger belongs.

* `created_time` - The creation time of the trigger.

* `last_updated_time` - The latest update time of the trigger.
