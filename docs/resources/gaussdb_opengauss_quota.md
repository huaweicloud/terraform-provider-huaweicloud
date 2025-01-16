---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_quota"
description: |-
  Manages a GaussDB OpenGauss quota resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_quota

Manages a GaussDB OpenGauss quota resource within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_gaussdb_opengauss_quota" "test" {
  enterprise_project_id = var.enterprise_project_id
  instance_quota        = 10
  vcpus_quota           = 0
  ram_quota             = -1
  volume_quota          = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Required, String, ForceNew) Specifies the enterprise project ID. Changing this parameter
  will create a new resource.

* `instance_quota` - (Optional, Int) Specifies the instance quantity quota. Value range: **-1** to **100000**. The value
  **-1** indicates no limit.

* `vcpus_quota` - (Optional, Int) Specifies the vCPU quota. Value range: **-1** to **2147483646**. The value **-1**
  indicates no limit.

* `ram_quota` - (Optional, Int) Specifies the memory quota in GB. Value range: **-1** to **2147483646**. The value **-1**
  indicates no limit.

* `volume_quota` - (Optional, Int) Specifies the storage quota in GB. Value range: **-1** to **2147483646**. The value
  **-1** indicates no limit.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `enterprise_project_id`.

* `enterprise_project_name` - Indicates the enterprise project name.

* `instance_used` - Indicates the used EPS instance quota.

* `vcpus_used` - Indicates the used EPS compute quota.

* `ram_used` - Indicates the used EPS memory quota in GB.

* `volume_used` - Indicates the used EPS storage quota, in GB.

## Import

The GaussDB OpenGauss quota can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_opengauss_quota.test <id>
```
