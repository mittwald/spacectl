version = "1"

space "test-space" {
  name = "Mein Test Space"
  team = "helmich"

  payment {
    paymentProfile = "UUID"
    plan = "spaces.flex/v1"
  }

  resource storage {
    quantity = 50
  }

  resource stages {
    quantity = 1
  }

  option backupIntervalMinutes {
    value = 60
  }

  stage production {
    application typo3 {
      version = "8.7.2"
      userData {
        initialUser {
          username = "admin"
          password = "test123"
        }
      }
    }

    database mysql {
      version = "5.7.*"
      storage {
        size = "32GB"
      }
    }

    cron typo3 {
      schedule = "*/5 * * * *"
      allowParallel = true
      command {
        command = "vendor/bin/typo3cmd"
        arguments = ["typo3:scheduler"]
        workingDirectory = "/var/www/typo3"
      }
      timezone = "Europe/Berlin"
    }
  }

  stage development {
    inherit = "production"
    onDemand = true
    application typo3 {
      version = "~8.7"
    }
  }
}