---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_quota"
description: |-
  Manages a GaussDB MySQL quota resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_quota

Manages a GaussDB MySQL quota resource within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_gaussdb_mysql_quota" "test" {
  enterprise_project_id = var.enterprise_project_id
  instance_quota        = 10
  vcpus_quota           = 0
  ram_quota             = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Required, String, ForceNew) Specifies the enterprise project ID. Changing this parameter
  will create a new resource.

* `instance_quota` - (Optional, Int) Specifies the instance quantity quota. Value range: **-1** to **100000**. The value
  **-1** indicates no limit. If there are already instances created, this parameter value must be greater than the number
  of existing instances. Defaults to **-1**.

* `vcpus_quota` - (Optional, Int) Specifies the vCPU quota. Value range: **-1** to **2147483646**. The value **-1**
  indicates no limit. If there are already instances created, this parameter value must be greater than the number of
  used vCPUs. Defaults to **-1**.

* `ram_quota` - (Optional, Int) Specifies the memory quota in GB. Value range: **-1** to **2147483646**. The value **-1**
  indicates no limit. If there are already instances created, this parameter value must be greater than the used memory
  size. Defaults to **-1**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `enterprise_project_id`.

* `enterprise_project_name` - Indicates the enterprise project name.

* `availability_instance_quota` - Indicates the remaining instance quota.

* `availability_vcpus_quota` - Indicates the remaining vCPU quota.

* `availability_ram_quota` - Indicates the remaining memory quota.

## Import

The GaussDB MySQL quota can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_quota.test <id>
```
