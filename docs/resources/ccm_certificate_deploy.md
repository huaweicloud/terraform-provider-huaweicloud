---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_deploy"
description: |-
  Manages a CCM SSL certificate deploy resource within HuaweiCloud.
---

# huaweicloud_ccm_certificate_deploy

Manages a CCM SSL certificate deploy resource within HuaweiCloud.

-> 1. Currently, this resource only supports deploying SSL certificates to **CDN**, **WAF** or **ELB**.
<br/>2. Each successful deployment of a certificate to a service resource incurs a cost.
<br/>3. The current resource is a one-time resource, and destroying this resource will not affect the result.

## Example Usage

### Deploy the certificate to CDN

```hcl
variable "certificate_id" {}
variable "domain_name" {}

resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = var.certificate_id
  service_name   = "CDN"

  resources {
    domain_name = var.domain_name
  }
}
```

### Deploy the certificate to WAF

```hcl
variable "certificate_id" {}
variable "project_name" {}
variable "waf_certificate_id" {}
variable "waf_type" {}
variable "waf_eps_id" {}

resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = var.certificate_id
  project_name   = var.project_name
  service_name   = "WAF"

  resources {
    id                    = var.waf_certificate_id
    type                  = var.waf_type
    enterprise_project_id = var.waf_eps_id
  }
}
```

### Deploy the certificate to ELB

```hcl
variable "certificate_id" {}
variable "project_name" {}
variable "elb_certificate_id" {}

resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = var.certificate_id
  project_name   = var.project_name
  service_name   = "ELB"

  resources {
    id = var.elb_certificate_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `certificate_id` - (Required, String, ForceNew) Specifies the CCM SSL certificate ID to be deployed.
  Changing this parameter will create a new resource.

  -> Certificates encrypted with SM series cryptographic algorithms cannot be deployed to other cloud services.

* `service_name` - (Required, String, ForceNew) Specifies the target service name to which the certificate is pushed.
  Valid values are **CDN**, **WAF**, and **ELB**.

  Changing this parameter will create a new resource.

* `resources` - (Required, List, ForceNew) Specifies the list of resources to be deployed.
  The [resources](#resources_struct) structure is documented below.

  Changing this parameter will create a new resource.

* `project_name` - (Optional, String, ForceNew) Specifies the project name where the deployed resources are located.
  This field is required only when `service_name` is set to **WAF** or **ELB**.

  Changing this parameter will create a new resource.

<a name="resources_struct"></a>
The `resources` block supports:

* `id` - (Optional, String, ForceNew) Specifies the certificate ID. This field is required only when `service_name` is
  set to **WAF** or **ELB**.

  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the resource type. This field is required only when `service_name` is
  set to **WAF**. Valid values are **premium** (exclusive mode) and **cloud** (cloud mode).

  Changing this parameter will create a new resource.

* `domain_name` - (Optional, String, ForceNew) Specifies the domain name to be deployed. This field is required only
  when `service_name` is set to **CDN**. The domain name must match the certificate.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the resources
  to be deployed. This field is required only when `service_name` is set to **WAF**.
  If omitted, default enterprise project will be used.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
