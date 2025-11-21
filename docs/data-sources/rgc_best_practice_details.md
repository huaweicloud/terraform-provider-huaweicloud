---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_best_practice_details"
description: |-
  Use this data source to list the details of best-practice in Resource Governance Center.
---

# huaweicloud_rgc_best_practice_details

Use this data source to list the details of best-practice in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_best_practice_details" "best_practice" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The data source ID.

* `details` - Information about the details of best-practice list.

The [details](#details) structure is documented below.

<a name="details"></a>
The `details` block supports:

* `check_item` - The check item number associated with the best practice.

* `check_item_name` - The name of the check item.
  
* `risk_description` - A description of the risk associated with the best practice.
  
* `result` - The result of the best practice check.
  
* `scene` - The scene or context in which the best practice is applied.
  
* `risk_level` - The risk level associated with the best practice.
