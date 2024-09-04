# COMBI (Config Combinator)

![GitHub Release](https://img.shields.io/github/v/release/freepik-company/combi)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/freepik-company/combi)
![GitHub License](https://img.shields.io/github/license/freepik-company/combi)

![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/freepik-company/combi/total)

![GitHub User's stars](https://img.shields.io/github/stars/sebastocorp)
![GitHub followers](https://img.shields.io/github/followers/sebastocorp)

## Description

Combi is a simple tool that consumes, update and merge diferents configurations in different formats to generate a final usable configuration and performs some defined actions at the end of the process. This tool consumes its own configuration from a source (like a file in a git repository), and perform a merge with the patches defined in that configuration in a local configuration file with the same format as the patches (libconfig, yaml, json, etc).

## Motivation

Many services used daily in the industry receive their configuration mostly from a file with a specific format (yaml, json, libconfig, etc.), and many of these services play a critical role, so downtime has to be minimal. or, if possible, it should be none.

For this reason, many of these services have functionalities to collect the updated configuration again without stopping its execution. Thanks to these functionalities, the configuration can be modified and applied without practically any downtime, but the problem appears, as always, when you have a multitude of instances of that service, with different configurations, in different environments.

The idea is to avoid restarting the service and easily update the configuration, so some problems arise here:

- In a container environment you do't want to rotate them so as not to generate downtime and lost requests, even if it is minimal and you can control it more or less with a rotation strategy.
- And whether with containers or vms, if you want to avoid the restart, you will have to modify the specific configuration they have, enter each of the instances to update it, and execute the functionality that refreshes the configuration.
- If you have different configurations of the same service in different instances, those configurations will often be splitted in different repositories, or if it is mono-repo, they will be splitted in different parts of the repository, or in different files.

It may also be the case that you want 2 parts of the same config separated from each other, and different teams modify one of these parts separately, one of the two configurations will be the base config, with optional fields, and the other it will be a chore config, with mandatory fields, and you want performing a merge of the configurations with precedence of the chore configuration.

Thinking about these problems and possible solutions, we have decided to create this tool, which is not only capable of centralizing the different configurations of the same format, but is also capable of performing checks on the final configuration and executing the commands necessary to refresh the configuration.

## Flags

| Name                  | Command  | Default                              | Description |
|:---                   |:---      |:---                                  |:---         |
| `log-level`           | `daemon` | `info`                               | Verbosity level for logs |
| `tmp-dir`             | `daemon` | `/tmp/combi`                         | Temporary directoty to store temporary objects like remote repos, scripts, etc |
| `sync-time`           | `daemon` | `15s`                                | Waiting time between source synchronizations (in duration type) |

## How to use

This project provides the binary files and container image in differents architectures to make it easy to use wherever wanted.

### Configuration

Current configuration version: `v1alpha2`

#### Root Parameters

| Name   | Default | Description |
|:---    |:---     |:---         |
| `kind` | `""`    | type of the configuration files to manage. The possible values ​​of this field can be found in the [supported configuration formats](#supported-configuration-formats) section |

#### Global Parameters

| Name                          | Default | Description |
|:---                           |:---     |:---         |
| `config.source`               | `""`    | source with the configuration to merge in the `mergedConfig` file |
| `global.conditions.mandatory` | `[]`    | list of mandatory conditions that `mergedConfig` have to achive |
| `global.conditions.optional`  | `[]`    | list of optional conditions that `mergedConfig` have to check |
| `global.actions.onSuccess`    | `[]`    | list of actions that `combi` has to execute in case of mandatory conditions success |
| `global.actions.onFailure`    | `[]`    | list of actions that `combi` has to execute in case of mandatory conditions fail |

#### Configs Parameters

| Name                                        | Default | Description |
|:---                                         |:---     |:---         |
| `config.mergedConfig`         | `""`    | filepath to create the configuration file after the merge and conditions |
| `config.source`               | `""`    | source with the configuration to merge in the `mergedConfig` file |
| `config.conditions.mandatory` | `[]`    | list of mandatory conditions that `mergedConfig` have to achive |
| `config.conditions.optional`  | `[]`    | list of optional conditions that `mergedConfig` have to check |
| `config.actions.onSuccess`    | `[]`    | list of actions that `combi` has to execute in case of mandatory conditions success |
| `config.actions.onFailure`    | `[]`    | list of actions that `combi` has to execute in case of mandatory conditions fail |

#### Condition List Item Parameters

| Name       | Default | Description |
|:---        |:---     |:---         |
| `name`     | `""`    | name of the current condition to evaluate |
| `template` | `""`    | golang template to compile and extract some config value to compare with condition `value` |
| `value`    | `""`    | value to compare with the result of `template` |

#### Action List Item Parameters

| Name       | Default | Description |
|:---        |:---     |:---         |
| `name`     | `""`    | name of the current action to execute |
| `command`  | `[]`    | string list with the command and his argument to execute |
| `script`   | `""`    | string with a script that combi generate to execute |

> **WARNING**
>
> The list of actions are commands that are executed on the machine after checking the conditions. Please be careful.

### Sources

You can consume the configuration from different sources.

| Name    | Description |
|:---     |:---         |
| `config.source.type`                 | type of the source to define where consume the configuration (`raw`, `filepath`, `git` y `kubernetes`) |
| `config.source.raw`                  | direct configuration string to merge |
| `config.source.filepath`             | file with the configuration to merge |
| `config.source.git`                  | git repository to find the file with the configuration to merge |
| `config.source.git.sshUrl`           | ssh url of the git repository |
| `config.source.git.sshKeyFilepath`   | ssh private key file to use in git clone |
| `config.source.git.branch`           | branch of the git repository |
| `config.source.git.filepath`         | filepath inside of repository of the config file |
| `config.source.kubernetes`           | kubernetes resource `ConfigMap` or `Secret` with the configuration to merge |
| `config.source.kubernetes.kind`      | kind of the resource (`ConfigMap` or `Secret`) |
| `config.source.kubernetes.namespace` | namespace of the resource |
| `config.source.kubernetes.name`      | name of the resource |
| `config.source.kubernetes.key`       | key inside of the resource to find the configuration to merge |

## How does it work?

Synchronization flow diagram:

```sh
               ┌─────────────┐                                   
               │             │                                   
 ┌────┬────┬───►  sync time  │                                   
 │    │    │   │             │                                   
 │    │    │   └──────┬──────┘                                   
 │    │    │          │                                          
 │    │    │  ┌───────▼───────┐    ┌──────────┐                  
 │    │    │  │               │    │          │  │local file     
 │    │    no │  get  config  ◄────┤  source  ├─►│git repo       
 │    │    │  │               │    │          │  │...            
 │    │    │  └───────┬───────┘    └──────────┘                  
 │    │    │          │                                          
 │    │    │    ┌─────▼─────┐                                    
 │    │    │    │           │                                    
 │    │    └────┤  update?  │                                    
 │    │         │           │                                    
 │    │         └─────┬─────┘                                    
 │    │               │                                          
 │    │              yes                                         
 │    │               │                                          
 │    │      ┌────────▼────────┐                                 
 │    │      │                 │                                 
 │    │      │  decode config  ◄─────────────┐                   
 │    │      │                 │             │                   
 │    │      └────────┬────────┘             │                   
 │    │               │                      │                   
 │    │      ┌────────▼────────┐             │                   
 │    │      │                 │             │                   
 │    │      │  merge configs  │             │                   
 │    │      │                 │       ┌─────┴─────┐  │libconfig 
 │    │      └────────┬────────┘       │           │  │nginx conf
 │    │               │                │  encoder  ├─►│json      
 │    │    ┌──────────▼──────────┐     │           │  │yaml      
 │    │    │                     │     └─────┬─────┘  │...       
 │    │    │  check  conditions  │           │                   
 │    │    │                     │           │                   
 │    │    └──────────┬──────────┘           │                   
 │    │               │                      │                   
 │    │     ┌─────────┴─────────┐            │                   
 │    │     │                   │            │                   
 │ ┌──┴─────▼────────┐ ┌────────▼─────────┐  │                   
 │ │                 │ │                  │  │                   
 │ │  fail acttions  │ │  encode  config  ◄──┘                   
 │ │                 │ │                  │                      
 │ └─────────────────┘ └────────┬─────────┘                      
 │                              │                                
┌┴───────────────────┐ ┌────────▼─────────┐                      
│                    │ │                  │                      
│  success acttions  ◄─┤  update  config  │                      
│                    │ │                  │                      
└────────────────────┘ └──────────────────┘                      
```

## Example

To consume a configuration in a git repository (with specific branch) and merge in a local config with `libconfig` format:

```sh
combi daemon \
    --sync-time=20s \
    --tmp-dir=/tmp/combi \
    --config=/config/combi.yaml
```

The combi config.yaml:

```yaml
# combi configuration file
kind: libconfig
global:
  source:
    type: raw
    raw: |
      int32field=2
      int64field=500L

      string_example="some-value"

      group_example=
      {
        admin_credentials="root:pass"
        addr="0.0.0.0:6032"
      }
  conditions:
    mandatory: []
    optional: []
  actions:
    onSuccess: []
    onFailure: []
configs:
  mergedConfig: ./path/to/merged/libconfig.cnf
  source:
    type: raw
    raw: |
      mysql_variables=
      {
        threads=2
        max_connections=500
      }

      mysql_servers =
      (
        { address="127.0.0.1" , port=3306 , hostgroup=0 , max_connections=1000, weight=1 },
        { address="127.0.0.2" , port=3306 , hostgroup=1 , max_connections=1000, weight=1 },
      )
      
      mysql_users:
      (
        { username = "writer" , password = "pass" , default_hostgroup = 0 , active = 1 },
        { username = "reader" , password = "pass" , default_hostgroup = 1 , active = 1 },
      )

  conditions:
    mandatory:
    - name: "search primitive value to check condition"
      template: |
        {{- $config := . -}}
        {{- printf "%s" $config.int64field -}}
      value: "500L"
    - name: "search group value to check condition"
      template: |
        {{- $config := . -}}
        {{- printf "%s" $config.mysql_variables.threads -}}
      value: "2"
    - name: "search list value to check condition"
      template: |
        {{- $config := . -}}
          {{- range $i, $v := $config.mysql_servers -}}
              {{- if (eq $v.hostgroup "0" ) -}}
                {{- printf "%s" $v.max_connections -}}
              {{- end -}}
          {{- end -}}
      value: "1000"
    - name: "search env variable to check condition"
      template: |
        {{- printf "%s" (env "MANDATORY_ENV_VAR") -}}
      value: "true"
    optional: []

  actions:
    onSuccess:
    - name: "execute success message config action"
      command:
      - echo
      - -e
      - "success in config for you\n"
    onFailure:
    - name: "execute success message config action"
      command:
      - echo
      - -e
      - "fail in config for you\n"
```

## Supported configuration formats

| Format      | Status |
|:---         |:---    |
| `json`      | ✅     |
| `nginx`     | ✅     |
| `libconfig` | ✅     |
| `yaml`      | ❌     |
| `hcl`       | ❌     |

## How to collaborate

We are open to external collaborations for this project: improvements, bugfixes, whatever.

For doing it, open an issue to discuss the need of the changes, then:

- Fork the repository
- Make your changes to the code
- Open a PR and wait for review
