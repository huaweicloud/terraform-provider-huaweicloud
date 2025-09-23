---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_instances"
description: ""
---

# huaweicloud_cnad_advanced_instances

Use this data source to get the list of CNAD advanced instances.

## Example Usage

```hcl
variable "instance_name" {}

data "huaweicloud_cnad_advanced_instances" "test" {
  instance_name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  + For Unlimited Protection Advanced Edition, support regions are: **cn-north-2**, **cn-north-4**, **cn-east-3**
    and **cn-south-1**.
  + For Unlimited Protection Basic Edition, support regions are: **cn-north-4**, **cn-east-3**, **cn-south-1**
    and **cn-southwest-2**.
  + For Cloud Native Anti-DDoS Standard, support regions are: **cn-north-4**, **cn-east-3** and **cn-south-1**.

* `instance_id` - (Optional, String) Specifies the instance id.

* `instance_name` - (Optional, String) Specifies the instance name.

* `instance_type` - (Optional, String) Specifies the instance type. Valid values are:
  + **cnad_pro**: Professional Edition.
  + **cnad_ip**: Standard Edition.
  + **cnad_ep**: Platinum Edition.
  + **cnad_full_high**: Unlimited Protection Advanced Edition.
  + **cnad_vic**: On demand Version.
  + **cnad_intl_ep**: International Station Platinum Edition.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of the advanced instances
  The [instances](#CNADAdvancedInstances_Instances) structure is documented below.

<a name="CNADAdvancedInstances_Instances"></a>
The `instances` block supports:

* `instance_id` - Indicates the instance id.

* `instance_name` - Indicates the instance name.

* `region` - Indicates the region where the instance belongs to.

* `instance_type` - Indicates the instance type of the instance.

* `protection_type` - Indicates the protection type of the instance.

* `ip_num` - Indicates the ip num of the instance.

* `ip_num_now` - Indicates the current ip num of the instance.

* `protection_num` - Indicates the protection num of the instance, value `9,999` means unlimited times.

* `protection_num_now` - Indicates the current protection num of the instance.

* `created_at` - Indicates the created time.
