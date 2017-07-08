version = "1"

space "test-space" {
  name = "Mein Test Space"
  team = "helmich"

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
    }
  }

  stage development {
    inherit = "production"
    application typo3 {
      version = "~8.7"
    }
  }
}