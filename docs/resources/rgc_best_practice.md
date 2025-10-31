---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_best_parctice"
description: |-
  Manages an RGC best-practice resource within HuaweiCloud.
---

# huaweicloud_rgc_best_parctice

Manages an RGC best-practice query resource within HuaweiCloud. This resource is only supported in cn-north-4 and ap-southeast-3.

## Example Usage

```hcl
resource "huaweicloud_rgc_best_practice" "best_practice" {}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

The following arguments are supported:

* `id` - The resource ID.

* `total_score` - The total score of the most recent best practice.

* `detect_time` - The detect time of the most recent best practice.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

## Import

The RGC best practice can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rgc_best_parctice.test <id>
```
