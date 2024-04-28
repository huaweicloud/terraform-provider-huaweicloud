---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protection_group"
description: ""
---

# huaweicloud_sdrs_protection_group

Manages an SDRS protection group resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_availability_zone" {}
variable "target_availability_zone" {}
variable "source_vpc_id" {}

data "huaweicloud_sdrs_domain" "test" {}

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = "test_protection_group"
  description              = "test description"
  source_availability_zone = var.source_availability_zone
  target_availability_zone = var.target_availability_zone
  source_vpc_id            = var.source_vpc_id
  domain_id                = data.huaweicloud_sdrs_domain.test.id
} 
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of a protection group. The name can contain a maximum of 64 characters.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), dots (.), underscores (_), and hyphens (-).

* `source_availability_zone` - (Required, String, ForceNew) Specifies the production site AZ of a protection group.
  `source_availability_zone` and `target_availability_zone` must be different availability zones in the same region.

  Changing this parameter will create a new resource.

* `target_availability_zone` - (Required, String, ForceNew) Specifies the disaster recovery site AZ of a protection
  group.
  `target_availability_zone` and `source_availability_zone` must be different availability zones in the same region.

  Changing this parameter will create a new resource.

* `domain_id` - (Required, String, ForceNew) Specifies the ID of an active-active domain.
  You can search `domain_id` with data source `huaweicloud_sdrs_domain`.

  Changing this parameter will create a new resource.

* `source_vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC for the production site.
  One protection group manages servers in one VPC. If you have multiple VPCs, create multiple protection groups.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of a protection group. The description can
  contain a maximum of 64 characters. The value angle brackets (<) and (>) are not allowed.

  Changing this parameter will create a new resource.

* `dr_type` - (Optional, String, ForceNew) Specifies the deployment model. If not specified, **migration** will be used.
  Indicating migration within a VPC. Currently, only **migration** supported.

  Changing this parameter will create a new resource.

* `enable` - (Optional, Bool) Specifies whether enable the protection group start protecting.
  The default value is **false**. It can only be set to true when there's replication pairs within the protection group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SDRS protection group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sdrs_protection_group.test <id>
```
