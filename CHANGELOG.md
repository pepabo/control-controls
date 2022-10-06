## [v0.6.3](https://github.com/pepabo/control-controls/compare/v0.6.2...v0.6.3) - 2022-10-06
- Expand env when load YAML by @k1LoW in https://github.com/pepabo/control-controls/pull/20
- Fix defaultTemplate by @k1LoW in https://github.com/pepabo/control-controls/pull/22

## [v0.6.2](https://github.com/pepabo/control-controls/compare/v0.6.1...v0.6.2) - 2022-10-05
- Remove homebrew-tap setting because updates in the homebrew-tap repository by @k1LoW in https://github.com/pepabo/control-controls/pull/15
- Use tagpr by @k1LoW in https://github.com/pepabo/control-controls/pull/16
- Support notification by @k1LoW in https://github.com/pepabo/control-controls/pull/18
- Bump up go version by @k1LoW in https://github.com/pepabo/control-controls/pull/19

## [v0.6.1](https://github.com/pepabo/control-controls/compare/v0.6.0...v0.6.1) (2022-06-17)

* Fix handling non region arn (eg. `arn:aws:s3:::` ) [#14](https://github.com/pepabo/control-controls/pull/14) ([k1LoW](https://github.com/k1LoW))

## [v0.6.0](https://github.com/pepabo/control-controls/compare/v0.5.0...v0.6.0) (2022-06-16)

* Support workflow status (and note) management [#13](https://github.com/pepabo/control-controls/pull/13) ([k1LoW](https://github.com/k1LoW))

## [v0.5.0](https://github.com/pepabo/control-controls/compare/v0.4.0...v0.5.0) (2022-06-09)

* Add Validate() [#12](https://github.com/pepabo/control-controls/pull/12) ([k1LoW](https://github.com/k1LoW))

## [v0.4.0](https://github.com/pepabo/control-controls/compare/v0.3.0...v0.4.0) (2022-06-08)

* Add `--overlay` option for patch [#11](https://github.com/pepabo/control-controls/pull/11) ([k1LoW](https://github.com/k1LoW))

## [v0.3.0](https://github.com/pepabo/control-controls/compare/v0.2.1...v0.3.0) (2022-06-07)

* Fix flag [#10](https://github.com/pepabo/control-controls/pull/10) ([k1LoW](https://github.com/k1LoW))
* Add reason of disabled in the configuration file. [#9](https://github.com/pepabo/control-controls/pull/9) ([k1LoW](https://github.com/k1LoW))

## [v0.2.1](https://github.com/pepabo/control-controls/compare/v0.2.0...v0.2.1) (2022-04-18)

* Fix nil pointer dereference [#8](https://github.com/pepabo/control-controls/pull/8) ([k1LoW](https://github.com/k1LoW))

## [v0.2.0](https://github.com/pepabo/control-controls/compare/v0.1.2...v0.2.0) (2022-04-18)

* exit status 2 when plan diff is not empty [#7](https://github.com/pepabo/control-controls/pull/7) ([k1LoW](https://github.com/k1LoW))

## [v0.1.2](https://github.com/pepabo/control-controls/compare/v0.1.1...v0.1.2) (2022-04-15)

* Fix contextcopy bug [#6](https://github.com/pepabo/control-controls/pull/6) ([k1LoW](https://github.com/k1LoW))

## [v0.1.1](https://github.com/pepabo/control-controls/compare/v0.1.0...v0.1.1) (2022-04-15)

* Fix sechub.Override behavior [#5](https://github.com/pepabo/control-controls/pull/5) ([k1LoW](https://github.com/k1LoW))

## [v0.1.0](https://github.com/pepabo/control-controls/compare/60006830255c...v0.1.0) (2022-04-14)

* Add option `--disabled-reason` [#4](https://github.com/pepabo/control-controls/pull/4) ([k1LoW](https://github.com/k1LoW))
* Add command `plan` [#3](https://github.com/pepabo/control-controls/pull/3) ([k1LoW](https://github.com/k1LoW))
* Fix apply [#2](https://github.com/pepabo/control-controls/pull/2) ([k1LoW](https://github.com/k1LoW))
* Add command `apply` [#1](https://github.com/pepabo/control-controls/pull/1) ([k1LoW](https://github.com/k1LoW))
