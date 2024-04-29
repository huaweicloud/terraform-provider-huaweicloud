---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_policy"
description: ""
---

# huaweicloud_cc_central_network_policy

Manages a central network policy resource of Cloud Connect within HuaweiCloud.

## Example Usage

### Create a basic policy

```hcl
variable "central_network_id" {}
variable "project_id" {}
variable "region_id" {}
variable "enterprise_router_id" {}
variable "enterprise_router_table_id" {}

resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = var.central_network_id

  planes {
    associate_er_tables {
      project_id                 = var.project_id
      region_id                  = var.region_id
      enterprise_router_id       = var.enterprise_router_id
      enterprise_router_table_id = var.enterprise_router_table_id
    }
  }

  er_instances {
    project_id           = var.project_id
    region_id            = var.region_id
    enterprise_router_id = var.enterprise_router_id
  }
}
```

### Create a policy that needs to exclude some ER connections
  
```hcl
variable "central_network_id" {}
variable "project_id_1" {}
variable "project_id_2" {}
variable "project_id_3" {}
variable "region_id_1" {}
variable "region_id_2" {}
variable "region_id_3" {}
variable "enterprise_router_id_1" {}
variable "enterprise_router_id_2" {}
variable "enterprise_router_id_3" {}
variable "enterprise_router_table_id_1" {}
variable "enterprise_router_table_id_2" {}
variable "enterprise_router_table_id_3" {}

resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = var.central_network_id
 
  planes {
    associate_er_tables {
      project_id                 = var.project_id_1
      region_id                  = var.region_id_1
      enterprise_router_id       = var.enterprise_router_id_1
      enterprise_router_table_id = var.enterprise_router_table_id_1
    }

    associate_er_tables {
      project_id                 = var.project_id_2
      region_id                  = var.region_id_2
      enterprise_router_id       = var.enterprise_router_id_2
      enterprise_router_table_id = var.enterprise_router_table_id_2
    }

    associate_er_tables {
      project_id                 = var.project_id_3
      region_id                  = var.region_id_3
      enterprise_router_id       = var.enterprise_router_id_3
      enterprise_router_table_id = var.enterprise_router_table_id_3
    }

    exclude_er_connections {
      exclude_er_instances {
        project_id           = var.project_id_1
        region_id            = var.region_id_1
        enterprise_router_id = var.enterprise_router_id_1
      }

      exclude_er_instances {
        project_id           = var.project_id_2
        region_id            = var.region_id_2
        enterprise_router_id = var.enterprise_router_id_2
      }
    }
  }
 
  er_instances {
    project_id           = var.project_id_1
    region_id            = var.region_id_1
    enterprise_router_id = var.enterprise_router_id_1
  }

  er_instances {
    project_id           = var.project_id_2
    region_id            = var.region_id_2
    enterprise_router_id = var.enterprise_router_id_2
  }

  er_instances {
    project_id           = var.project_id_3
    region_id            = var.region_id_3
    enterprise_router_id = var.enterprise_router_id_3
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `central_network_id` - (Required, String, ForceNew) Central network ID.

  Changing this parameter will create a new resource.

* `er_instances` - (Required, List, ForceNew) List of the enterprise routers on the central network policy.

  Changing this parameter will create a new resource.
  The [er_instances](#centralNetworkPolicy_AssociateErInstanceDocument) structure is documented below.

* `planes` - (Optional, List, ForceNew) List of the central network policy planes.

  Changing this parameter will create a new resource.
  The [planes](#centralNetworkPolicy_CentralNetworkPolicyPlaneDocument) structure is documented below.

<a name="centralNetworkPolicy_AssociateErInstanceDocument"></a>
The `er_instances` block supports:

* `project_id` - (Required, String, ForceNew) Project ID.

  Changing this parameter will create a new resource.

* `region_id` - (Required, String, ForceNew) Region ID.

  Changing this parameter will create a new resource.
  
* `enterprise_router_id` - (Required, String, ForceNew) Enterprise router ID.

  Changing this parameter will create a new resource.

<a name="centralNetworkPolicy_CentralNetworkPolicyPlaneDocument"></a>
The `planes` block supports:

* `associate_er_tables` - (Required, List, ForceNew) List of route tables associated with the central network policy.
  The [associate_er_tables](#centralNetworkPolicy_AssociateErTableDocument) structure is documented below.
  Changing this parameter will create a new resource.

* `exclude_er_connections` - (Optional, List, ForceNew) List of the enterprise router connections excluded from the
  central network policy.
  The [exclude_er_connections](#centralNetworkPolicy_ExcludeErConnectionDocument) structure is documented below.
  Changing this parameter will create a new resource.

<a name="centralNetworkPolicy_AssociateErTableDocument"></a>
The `associate_er_tables` block supports:

* `project_id` - (Required, String, ForceNew) Project ID.
  Changing this parameter will create a new resource.

* `region_id` - (Required, String, ForceNew) Region ID.
  Changing this parameter will create a new resource.

* `enterprise_router_id` - (Required, String, ForceNew) Enterprise router ID.
  Changing this parameter will create a new resource.

* `enterprise_router_table_id` - (Required, String, ForceNew) Enterprise router table ID.
  Changing this parameter will create a new resource.

<a name="centralNetworkPolicy_ExcludeErConnectionDocument"></a>
The `exclude_er_connections` block supports:

* `exclude_er_instances` - (Required, List, ForceNew) List of enterprise routers that will not establish a connection.
  The [exclude_er_instances](#centralNetworkPolicy_AssociateErInstanceDocument) structure is the same as `er_instances`.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `document_template_version` - Central network policy document template version.

* `is_applied` - Whether the central network policy is applied.

* `version` - Central network policy version.

* `state` - Central network policy status.
  The valid values are as follows:
    - **AVAILABLE**: The policy is available.
    - **CANCELING**: The policy is being cancelled.
    - **APPLYING**: The policy is being applied.
    - **FAILED**: The operation on the policy failed.
    - **DELETED**: The policy is deleted.

## Import

The central network policy can be imported using
`central_network_id`, `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cc_central_network_policy.test <central_network_id>/<id>
```
