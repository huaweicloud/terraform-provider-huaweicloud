---
layout: "huaweicloud"
page_title: "Provider: HuaweiCloud"
sidebar_current: "docs-huaweicloud-index"
description: |-
  The HuaweiCloud provider is used to interact with the many resources supported by HuaweiCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# HuaweiCloud Provider

The HuaweiCloud provider is used to interact with the
many resources supported by HuaweiCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  region     = "cn-north-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "huaweicloud_vpc_v1" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

## Authentication

The Huawei Cloud provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials ###

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "huaweicloud" {
  region     = "cn-north-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```
-> **NOTE:** `domain_name`, [Account](https://support.huaweicloud.com/en-us/usermanual-iam/iam_01_0552.html) need to be set if using IAM or prePaid resources.


By adding `user_name` and `password`:

```hcl
provider "huaweicloud" {
  region      = "cn-north-1"
  domain_name = "my-account-name"
  user_name   = "my-username"
  password    = "my-password"
}
```

By adding `token`:

```hcl
provider "huaweicloud" {
  region      = "cn-north-1"
  domain_name = "my-account-name"
  token       = "my-token"
}
```


Delegating Resource Access to Another Account by adding `agency_name`, `agency_domain_name` and `delegated_project`:

Usage:

```hcl
provider "huaweicloud" {
  agency_name        = "agency-name"
  agency_domain_name = "agency-domain-name"
  delegated_project  = "delegated-project"

  region     = "cn-north-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `OS_ACCESS_KEY` and
`OS_SECRET_KEY`, environment variables, representing your Huawei
Cloud Access Key and Secret Key, respectively.
The `OS_USERNAME` and `OS_PASSWORD` environment variables
are also used, if applicable:

```hcl
provider "huaweicloud" {}
```

Usage:

```sh
$ export OS_ACCESS_KEY="anaccesskey"
$ export OS_SECRET_KEY="asecretkey"
$ export OS_REGION_NAME="cn-north-1"
$ export OS_DOMAIN_NAME="account-name"
$ terraform plan
```


## Configuration Reference

The following arguments are supported:

* `region` - (Required) This is the Huawei Cloud region. It must be provided,
  but it can also be sourced from the `OS_REGION_NAME` environment variables.

* `access_key` - (Optional) The access key of the HuaweiCloud to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the HuaweiCloud to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `user_id` - (Optional) The User ID to login with. If omitted, the
  `OS_USER_ID` environment variable is used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `tenant_name` - (Optional) The Name of the Tenant/Project to login with.
  If omitted, the `OS_TENANT_NAME` or `OS_PROJECT_NAME` environment variable are used.

* `tenant_id` - (Optional) The ID of the Tenant/Project to login with. If omitted,
  the `OS_TENANT_ID` or `OS_PROJECT_ID` environment variables are used.

* `token` - (Optional) A token is an expiring, temporary means of access issued via
  the IAM service. If omitted, the `OS_AUTH_TOKEN` environment variable is used.

* `domain_name` - (Optional) The Account name of IAM to scope to. If omitted, the
  `OS_DOMAIN_NAME` environment variable is used.

* `domain_id` - (Optional) The Account ID of IAM to scope to. If omitted, the
  `OS_DOMAIN_ID` environment variable is used.

* `auth_url` - (Optional, Required before 1.14.0) The Identity authentication URL. If omitted, the
  `OS_AUTH_URL` environment variable is used. This is not required if you use Huawei Cloud.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

* `agency_name` - (Optional) if authorized by Agencies, it must be set. The
  name of agency.

* `agency_domain_name` - (Optional) if authorized by Agencies, it must be set.
  The name of the account who created the agency.

* `delegated_project` - (Optional) The name of the delegated project.

## Additional Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the HuaweiCloud cloud which is useful for troubleshooting and
debugging.

To enable these logs, set the `OS_DEBUG` environment variable to `1` along
with the usual `TF_LOG=DEBUG` environment variable:

```shell
$ OS_DEBUG=1 TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_REGION_NAME` - The region in which to create the server instance.

* `OS_IMAGE_ID` or `OS_IMAGE_NAME` - a UUID or name of an existing image in
    Glance.

* `OS_FLAVOR_ID` or `OS_FLAVOR_NAME` - an ID or name of an existing flavor.

* `OS_POOL_NAME` - The name of a Floating IP pool.

* `OS_NETWORK_ID` - The UUID of a network in your test environment.

* `OS_EXTGW_ID` - The UUID of the external gateway.

You should be able to use any HuaweiCloud environment to develop on as long as the
above environment variables are set.
