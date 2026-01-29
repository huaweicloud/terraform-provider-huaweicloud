---
subcategory: "LakeFormation"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lakeformation_instance"
description: |-
  Manages a LakeFormation instance resource within HuaweiCloud.
---

# huaweicloud_lakeformation_instance

Manages a LakeFormation instance resource within HuaweiCloud.

-> For more information on the specification configurations of the dedicated instance, please contact service OnCall via
   O&M (Tickets) system.

## Example Usage

### Create a shared instance

```hcl
variable "instance_name" {}

resource "huaweicloud_lakeformation_instance" "test" {
  name        = var.instance_name
  shared      = true
  description = "Created by terraform script"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }

  # When deleting, move into the recycle bin instead of deleting it directly.
  # You can manually manage it after 15 minutes.
  to_recycle_bin = true
}
```

### Create a dedicated instance

```hcl
variable "instance_name" {}

resource "huaweicloud_lakeformation_instance" "test" {
  name        = var.instance_name
  shared      = false
  description = "Created by terraform script"

  specs {
    spec_code  = "lakeformation.unit.basic.qps"
    stride_num = 1 # 2000 QPS (1 * 2000)
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the instance is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the instance.

* `shared` - (Required, Bool, NonUpdatable) Specifies whether the instance is shared.

* `specs` - (Optional, List) Specifies the list of custom specifications for dedicated instance.  
  The [specs](#lakeformation_instance_specs) structure is documented below.

* `description` - (Optional, String) Specifies the description of the instance.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

* `to_recycle_bin` - (Optional, Bool) Specifies whether to put the instance into the recycle bin when deleting.  
  Default value is `false`.

  -> This parameter is only useful during the deletion phase. It cannot be deleted via the console or other ways within
     fifteen minutes of being placed in the recycle bin.

  ~> Instances in the recycle bin will continue to be billed.

<a name="lakeformation_instance_specs"></a>
The `specs` block supports:

* `spec_code` - (Required, String) Specifies the specification code.  
  The valid values ​​can be obtained by data source `huaweicloud_lakeformation_specifications`.

* `stride_num` - (Optional, Int) Specifies the stride number of the specification.  

  -> This parameter is only available when the value of parameter `shared` is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `specs` - The list of custom specifications for dedicated instance.  
  The [specs](#lakeformation_instance_specs_attr) structure is documented below.

* `status` - The status of the instance.

* `is_default` - Whether the instance is the default instance.

* `create_time` - The creation time of the instance, in RFC3339 format.

* `update_time` - The update time of the instance, in RFC3339 format.

<a name="lakeformation_instance_specs_attr"></a>
The `specs` block supports:

* `product_id` - (Optional, String) Specifies the product ID of the specification.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 hours.
* `update` - Default is 2 hours.
* `delete` - Default is 20 minutes.

## Import

LakeFormation instances can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_lakeformation_instance.test <id>
```

After importing this resource, the value of the `to_recycle_bin` parameter stored in the `tfstate` file defaults to
`false`. If this value differs from the current script, please apply the script value by executing the `terraform apply`
command.
