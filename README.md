
A CLi to .......


## Usage

The command `devai` (/bin/devai) can be used to ....


## Structure 

The following strucure was selected to aid with;

* **Code Reusability**: Avoids code duplication by centralizing shared functionality.
* **Maintainability**: Makes it easier to update and test shared code in one place.
* **Modularity**: Encourages better organization and separation of concerns.
* **Testability**: Shared packages can have their own unit tests.
```
devai/
├── cmd/
│   ├── root.go
│   ├── echo/
│   │   ├── echo.go
│   │   └── flags.go (if needed)
│   ├── release/
│   │   ├── release.go
│   │   └── ...
│   └── review/
│       ├── review.go
│       └── ...
├── pkg/
│   ├── utils/
│   │   ├── helpers.go
│   │   └── ...
│   └── config/
│       ├── config.go
│       └── ...
├── go.mod
├── go.sum
└── README.md
```

* `cmd/`: Remains the same, housing your root command and subcommand directories.
* `pkg/`: This is where you'll place shared packages that can be imported by multiple subcommands (or even other parts of your project).
  * `utils/`: A common package for utility functions, like string manipulation, data validation, logging, etc.
  * `config/`: A package for handling configuration loading, parsing, and access (using Viper or similar).
  * Other Packages: You can create more packages as needed, e.g., `api/` for interacting with external APIs, `db/` for database access, etc.


  ## Development

  
  ```
  go run main.go
  go run main.go --help
  go run main.go review code
  ```


  Migrate `buildey` to `devai` from [inital experiments](https://gitlab.com/robedwards/buildey.git).
