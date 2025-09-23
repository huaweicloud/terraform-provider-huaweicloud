---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_ext_data_source"
description: |-
  Manages a GaussDB(DWS) external data source resource within HuaweiCloud.
---

# huaweicloud_dws_ext_data_source

Manages a GaussDB(DWS) external data source resource within HuaweiCloud.

## Example Usage

### Create an external data source to MRS

```hcl
variable "mrs_id" {}
variable "mrs_username" {}
variable "mrs_password" {}
variable "cluster_id" {}

resource "huaweicloud_dws_ext_data_source" "test" {
  type           = "MRS"
  name           = "demo"
  data_source_id = var.mrs_id
  user_name      = var.mrs_username
  user_pwd       = var.mrs_password
  cluster_id     = var.cluster_id
  description    = "This is a demo"
}
```

### Create an external data source to OBS

```hcl
variable "dws_agency_obs" {}
variable "dws_database" {}
variable "cluster_id" {}

resource "huaweicloud_dws_ext_data_source" "test" {
  type         = "OBS"
  name         = "demo"
  user_name    = var.dws_agency_obs
  connect_info = var.dws_database
  cluster_id   = var.cluster_id
  description  = "This is a demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) The DWS cluster ID to which the external data source belongs.
  The cluster **type** must be **ANALYSIS**.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of the external data source.  
  The name can contain `3` to `64` characters. Only letters, digits, and underscores (_) are allowed, and must start with
  a lowercase letter.

  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) The type of the external data source.  
  The valid values are **OBS**, and **MRS**.
  Changing this parameter will create a new resource.
  + If `type` is set to **MRS**, the `type` parameter of the `huaweicloud_mapreduce_cluster` resource must be **ANALYSIS**.
  + If `type` is set to **OBS**, the DWS cluster version must be `8.2.0` or later.

* `user_name` - (Required, String) Specifies the user name of the external data source.  
  It is OBS agency name when `type` is **OBS**.
  This parameter can be modified only when `type` is **OBS**.

  Changing this parameter will create a new resource.

* `data_source_id` - (Optional, String, ForceNew) ID of the external data source. It is mandatory when **type** is **MRS**.

  Changing this parameter will create a new resource.

* `user_pwd` - (Optional, String, ForceNew) The password of the external data source. It is mandatory when **type** is **MRS**.

  Changing this parameter will create a new resource.

* `connect_info` - (Optional, String, ForceNew) The connection information of the external data source.
  It is mandatory when **type** is **OBS**.
  The value is the **database** where the OBS data source connection is to be created.

  Changing this parameter will create a new resource.

* `reboot` - (Optional, Bool) Whether to reboot the cluster.  

* `description` - (Optional, String, ForceNew) The description of the external data source.  

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `configure_status` - The configure status of the external data source.  

* `status` - The status of the external data source.  

* `created_at` - The creation time of the external data source.  

* `updated_at` - The updated time of the external data source.  

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

The dws external data source can be imported using `cluster_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_ext_data_source.test <cluster_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `user_pwd`, `reboot`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dws_ext_data_source" "test" {
  ...

  lifecycle {
    ignore_changes = [
      user_pwd,reboot
    ]
  }
}
```
