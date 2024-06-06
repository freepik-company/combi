
# TODO lists

## Features

- [ ] Add patch system post merge
- [ ] Consume combi config from k8s ConfigMap/Secrets
- [ ] Consume target config from git repository file ConfigMap/Secrets
- [ ] Consume target config from ConfigMap/Secrets field

## Supports

- [x] Add support to nginx conf files with custom parser
- [ ] Add support to hcl files
- [x] Add support to json files with golang parser
- [ ] Add support to yaml files with golang parser

## Code

- [ ] Change merge order beetween global and specific configs
- [ ] Add source and target structure to consume configs from different sources
- [ ] Add type DaemonT with attached flow functions
- [ ] Refactor the code to clean it and add comments
- [ ] Change to custom parser in libconfig kind instead use library
