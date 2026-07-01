# Changelog

## [1.3.0](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.2.2...v1.3.0) (2026-07-01)


### Features

* Exit if seeding fails ([f9aa66f](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/f9aa66f42022b182b3af03a368a061303a910848))
* Fail-fast exit with code 1 when error propagates to Seed call ([5f59470](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/5f59470eb6dab390118413c5c7207ff769d3a6cd))
* Return error when JSON source file cannot be found. ([9319574](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/9319574082f0f3104478522b81a7b659359aad0a))
* Stop parsing virtual sources as files ([7e7f8fa](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/7e7f8fa513525064958290a937b270404fb56101))
* Tighten conversion safety for nullable integers ([6a5640c](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/6a5640c6d7f81d004e57bb978a11afb6bd5f6d2a))
* Wrap insert phase into a single transaction ([611d817](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/611d81770610ee2de7a30bc78bfb06a4819848d3))


### Bug Fixes

* Add go user before using it ([e3d53aa](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/e3d53aa12b21c74235ed90653eb1df66884606e0))
* Rename dBField to dbField to follow linting standards ([8115de2](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/8115de2b47500d9cdccc0d88ea588ee0740c143a))
* Use root:root ownership for copied artifacts and remove write bits for non-root users ([426af63](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/426af63a1c8f25b534b405d720994ff5e888cfab))

## [1.2.2](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.2.1...v1.2.2) (2026-03-07)


### Bug Fixes

* Image user set to non root user ([7d42706](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/7d42706563322b57ba169b497b6b9df0ae254123))

## [1.2.1](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.2.0...v1.2.1) (2026-03-04)


### Bug Fixes

* Add relics table ([855ba96](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/855ba9625fa7397f3b276c5cbf3ed311b5c0071a))
* Forgot to seed material_stats ([071ddca](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/071ddcab96e551844ba113d08cf3a434e64808fb))
* Massively refactored dbseeder ([573fdb1](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/573fdb1db7d69acfd8312658571fc11330debf2c))
* material stats ([b613558](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/b613558d656cf9d7d26f01078162247dd21f75c7))

## [1.2.0](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.1.2...v1.2.0) (2026-03-02)


### Features

* DB Seeder runs when db is empty ([6520792](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/65207925c759784b73888bc1edf957badeb923a0))


### Bug Fixes

* Add possibility to supply address through env ([706118f](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/706118f3489eea96f81f0d9c67951939fd1dd270))
* materials api ([cc3fbb9](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/cc3fbb9a936921c75a21a9b4e7de446c988de077))
* Remove default port expose ([3bc9bd7](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/3bc9bd78a76171b5b5dd573067da7cb1d159803a))
* Remove error log when shutting down gracefully ([9dfe585](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/9dfe585921d44a0de48e6adf96a42731ea2f2fff))

## [1.1.2](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.1.1...v1.1.2) (2026-03-01)


### Bug Fixes

* remove dbseeder from release for now ([ab55004](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/ab55004d80df752ade2b90bc6759f84e71e0e4bf))
* remove docker-compose ([948b531](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/948b531017539c1a1b01fe744b8af021142fc7af))

## [1.1.1](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.1.0...v1.1.1) (2026-03-01)


### Bug Fixes

* Correct docker repo ([19f5ec7](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/19f5ec79561c727b64ee85e07e51e2c00b1014cb))

## [1.1.0](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.0.1...v1.1.0) (2026-03-01)


### Features

* Add build and push image ([2becfa6](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/2becfa6ef40fdd236faf6830ad4b07e9ebb5a26e))

## [1.0.1](https://github.com/Ranoth/siralim-ultimate-team-builder-db/compare/v1.0.0...v1.0.1) (2026-03-01)


### Bug Fixes

* fix autorelease branch ([caaae6d](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/caaae6dfc183d5b258b62692439c321150821813))

## 1.0.0 (2026-03-01)


### Features

* API is DONE ([1778778](https://github.com/Ranoth/siralim-ultimate-team-builder-db/commit/1778778c73fc6cb5fb44fe5506a54995a544d8f1))
