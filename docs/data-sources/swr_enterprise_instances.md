---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instances"
description: |-
  Use this data source to get the list of SWR enterprise instances.
---

# huaweicloud_swr_enterprise_instances

Use this data source to get the list of SWR enterprise instances.

## Example Usage

```hcl
data "huaweicloud_swr_enterprise_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the instance.

* `status` - (Optional, String) Specifies the instance status.
  Valid value can be **Initial**, **Creating**, **Running**, **Unavailable**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the instances.
  The [instances](#attrblock--instances) structure is documented below.

<a name="attrblock--instances"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the name of the instance.

* `access_address` - Indicates the access address of instance.

* `charge_mode` - Indicates the charge mode of instance.

* `description` - Indicates the description.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `expires_at` - Indicates the expired time.

* `obs_bucket_name` - Indicates the OBS bucket name.

* `spec` - Indicates the specification of the instance.

* `status` - Indicates the instance status.

* `user_def_obs` - Indicates whether the user specifies the OBS bucket.

* `version` - Indicates the instance version.

* `vpc_id` - Indicates the VPC ID .

* `vpc_name` - Indicates the VPC name.

* `vpc_cidr` - Indicates the range of available subnets for the VPC.

* `subnet_id` - Indicates the subnet ID.

* `subnet_name` - Indicates the subnet name.

* `subnet_cidr` - Indicates the range of available subnets for the subnet.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.
