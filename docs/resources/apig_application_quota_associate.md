---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_quota_associate"
description: |-
  Use this resource to bind one or more applications to an application quota (policy) within HuaweiCloud.
---

# huaweicloud_apig_application_quota_associate

Use this resource to bind one or more applications to an application quota (policy) within HuaweiCloud.

~> Please note the following key points before use:<br>1. An application can only be associated with one quota.
   Please ensure that all applications in the script have only one association relationship (if an association
   relationship between an application and a quota is established, the corresponding relationship between the
   application and the original quota should be deleted).<br>2. The action of changing the quota association
   relationship may generate a large number of API calls. Please ensure that the API flow control of the APIG service
   is sufficient.

## Example Usage

```hcl
variable "instance_id" {}
variable "quota_id" {}
variable "associated_application_ids" {
  type = list(string)
}

resource "huaweicloud_apig_application_quota_associate" "test" {
  instance_id = var.instance_id
  quota_id    = var.quota_id

  dynamic "applications" {
    for_each = var.associated_application_ids

    content {
      id = applications.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application quota (policy) is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  quota (policy) belongs.  
  Changing this will create a new resource.

* `quota_id` - (Required, String, ForceNew) Specifies the ID of the application quota (policy).  
  Changing this will create a new resource.

* `applications` - (Required, List) Specifies the configuration of applications bound to the quota (policy).  
  The [applications](#application_quota_associate_config) structure is documented below.

<a name="application_quota_associate_config"></a>
The `applications` block supports:

* `id` - (Required, String) Specifies the application ID bound to the application quota.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `applications` - The configuration of applications bound to the quota (policy).  
  The [applications](#application_quota_associate_config) structure is documented below.

<a name="application_quota_associate_config"></a>
The `applications` block supports:

* `bind_time` - The binding time, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

Quota associate relationship can be imported using related `instance_id` and `id` (also `quota_id`), separated by
a slash, e.g.

```bash
$ terraform import huaweicloud_apig_application_quota_associate.test <instance_id>/<id>
```
