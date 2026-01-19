

# codecov.yml Reference

> 🚧 Changing your YAML?
>
> Changing your YAML? A reminder to always validate your YAML before you deploy [https://api.codecov.io/validate](https://api.codecov.io/validate)

# Top-level sections used in codecov.yml

## codecov

```yaml
codecov:
  token: "<some token here>"
  bot: "codecov-io"
  ci:
    - "travis.org"
  strict_yaml_branch: "yaml-config"
  max_report_age: 24
  disable_default_path_fixes: no
  require_ci_to_pass: yes
  notify:
    after_n_builds: 2
    wait_for_ci: yes
```

### codecov.token

The repository upload token

* Type: String
* Default: N/A

***

### codecov.bot

The username you want to use for Codecov operations

* Type: String
* Default: N/A
* More Info: [https://docs.codecov.io/docs/team-bot](https://docs.codecov.io/docs/team-bot)

***

### codecov.ci

Additional CI provider URLs you want Codecov to recognize.

* Type: Map
* Default: N/A
* More Info: [https://docs.codecov.io/docs/detecting-ci-services](https://docs.codecov.io/docs/detecting-ci-services)

***

### codecov.strict\_yaml\_branch

Specify a branch you want Codecov to always only read the YAML from

* Type: String
* Default: N/A
* More Info: [https://docs.codecov.io/docs/codecov-yaml#section-restricting-changes](https://docs.codecov.io/docs/codecov-yaml#section-restricting-changes)

***

### codecov.max\_report\_age

The age you want coverage reports to expire at, or if you want to disable this check. Expired reports will not be processed by codecov.

* Type: Int, String, Boolean
* Default: 12h
* More Info: [https://docs.codecov.io/docs/codecov-yaml#section-expired-reports](https://docs.codecov.io/docs/codecov-yaml#section-expired-reports)

***

### codecov.disable\_default\_path\_fixes

Should Codecov's default path fixes be disabled

* Type: Boolean
* Default: false
* More Info: [https://docs.codecov.io/docs/fixing-paths](https://docs.codecov.io/docs/fixing-paths)

***

### codecov.require\_ci\_to\_pass

Should Codecov wait for all other statuses to pass before sending its status.

* Type: Boolean
* Default: true

***

### codecov.notify

#### codecov.notify.after\_n\_builds

How many uploaded reports Codecov should wait to receive before sending statuses

* Type: Int
* Default: 1
* More Info: [https://docs.codecov.io/docs/notifications#section-preventing-notifications-until-after-n-builds](https://docs.codecov.io/docs/notifications#section-preventing-notifications-until-after-n-builds)

***

#### codecov.notify.wait\_for\_ci

Should Codecov wait for all CI statuses to complete before sending ours.\
Note: Codecov considers all non-codecov statuses to be CI statuses

* Type: Boolean
* Default: true

***

#### codecov.notify.notify\_error

If any uploads fail to process replace the regular coverage comment with a comment containing the number of uploads that failed.

* Type: Boolean
* Default: false

***

#### codecov.notify.manual\_trigger

```yaml YAML
codecov:
  notify:
    manual_trigger: true
```

Codecov will trigger no notifications or statuses until explicitly told do so by the CLI. This setting is helpful if your CI upload pipeline uploads a variable number of coverage reports per commit, and you would prefer to get no notifications from Codecov until Codecov is told to do so (e.g., at the end of the CI process).

Note that this feature *requires* the use of Codecov's CLI, and no notifications will be triggered until the CLI's [send-notifications command](https://github.com/codecov/codecov-cli#send-notifications) is executed in your CI run.

* Type: Boolean
* Default: false

***

## component\_management

```yaml YAML
component_management:
    default_rules:  # Dict. rules inherited by all components that don't define a tag for themselves.
    paths: # List. Path filters.
        - "specific_file.txt
      - "^some.*regex$"
      - "glob/*"
    flag_regexes: # List. Flags to be included in the component.
        - "fruit.*"
    statuses: # List. Status definitions.
    # These statuses are the same as for flag_management
    # Except they don't accept 'flags'
  individual_components:  # List. These are the actual components.
    - component_id: component_1  # String. Required.
        name: display_name_1    # String. Optional.
      # Individual components also accept paths, flag_regexes and statuses
    - component_id: other_component
    #...
```

More info: [https://docs.codecov.com/docs/components](https://docs.codecov.com/docs/components)

## coverage

```yaml
coverage:
  precision: 2
  round: down
  range: "70...100"
  notify:
    # notification blocks. See: https://docs.codecov.io/docs/codecovyml-reference#section-coverage-notify
  status:
    project:
    patch:
    changes:
```

***

### coverage.precision

What precision do you want the coverage value to be

* Type: Range(0,5)
* Default: 2

***

### coverage.round

Which direction to you want to round the coverage value

* Type: One of 'down', 'up', 'nearest'
* Default: down

***

### coverage.range

The value range where you want the value to be green

* Type: start...end
* Default: 70..100
* More Info: [https://docs.codecov.io/docs/coverage-configuration#section-range](https://docs.codecov.io/docs/coverage-configuration#section-range)

***

### coverage.notify

The standard notification settings are shown below, but it is recommend to view [https://docs.codecov.io/docs/notifications](https://docs.codecov.io/docs/notifications) for more information.

```yaml
<notification_provider>: #see: ttps://docs.codecov.io/docs/notifications
  url: "https://hooks.example.com/hook/8675309"
  branches': 
    - master
    - dev
    - staging
  threshold': 1%
  flags: 
    - backend
    - frontend
  base: "parent"
  only_pulls: false
  paths: "*/**/*"
```

Note that the following notification providers are supported and will need to be defined in the YAML (see [Notifications](https://docs.codecov.io/docs/notifications)):

* Gitter
* IRC
* Slack

***

### coverage.status

More info: [https://docs.codecov.io/docs/commit-status](https://docs.codecov.io/docs/commit-status)

```yaml
coverage:
  status:
    project:
      default: # This can be anything, but it needs to exist as the name
        # basic settings
        target: auto
        threshold: 5%
        base: auto 
        # advanced settings
        branches: 
          - master
        if_ci_failed: error #success, failure, error, ignore
        only_pulls: false
        flags: 
          - frontend
        paths: 
          - src/frontend
        removed_code_behavior: off #removals_only, fully_covered_patch, adjust_base",
```

***

## parsers

```yaml
parsers:
  javascript:
    enable_partials: yes # default yes
  v1:
    include_full_missed_files: true  # default false
  gcov:
    branch_detection:
      conditional: yes
      loop: yes
      method: no
      macro: no
  cobertura:
  	handle_missing_conditions: true #default false
    partials_as_hits: true #default false
```

***

#### parsers.jacoco

Controls how Codecov counts partial coverage hits in Jacoco

* Type: map

```yaml YAML
parsers:
  jacoco:
    partials_as_hits: true #false by default
```

***

#### parsers.javascript

Unknown

* Type: map
* More Info: [https://docs.codecov.io/docs/node](https://docs.codecov.io/docs/node)

***

#### parsers.v1

Unknown

* Type: map
* More Info: [https://docs.codecov.io/docs/ruby](https://docs.codecov.io/docs/ruby)

***

#### parsers.gcov

Controls how Codecov parses branch coverage in gcov reports. The coverage must exist in the uploaded report for us to parse it.

* Type: map
* More Info: [https://docs.codecov.io/docs/codecov-yaml#section-default-yaml](https://docs.codecov.io/docs/codecov-yaml#section-default-yaml)

```yaml Default settings
parsers:
  gcov:
    branch_detection:
      conditional: yes
      loop: yes
      method: no
      macro: no
```

***

#### parsers.go

Controls how Codecov parses partial coverage in Golang reports. The coverage must exist in the uploaded report for us to parse it.

* Type: map

```yaml
parsers:
  go:
    partials_as_hits: true #false by default
```

***

#### parsers.cobertura

Controls how Codecov parses partial coverage in Cobertura reports. The coverage must exist in the uploaded report for us to parse it.

* Type: map

```yaml
parsers:
  cobertura:
    partials_as_hits: true # false by default
    handle_missing_conditons: true # false by default
```

<br />

#### parsers.lcov

Controls how Codecov parses partial coverage in lcov reports. The coverage must exist in the uploaded report for us to parse it.

```yaml
parsers:
  lcov:
    partials_as_hits: true #false by default
```

***

## ignore

* More info: [https://docs.codecov.io/docs/ignoring-paths](https://docs.codecov.io/docs/ignoring-paths)

```yaml
ignore:
  - "path/to/folder"  # ignore folders and all its contents
  - "test_*.rb"       # wildcards accepted
  - "**/*.py"         # glob accepted
```

***

## fixes

* More info: [https://docs.codecov.io/docs/fixing-paths](https://docs.codecov.io/docs/fixing-paths)

```yaml
fixes:
  - "before/::after/"  # move path   e.g., "before/path" => "after/path"
  - "::after/"         # move root   e.g., "path/" => "after/path/"
  - "before/::"        # reduce root e.g., "before/path/" => "path/"
```

***

## flags

* More info: [https://docs.codecov.io/docs/flags](https://docs.codecov.io/docs/flags)

```yaml
flags:
  projectA:  
    paths:
      - projectA/src
     carryforward: false #default -- false
  projectB:
    paths:
      - projectB/src
     carryforward: true
```

***

## comment

* More info: [https://docs.codecov.io/docs/pull-request-comments](https://docs.codecov.io/docs/pull-request-comments)

```yaml Default settings
comment:
  layout: "diff, flags, files"
  behavior: default
  require_changes: false  # if true: only post the comment if coverage changes
```

***

## slack\_app

Toggle the sending of notifications through the Codecov Slack App

* Type: boolean
* More info: [https://docs.codecov.com/docs/slack-integration](https://docs.codecov.com/docs/slack-integration)

```yaml
slack_app: false
```

***

## turn off github\_checks annotations (GitHub users only)

```yaml
github_checks:
  annotations: false
```

***

## Codecov Cloud Report Archiving

* How to disable archiving - In your codecov.yml, set the following line to false as follows:

```yaml
codecov:
  archive:
    uploads: false
```

### annotations

Specify whether to use GitHub Checks annotations or normal statuses. GitHub Checks are enabled by default and can be modified on the top level in your codecov.yml.

* Type: Boolean
* Default: true
* More Info: [GitHub Checks](https://docs.codecov.com/docs/github-checks)
