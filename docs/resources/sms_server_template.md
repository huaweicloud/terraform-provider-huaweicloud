---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_server_template"
description: ""
---

# huaweicloud_sms_server_template

Manages an SMS server template resource within HuaweiCloud.

## Example Usage

### A template will create networks during migration

```hcl
data "huaweicloud_availability_zones" "demo" {}

resource "huaweicloud_sms_server_template" "demo" {
  name              = "demo"
  availability_zone = data.huaweicloud_availability_zones.demo.names[0]
}
```

### A template will use the existing networks during migration

```hcl
variable "vpc_id" {}
variable "subent_id" {}
variable "secgroup_id" {}

data "huaweicloud_availability_zones" "demo" {}

resource "huaweicloud_sms_server_template" "demo" {
  name               = "demo"
  availability_zone  = data.huaweicloud_availability_zones.demo.names[0]
  vpc_id             = var.vpc_id
  subnet_ids         = [ var.subent_id ]
  security_group_ids = [ var.secgroup_id ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the server template name.

* `availability_zone` - (Required, String) Specifies the availability zone where the target server is located.

* `region` - (Optional, String) Specifies the region where the target server is located.
  If omitted, the provider-level region will be used.

* `project_id` - (Optional, String) Specifies the project ID where the target server is located.
  If omitted, the default project in the region will be used.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC which the target server belongs to.
  If omitted or set to "autoCreate", a new VPC will be created automatically during migration.

* `subnet_ids` - (Optional, List) Specifies an array of one or more subnet IDs to attach to the target server.
  If omitted or set to ["autoCreate"], a new subnet will be created automatically during migration.

* `security_group_ids` - (Optional, List) Specifies an array of one or more security group IDs to associate with
  the target server.  
  If omitted or set to ["autoCreate"], a new security group will be created automatically during migration.

* `volume_type` - (Optional, String) Specifies the disk type of the target server.
  Available values are: **SAS**, **SSD**, defaults to **SAS**.

* `flavor` - (Optional, String) Specifies the flavor ID for the target server.

* `target_server_name` - (Optional, String) Specifies the name of the target server. Defaults to the template name.

* `bandwidth_size` - (Optional, Int) Specifies the bandwidth size in Mbit/s about the public IP address
  that will be used for migration.  
  The valid value is range from `1` to `2,000`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `vpc_name` - The name of the VPC which the target server belongs to.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

SMS server templates can be imported by `id`, e.g.

```sh
terraform import huaweicloud_sms_server_template.demo 4618ccaf-b4d7-43b9-b958-3df3b885126d
```
