---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack"
description: ""
---

# huaweicloud_rfs_stack

Provides an RFS resource stack.

## Example Usage

### Create an RFS resource stack with resource deployment (using OBS URIs)

```hcl
variable "stack_name" {}
variable "agency_name" {}
variable "template_obs_uri" {}
variable "variable_obs_uri" {}

resource "huaweicloud_rfs_stack" "test" {
  name = var.stack_name

  agency {
    name          = var.agency_name
    provider_name = "huaweicloud"
  }

  template_uri = var.template_obs_uri
  vars_uri     = var.variable_obs_uri

  retain_all_resources = true
}
```

### Create an RFS resource stack with VPC deployment (using template and variable files)

```hcl
variable "stack_name" {}
variable "agency_name" {}
variable "template_path" {}
variable "variable_path" {}

resource "huaweicloud_rfs_stack" "test" {
  name = var.stack_name

  agency {
    name          = var.agency_name
    provider_name = "huaweicloud"
  }

  template_body = file(var.template_path) // local storage path of HCL/JSON script
  vars_body     = file(var.variable_path) // local storage path of .vars file
}
```

The content of the template file (in JSON format) is as follows:

```json
{
  "terraform": {
    "required_providers": [
      {
        "huaweicloud": {
          "source": "huawei.com/provider/huaweicloud",
          "version": ">= 1.41.0"
        }
      }
    ]
  },
  "provider": {
    "huaweicloud": {
      "region": "${var.region_name}"
    }
  },
  "resource": {
    "huaweicloud_vpc": {
      "test": {
        "name": "${var.vpc_name}",
        "cidr": "192.168.0.0/16"
      }
    },
    "huaweicloud_vpc_subnet": {
      "test": {
        "vpc_id": "${huaweicloud_vpc.test.id}",
        "name": "${var.subnet_name}",
        "cidr": "${cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)}",
        "gateway_ip": "${cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)}"
      }
    }
  },
  "variable": {
    "region_name": {
      "type": "string"
    },
    "vpc_name": {
      "type": "string"
    },
    "subnet_name": {
      "type": "string"
    }
  }
}
```

The content of the `.vars` file is as follows:

```hcl
region_name = "cn-north-4"
vpc_name    = "tf-example-vpc"
subnet_name = "tf-example-vpc-subnet"
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the RFS resource stack is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the resource stack.  
  The valid length is limited from `1` to `64`, only letters, digits and hyphens (-) are allowed.
  The name must start with a lowercase letter and end with a lowercase letter or digit.
  Change this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the resource stack, which contain maximum of
  `255` characters.  
  Change this parameter will create a new resource.

* `agency` - (Optional, List, ForceNew) Specifies the configuration of the agencies authorized to IAC.  
  Change this parameter will create a new resource.
  The [object](#stack_agency) structure is documented below.

* `template_body` - (Optional, String) Specifies the HCL/JSON template content for deployment resources.  
  This parameter and `template_uri` are alternative and required if `vars_body` is set.

* `vars_body` - (Optional, String) Specifies the variable content for deployment resources.  
  This parameter and `vars_uri` are alternative.

* `template_uri` - (Optional, String) Specifies the OBS address where the HCL/JSON template archive (**.zip** file,
  which contains all resource **.tf.json** script files to be deployed) or **.tf.json** file is located, which describes
  the target status of the deployment resources.

* `vars_uri` - (Optional, String) Specifies the OBS address where the variable (**.tfvars**) file corresponding to the
  HCL/JSON template located, which describes the target status of the deployment resources.

* `enable_auto_rollback` - (Optional, Bool, ForceNew) Specifies whether to enable automatic rollback.  
  If enabled, the stack resources will rollback automatically to the last stable state with deployment failure.
  The default value is **false**.
  Change this parameter will create a new resource.

* `enable_deletion_protection` - (Optional, Bool, ForceNew) Specifies whether to enable delete protection.  
  The default value is **false**.
  Change this parameter will create a new resource.

* `retain_all_resources` - (Optional, Bool) Specifies whether to reserve resources when deleting the resource stack.  
  The default value is **false**.

<a name="stack_agency"></a>
The `agency` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of IAM agency authorized to IAC account.  
  Change this parameter will create a new resource.

* `provider_name` - (Required, String, ForceNew) Specifies the name of the provider corresponding to the IAM agency.  
  Change this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The stack ID.

* `status` - The current status of the resource stack.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 5 minutes.

For most HCL templates, the timeout parameters needs to be manually configured by the user to ensure that resources can
be deployed successfully on the RFS resource stack, e.g.

```hcl
resource "huaweicloud_rfs_stack" "test" {
  ...

  timeouts {
    create = "1h" // Such as the creation of GaussDB instances.
    update = "1h"
  }
}
```

## Import

Stacks can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_rfs_stack.test edd2f099-e1ac-4bd0-be32-8b2185620a90
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `agency`, `template_body`, `vars_body`, `template_uri`, `vars_uri`,
`enable_auto_rollback`, `enable_deletion_protection` and `retain_all_resources`.
It is generally recommended running `terraform plan` after importing a stack.
You can keep the resource the same with its definition bo choosing any of them to update.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_rfs_stack" "test" {
  ...

  lifecycle {
    ignore_changes = [
      agency,
      template_body,
      ...
    ]
  }
}
```
