<!-- markdownlint-disable-file MD041 -->
## upcoming release

* switch to the standalone SDK v2 for compatibility with last Terraform version
* move docs in dedicated directory
* bump golang version
* refactor release workflow (GH Actions) to generate files compatible with Terraform registry
* default `auth_pass` for `lvsnetwork_ifacevrrp` is now configurable on provider config with `default_auth_pass`

## 1.1.2 (July 6, 2021)

* fix permanent conflict between `vault_enable` and `login`/`password` provider arguments

## 1.1.1 (December 18, 2020)

* fix track_script not computed if there isn't vip

## 1.1.0 (October 19, 2019)

* add vrrp_script resource and option for ifacevrrp

## 1.0.0 (October 10, 2019)

First release
