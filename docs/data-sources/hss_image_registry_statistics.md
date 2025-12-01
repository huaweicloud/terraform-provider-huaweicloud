---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_registry_statistics"
description: |-
  Use this data source to get the image registry statistics of HSS within HuaweiCloud.
---

# huaweicloud_hss_image_registry_statistics

Use this data source to get the image registry statistics of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_registry_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

* `registry_type` - (Optional, String) Specifies the image repository type. If this parameter is not specified, all types
  are returned by default. To query multiple types, separate them with commas (,). Valid values are:
  + **Harbor**: harbor
  + **Jfrog**: jfrog
  + **SwrPrivate**: SWR private repository
  + **SwrShared**: SWR shared repository
  + **SwrEnterprise**: SWR enterprise repository
  + **Other**: other repository

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `fail_num` - The number of access exceptions.

* `success_num` - The number of successful accesses.
