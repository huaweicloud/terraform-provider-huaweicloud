---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_link_flavors"
description: |-
  Use this data source to query the task flavor info of DRS links within HuaweiCloud.
---

# huaweicloud_drs_link_flavors

Use this data source to query the task flavor info of DRS links within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_drs_link_flavors" "test" {}
```

### Filter by Parameters

```hcl
data "huaweicloud_drs_link_flavors" "test" {
  engine_type   = "mysql"
  job_type      = "sync"
  job_direction = "up"
  net_type      = "vpc"
  node_type     = "high"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the link flavors.  
  If omitted, the provider-level region will be used.

* `engine_type` - (Optional, String) Specifies the engine type of the link.

* `job_type` - (Optional, String) Specifies the job type of the link.  
  The valid values are as follows:
  + **migration**
  + **sync**
  + **cloudDataGuard**
  + **subscription**
  + **replay**
  + **verify**
  + **cdc**

* `job_direction` - (Optional, String) Specifies the job direction of the link.  
  The valid values are as follows:
  + **up**
  + **down**
  + **non-dbs**

* `net_type` - (Optional, String) Specifies the network type of the link.  
  The valid values are as follows:
  + **eip**
  + **vpc**
  + **vpn**

* `node_type` - (Optional, String) Specifies the node type of the link.  
  The valid values are as follows:
  + **micro**
  + **small**
  + **medium**
  + **high**
  + **xlarge**
  + **2xlarge**

* `is_multi_write` - (Optional, Bool) Specifies whether the disaster recovery is dual-master.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vm_flavor` - The VM flavor information.  
  The [vm_flavor](#drs_link_flavors_vm_flavor) structure is documented below.

* `volume_flavor` - The volume flavor information.  
  The [volume_flavor](#drs_link_flavors_volume_flavor) structure is documented below.

* `flow_flavor` - The flow flavor information.  
  The [flow_flavor](#drs_link_flavors_flow_flavor) structure is documented below.

<a name="drs_link_flavors_vm_flavor"></a>
The `vm_flavor` block supports:

* `id` - The flavor ID.

* `cloud_service_type` - The cloud service type.

* `spec_type_code` - The spec type code.

* `spec_code` - The spec code.

<a name="drs_link_flavors_volume_flavor"></a>
The `volume_flavor` block supports:

* `id` - The flavor ID.

* `cloud_service_type` - The cloud service type.

* `spec_type_code` - The spec type code.

* `spec_code` - The spec code.

<a name="drs_link_flavors_flow_flavor"></a>
The `flow_flavor` block supports:

* `id` - The flavor ID.

* `cloud_service_type` - The cloud service type.

* `spec_type_code` - The spec type code.

* `spec_code` - The spec code.
