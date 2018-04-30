job "dp-notification-sender" {
  datacenters = ["eu-west-1"]
  region      = "eu"
  type        = "service"

  group "publishing" {
    count = "{{PUBLISHING_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "publishing"
    }

    task "dp-notification-sender" {
      driver = "exec"

      artifact {
        source = "s3::https://s3-eu-west-1.amazonaws.com/{{BUILD_BUCKET}}/dp-notification-sender/{{REVISION}}.tar.gz"
      }

      artifact {
        source = "s3::https://s3-eu-west-1.amazonaws.com/{{DEPLOYMENT_BUCKET}}/dp-notification-sender/{{REVISION}}.tar.gz"
      }

      config {
        command = "${NOMAD_TASK_DIR}/start-task"

        args    = [
          "${NOMAD_TASK_DIR}/dp-notification-sender",
        ]
      }

      service {
        name = "dp-notification-sender"
        port = "http"
        tags = ["publishing"]
        check {
          type     = "http"
          path     = "/healthcheck"
          interval = "10s"
          timeout  = "2s"
        }
      }

      resources {
        cpu    = "{{PUBLISHING_RESOURCE_CPU}}"
        memory = "{{PUBLISHING_RESOURCE_MEM}}"

        network {
          port "http" {}
        }
      }

      template {
        source      = "${NOMAD_TASK_DIR}/vars-template"
        destination = "${NOMAD_TASK_DIR}/vars"
      }

      vault {
        policies = ["dp-notification-sender"]

        change_mode   = "signal"
        change_signal = "SIGUSR1"
      }
    }
  }
}