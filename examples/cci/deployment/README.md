# Create a CCI Deployment

This example provides best practice code for using Terraform to create a Cloud Container Instance (CCI) deployment
in HuaweiCloud. The example demonstrates how to create a CCI namespace and CCI deployment with configurable container
specifications.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* CCI service enabled in target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where CCI deployment is located
* `access_key` - The access key of IAM user
* `secret_key` - The secret key of IAM user

### Resource Variables

#### Required Variables

* `deployment_name` - The name of CCI deployment
* `namespace_name` - The name of CCI namespace

#### Optional Variables

* `instance_type` - The instance type of CCI pod (default: "general-computing")
* `container_name` - The name of container (default: "c1")
* `container_image` - The image of container (default: "alpine:latest")
* `cpu_limit` - The CPU limit of container (default: "1")
* `memory_limit` - The memory limit of container (default: "2G")
* `image_pull_secret_name` - The name of image pull secret (default: "imagepull-secret")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  deployment_name = "tf-test-deployment"
  namespace_name  = "tf-test-namespace"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Deployment Configuration

### Selector Configuration

The deployment uses a selector to match pods:

```hcl
selector {
  match_labels = {
    app = "template1"
  }
}
```

The selector labels must match the pod template labels.

### Pod Template Configuration

The pod template defines the specification for pods created by the deployment:

```hcl
template {
  metadata {
    labels = {
      app = "template1"
    }

    annotations = {
      "resource.cci.io/instance-type" = "general-computing"
    }
  }

  spec {
    containers {
      name  = "c1"
      image = "alpine:latest"

      resources {
        limits = {
          cpu    = "1"
          memory = "2G"
        }

        requests = {
          cpu    = "1"
          memory = "2G"
        }
      }
    }

    image_pull_secrets {
      name = "imagepull-secret"
    }
  }
}
```

### Lifecycle Configuration

The deployment ignores changes to the template annotations to prevent unnecessary updates:

```hcl
lifecycle {
  ignore_changes = [
    template.0.metadata.0.annotations,
  ]
}
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The deployment must be created within an existing namespace
* The selector labels must match the pod template labels
* The image pull secret must be created beforehand if needed
* All resources will be created in the specified region
* Deployment names must be unique within the namespace
* Make sure to have sufficient quota for the resources you plan to create
* The instance type affects the pricing and performance of the pods

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.74.0 |
