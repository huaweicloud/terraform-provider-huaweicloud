resource "huaweicloud_cciv2_namespace" "test" {
  name = var.namespace_name
}

resource "huaweicloud_cciv2_deployment" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = var.deployment_name

  selector {
    match_labels = {
      app = "template1"
    }
  }

  template {
    metadata {
      labels = {
        app = "template1"
      }

      annotations = {
        "resource.cci.io/instance-type" = var.instance_type
      }
    }

    spec {
      containers {
        name  = var.container_name
        image = var.container_image

        resources {
          limits = {
            cpu    = var.cpu_limit
            memory = var.memory_limit
          }

          requests = {
            cpu    = var.cpu_limit
            memory = var.memory_limit
          }
        }
      }

      image_pull_secrets {
        name = var.image_pull_secret_name
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template.0.metadata.0.annotations,
    ]
  }
}
