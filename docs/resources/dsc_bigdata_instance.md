---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_bigdata_instance"
description: |-
  Manages a resource to add a bigdata instance to the DSC asset center within HuaweiCloud.
---

# huaweicloud_dsc_bigdata_instance

Manages a resource to add a bigdata instance to the DSC asset center within HuaweiCloud.

-> This resource is a one-time action resource used to add a bigdata instance to the DSC asset center.
Deleting this resource will not remove the added bigdata instance or undo the add action, but will only
remove the resource information from the tf state file.

## Example Usage

```hcl
variable "asset_name" {}
variable "ds_type" {}
variable "ds_name" {}
variable "ds_address" {}
variable "ds_port" {}
variable "ds_version" {}
variable "ins_type" {}
variable "ins_id" {}
variable "ins_name" {}
variable "ds_user" {}
variable "ds_password" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "huaweicloud_dsc_bigdata_instance" "test" {
  asset_name        = var.asset_name
  ds_type           = var.ds_type
  ds_name           = var.ds_name
  ds_address        = var.ds_address
  ds_port           = var.ds_port
  ds_version        = var.ds_version
  ins_type          = var.ins_type
  ins_id            = var.ins_id
  ins_name          = var.ins_name
  ds_user           = var.ds_user
  ds_password       = var.ds_password
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  scan_metadata     = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `asset_name` - (Optional, String, NonUpdatable) Specifies the asset name.

* `create_time` - (Optional, Int, NonUpdatable) Specifies the asset creation time in milliseconds.

* `ds_address` - (Optional, String, NonUpdatable) Specifies the data source address.

* `ds_name` - (Optional, String, NonUpdatable) Specifies the data source name.

* `ds_password` - (Optional, String, NonUpdatable) Specifies the data source password.

* `ds_port` - (Optional, Int, NonUpdatable) Specifies the data source port.

* `ds_type` - (Optional, String, NonUpdatable) Specifies the data source type.
  Valid values include **Elasticsearch**, **DLI**, **Hive**, **HBase**, **MRS_HIVE**, **ALL**, **LTS**,
  **HIVE_ONLY**, **JUST_BIG_DATA**.

* `ds_user` - (Optional, String, NonUpdatable) Specifies the data source username.

* `ds_version` - (Optional, String, NonUpdatable) Specifies the data source version.

* `ins_id` - (Optional, String, NonUpdatable) Specifies the instance ID.

* `ins_name` - (Optional, String, NonUpdatable) Specifies the instance name.

* `ins_type` - (Optional, String, NonUpdatable) Specifies the instance type.
  The valid values are as follows:
  + **CSS**: Huawei Cloud CSS Elasticsearch service. (ds_type: Elasticsearch)
  + **ECS**: Huawei Cloud ECS self-built Elasticsearch. (ds_type: Elasticsearch)
  + **PUB**: Public network self-built Elasticsearch. (ds_type: Elasticsearch)
  + **EXTERNAL**: External self-built Elasticsearch. (ds_type: Elasticsearch)
  + **MRS**: MRS cluster Hive. (ds_type: Hive/MRS_HIVE)
  + **ECS**: Huawei Cloud ECS self-built Hive. (ds_type: Hive)
  + **PUB**: Public network self-built Hive. (ds_type: Hive)
  + **EXTERNAL**: External self-built Hive. (ds_type: Hive)
  + **DLI**: Data Lake Insight service. (ds_type: DLI)
  + **LTS**: Log Tank Service. (ds_type: LTS)
  + **ECS/PUB/EXTERNAL**: Self-built HBase. (ds_type: HBase)

* `lts_group_id` - (Optional, String, NonUpdatable) Specifies the LTS log group ID.

* `lts_group_name` - (Optional, String, NonUpdatable) Specifies the LTS log group name.

* `lts_stream_id` - (Optional, String, NonUpdatable) Specifies the LTS log stream ID.

* `lts_stream_name` - (Optional, String, NonUpdatable) Specifies the LTS log stream name.

* `queue_name` - (Optional, String, NonUpdatable) Specifies the queue name.

* `region_name` - (Optional, String, NonUpdatable) Specifies the region information.

* `scan_metadata` - (Optional, Bool, NonUpdatable) Specifies whether to scan metadata.

* `security_group_id` - (Optional, String, NonUpdatable) Specifies the security group ID.

* `subnet_id` - (Optional, String, NonUpdatable) Specifies the subnet ID.

* `vpc_id` - (Optional, String, NonUpdatable) Specifies the VPC ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
