# Changelog
Changelog for NeoFS Node

## [Unreleased]

## [0.27.5] - 2022-01-31

### Fixed
- Flush small objects when persist write cache (#1088)
- Empty response body in object.Search request (#1098)
- Inner ring correctly checks session token in container.SetEACL request (#1110)
- Printing in verbose mode in CLI (#1120)
- Subnet removal event processing (#1123)

### Added
- Password support in CLI config (#1103)
- Shard dump restore commands in CLI (#1085, #1086)
- `acl extended create` command in CLI (#1092)

### Changed
- Adopt new `owner.ID` API from SDK (#1100)
- Use `go install` instead of `go get` in Makefile (#1102)
- Storage node returns Fixed12 decimal on accounting.Balance request. CLI
  prints Fixed8 rounded value by default. (#1084)
- Support new update interface for NNS contract in NeoFS Adm (#1091)
- Rename `use_write_cache` to `writecache.enabled` in stoarge config (#1117)
- Preallocate slice in `headersFromObject` (#1115)
- Unify collection of expired objects (#1115)
- Calculate blobovnicza size at initialization properly (#1115)
- Process fast search filters outside bbolt transaction (#1115)
- Update TZHash library to v1.5.1

### Removed
- `--wif` and `--binary-key` keys from CLI (#1083)
- Extended ACL validator moved to SDK library (#1096)
- `--generate-key` flag in CLI control commands (#1103)
- Various unused code (#1123)

### Upgrading from v0.27.4
Use `--wallet` key in CLI to provide WIF or binary key file instead of `--wif`
and `--binary-key`.

Replace `NEOFS_STORAGE_SHARD_N_USE_WRITE_CACHE` with 
`NEOFS_STORAGE_SHARD_N_WRITECACHE_ENABLED` in Storage node config.

Specify `password: xxx` in config file for NeoFS CLI to avoid password input.

## [0.27.4] - 2022-01-13

### Fixed
- ACL check did not produce status code (#1062)
- Asset transfer wrapper used incorrect receiver (#1069)
- Empty search response missed meta header and body (#1063)
- IR node in single chain environment used incorrect source of IR list (#1025)
- Incorrect message sequence in object.Range request (#1077)

### Added
- Option to disable compression of object based on their content-type attribute
  (#1060)

### Changed
- Factor out autocomplete command in CLI and Adm (#1041)
- Single crypto rand source (#851)

### Upgrading from v0.27.3
To disable compression for object with specific content-types, specify them
as a string array in blobstor section:
`NEOFS_STORAGE_SHARD_N_BLOBSTOR_COMPRESSION_EXCLUDE_CONTENT_TYPES`. Use
asterisk as wildcard, e.g. `video/*`.

## [0.27.3] - 2021-12-30

### Added
- `SetShardMode` RPC in control API, available in CLI (#1044)
- Support of basic ACL constants without final flag in CLI (#1066)

### Changed
- `neofs-adm` updates contracts in single tx (#1035)
- Proxy contract arguments for deployment in `neofs-adm` (#1056)

## [0.27.2] - 2021-12-28

### Fixed
- Goroutine leak due to infinite response message await ([neofs-api-go#366](https://github.com/nspcc-dev/neofs-api-go/pull/366))
- Inconsistency in placement function ([neofs-sdk-go#108](https://github.com/nspcc-dev/neofs-sdk-go/pull/108))

### Added
- `ListShards` RPC in control API, available in CLI (#1043)
- Epoch metric in Storage and Inner Ring applications (#1054)

### Changed
- Some object replication related logs were moved to DEBUG level (#1052)

## [0.27.1] - 2021-12-20

### Fixed
- Big objects now flushed from WriteCache after write (#1028)
- WriteCache big object counter (#1022)
- Panic in the container estimation routing (#1016)
- Shutdown freeze in policer component (#1047)

### Added
- Shorthand `-g` for `--generate-key` in NeoFS CLI (#1034)
- Autocomplete generator command for neofs-adm (#1013)
- Max connection per host config value for neo-go client (#780)
- Sanity check of session token context in container service (#1045)

### Changed
- CLI now checks NeoFS status code for exit code (#1039)
- New `Update` method signature for NNS contract in neofs-adm (#1038)

## [0.27.0] - 2021-12-09 - Sinjido (신지도, 薪智島)

NeoFS API v2.11.0 support with response status codes and storage subnetworks.

### Fixed
- CLI now opens LOCODE database in read-only mode for listing command (#958)
- Tombstone owner now is always set (#842)
- Node in relay mode does not require shard config anymore (#969)
- Alphabet nodes now ignore notary notifications with non-HALT main tx (#976) 
- neofs-adm now prints version of NNS contract (#1014)
- Possible NPE in blobovnicza (#1007)
- More precise calculation of blobovnicza size (#915)

### Added
- Maintenance mode for Storage node (#922)
- Float values in Storage node config (#903)
- Status codes for NeoFS API Response messages (#961)
- Subnetwork support (#977, #973, #983, #974, #982, #979, #998, #995, #1001, #1004)
- Customized fee for named container registration (#1008)

### Changed
- Alphabet contract number is not mandatory (#880)
- Alphabet nodes resign `AddPeer` request if it updates Storage node info (#938)
- All applications now use client from neofs-sdk-go library (#966)
- Some shard configuration records were renamed, see upgrading section (#859)
- `Nonce` and `VUB` values of notary transactions generated from notification 
  hash (#844)
- Non alphabet notary invocations now have 4 witnesses (#975)
- Object replication is now async and continuous (#965)
- NeoFS ADM updated for the neofs-contract v0.13.0 deploy (#984)
- Minimal TLS version is set to v1.2 (#878)
- Alphabet nodes now invoke `netmap.Register` to add node to the network map 
  candidates in notary enabled environment (#1008)

### Upgrading from v0.26.1
`NEOFS_IR_CONTRACTS_ALPHABET_AMOUNT` is not mandatory env anymore. If it
is not set, Inner Ring would try to read maximum from config and NNS contract.
However, that parameter still can be set in order to require the exact number
of contracts.

Shard configuration records were renamed:
- `refill_metabase` -> `resync_metabase`
- `writecache.max_size` -> `writecache.max_object_size`
- `writecache.mem_size` -> `writecache.memcache_capacity`
- `writecache.size_limit` -> `writecache_capcity`
- `blobstor.blobovnicza.opened_cache_size` -> `blobstor.blobovnicza.opened_cache_capacity`
- `*.shallow_depth` -> `*.depth`
- `*.shallow_width` -> `*.width`
- `*.small_size_limit` -> `*.small_object_size`

Specify storage subnetworks in `NEOFS_NODE_SUBNET_ENTRIES` as the list of
integer numbers. To exit default subnet, use `NEOFS_NODE_SUBNET_EXIT_ZERO=true`

Specify fee for named container registration in notary disabled environment
with `NEOFS_IR_FEE_NAMED_CONTAINER_REGISTER`.

## [0.26.1] - 2021-11-02

### Fixed
- Storage Node handles requests before its initialization is finished (#934)
- Release worker pools gracefully (#901)
- Metabase ignored containers of storage group and tombstone objects 
  in listing (#945)
- CLI missed endpoint flag in `control netmap-snapshot` command (#942)
- Write cache object persisting (#866)

### Added
- Quote symbol support in `.env` example tests (#935)
- FSTree object counter (#821)
- neofs-adm prints contract version in `dump-hashes` command (#940)
- Default values section in shard configuration (#877)
- neofs-adm downloads contracts directly from GitHub (#733)

### Changed
- Use FSTree counter in write cache (#821)
- Calculate notary deposit `till` parameter depending on available 
  deposit (#910)
- Storage node returns session token error if attached token's private key 
  is not available (#943)
- Refactor of NeoFS API client in inner ring (#946)
- LOCODE generator tries to find the closest continent if there are 
  no exact match (#955)

### Upgrading from v0.26.0
You can specify default section in storage engine configuration. 
See [example](./config/example/node.yaml) for more details.

## [0.26.0] - 2021-10-19 - Udo (우도, 牛島)

NeoFS API v2.10 support

### Fixed
- Check remote node public key in every response message (#645)
- Do not lose local container size estimations (#872)
- Compressed and uncompressed objects are always available for reading 
  regardless of compression configuration (#868)
- Use request session token in ACL check of object.Put (#881)
- Parse URI in neofs-cli properly (#883)
- Parse minutes in LOCODE DB properly (#902)
- Remove expired tombstones (#884)
- Close all opened blobovniczas properly (#896)
- Do not accept objects with empty OwnerID field (#841)

### Added
- More logs in governance and policer components (#867, #882)
- Contract address getter in static blockchain clients (#627)
- Alphabet configuration option to disable governance sync (#869)
- neofs-lens app implementation (#791)
- Detailed comments in neofs-node config example (#858)
- Size suffixes support in neofs-node config (#857)
- Docs for neofs-adm (#906)
- Side chain block size duration and global NeoFS configuration in
  NetworkConfig response (#833)
- Support native container names (#889)

### Changed
- Updated grpc to v1.41.0 (#860)
- Updated neo-go to v0.97.3 (#833)
- Updated neofs-api-go to v1.30.0
- Adopt neofs-adm for new contracts release (#835, #888)
- Adopt neofs-node for new contracts release (#905)
- SN and IR notary deposits are made dynamically depending on the Notary and
  GAS balances (#771)
- VMagent port in testnet config is now 443 (#908)
- Use per-shard worker pools for object.Put operations (#674)
- Renamed `--rpc-endpoint` CLI flag for `control command` to `--endpoint` (#879)

### Removed
- Global flags in CLI. Deleted useless flags from `accounting balance`
  command (#810).
- Interactive mode in docker run command (#916)

### Upgrading from v0.25.1
Deleted `NEOFS_IR_NOTARY_SIDE_DEPOSIT_AMOUNT`, `NEOFS_IR_NOTARY_MAIN_DEPOSIT_AMOUNT`
and `NEOFS_IR_TIMERS_SIDE_NOTARY`, `NEOFS_IR_TIMERS_MAIN_NOTARY` Inner Ring envs.
Deleted `NEOFS_MORPH_NOTARY_DEPOSIT_AMOUNT` and `NEOFS_MORPH_NOTARY_DEPOSIT_DURATION`
Storage Node envs.
`control` CLI command does not have `--rpc-endpoint`/`r` flag, use `endpoint`
instead.

## [0.25.1] - 2021-09-29

### Fixed
- Panic caused by missing Neo RPC endpoints in storage node's config (#863)

### Added
- Support of multiple Neo RPC endpoints in Inner Ring node (#792)

`mainchain` section of storage node config is left unused by the application. 

## [0.25.0] - 2021-09-27 - Mungapdo (문갑도, 文甲島)

### Fixed
- Work of a storage node with one Neo RPC endpoint instead of a list (#746)
- Lack of support for HEAD operation on the object write cache (#762)
- Storage node attribute parsing is stable now (#787)
- Inner Ring node now logs transaction hashes of Deposit and Withdraw events 
  in LittleEndian encoding (#794)
- Storage node uses public keys of the remote nodes in placement traverser
  checks (#645)
- Extended ACL `Target` check of role and public keys is mutual exclusive now
  (#816)
- neofs-adm supports update and deploy of neofs-contract v0.11.0 (#834, #836)
- Possible NPE in public key conversion (#848)
- Object assembly routine do not forward existing request instead of creating 
  new one (#839)
- Shard now returns only physical stored objects for replication (#840)

### Added
- Support events from P2P notary pool
- Smart contract address auto negotiation with NNS contract (#736)
- Detailed logs for all data writing operations in storage engine (#790)
- Docker build and release targets in Makefile (#785)
- Metabase restore option in the shard config (#789) 
- Write cache used size limit in bytes (#776)

### Changed
- Reduce container creation delay via listening P2P notary pool (#519)
- Extended ACL table is not limited to 1KiB (#731)
- Netmap side chain client wrapper now has `TryNotary` option (#793)
- Sticky bit is ignored in requests with `SYSTEM` role (#818)
- Incomplete object put error now contains last RPC error (#778)
- Container service invalidates container cache on writing operations (#803)
- Improved write cache size counters (#776)
- Metabase returns `NotFound` error instead of `AlreadyRemoved` on GCMarked
  objects (#840)
- Object service uses two routine pools for remote and local GET requests (#845)

### Removed
- Dockerfile for AllInOne image moved to a separate repository (#796)

### Upgrading from v0.24.1
Added `NEOFS_CONTRACTS_PROXY` env for Storage Node; mandatory in
notary enabled environments only. It should contain proxy contract's
scripthash in side chain.

Added `NEOFS_MORPH_NOTARY_DEPOSIT_AMOUNT` and
`NEOFS_MORPH_NOTARY_DEPOSIT_DURATION` envs for Storage Node, that
have default values, not required. They should contain notary deposit
amount and frequency(in blocks) respectively.

All side chain contract address config values are optional. If side chain
contract address is not specified, then value gathered from NNS contract.

Added `NEOFS_STORAGE_SHARD_<N>_WRITECACHE_SIZE_LIMIT` where `<N>` is shard ID.
This is the size limit for the all write cache storages combined in bytes. Default
size limit is 1 GiB.

Added `NEOFS_STORAGE_SHARD_<N>_REFILL_METABASE` bool flag where `<N>` is shard 
ID. This flag purges metabase instance at the application start and reinitialize
it with available objects from the blobstor.

Object service pool size now split into `NEOFS_OBJECT_PUT_POOL_SIZE_REMOTE` and
`NEOFS_OBJECT_PUT_POOL_SIZE_LOCAL` configuration records.

## [0.24.1] - 2021-09-07

### Fixed 
- Storage and Inner Ring will not start until Neo RPC node will have the height
of the latest processed block by the nodes (#795)

### Upgrading from v0.24.0
Specify path to the local state DB in Inner Ring node config with 
`NEOFS_IR_NODE_PERSISTENT_STATE_PATH`. Specify path to the local state DB in
Storage node config with `NEOFS_NODE_PERSISTENT_STATE_PATH`.

## [0.24.0] - 2021-08-30 Anmyeondo (안면도, 安眠島)

### Fixed
- Linter warning messages (#766)
- Storage Node does not register itself in network in relay mode now (#761)

### Changed
- `neofs-adm` fails when is called in a notary-disabled environment (#757)
- `neofs-adm` uses `neo-go` client's native NNS resolving method instead of the custom one (#756)
- Node selects pseudo-random list of objects from metabase for replication (#715)

### Added
- Contract update support in `neofs-adm` utility (#748)
- Container transferring support in `neofs-adm` utility (#755)
- Storage Node's balance refilling support in `neofs-adm` utility (#758)
- Support `COMMON_PREFIX` filter for object attributes in storage engine and `neofs-cli` (#760) 
- Node's and IR's notary status debug message on startup (#758)
- Go `1.17` unit tests in CI (#766)
- Supporting all eACL filter fields from the specification (#768)
- Cache for Container service's read operations (#676)

### Updated
- `neofs-api-go` library to `v1.29.0`

### Removed
- Unused `DB_SIZE` parameter of writecache (#773)

### Upgrading from v0.23.1
Storage Node does not read unused `NEOFS_STORAGE_SHARD_XXX_WRITECACHE_DB_SIZE`
config parameter anymore.

## [0.23.1] - 2021-08-06

N3 Mainnet launch release with minor fixes.

### Added
- Initial version of `neofs-adm` tool for fast side chain deployment and
  management in private installations
- Notary support auto negotiation (#709)
- Option to disable side chain cache in Storage node (#704)
- Escape symbols in Storage node attributes (#700)

### Changed
- Updated neo-go to v0.97.1
- Updated multiaddr lib to v0.4.0 with native TLS protocol support (#661)
- Default file permission in storage engine is 660 (#646)

### Fixed
- Container size estimation routine now aggregates values by cid-epoch tuple
  (#723)
- Storage engine always creates executable dirs (#646)
- GC routines in storage engine shards shutdown gracefully (#745)
- Handle context shutdown at NeoFS multi client group address switching (#737)
- Scope for main chain invocations from Inner Ring nodes (#751)

### Upgrading from v0.23.0
Added `NEOFS_MORPH_DISABLE_CACHE` env. If `true`, none of
the `eACL`/`netmap`/`container` RPC responses cached.

Remove `WITHOUT_NOTARY` and `WITHOUT_MAIN_NOTARY` records from Inner Ring node
config. Notary support is now auto negotiated.

## [0.23.0] - 2021-07-23 - Wando (완도, 莞島)

Improved stability for notary disabled environment.

### Added
- Alphabet wallets generation command in neofs-adm (#684)
- Initial epoch timer tick synchronization at Inner Ring node startup (#679)

### Changed
- `--address` flag is optional in NeoFS CLI (#656)
- Notary subsystem now logs `ValidUntilBlock` (#677)
- Updated neo-go to v0.96.1
- Storage Node configuration example contains usable parameters (#699)

### Fixed
- Do not use side chain RoleManagement contract as source of Inner Ring list 
  when notary disabled in side chain (#672)
- Alphabet list transition is even more effective (#697)
- Inner Ring node does not require proxy and processing contracts if notary
  disabled (#701, #714)

### Upgrading from v0.22.3
To upgrade Storage node or Inner Ring node from v0.22.3, you don't need to
change configuration files. Make sure, that NEO RPC nodes, specified in config,
are connected to N3 RC4 (Testnet) network.

## [0.22.3] - 2021-07-13

### Added
- Support binary eACL format in container CLI command ([#650](https://github.com/nspcc-dev/neofs-node/issues/650)).
- Dockerfile for neofs-adm utility ([#680](https://github.com/nspcc-dev/neofs-node/pull/680)).

### Changed
- All docker files moved to `.docker` dir ([#682](https://github.com/nspcc-dev/neofs-node/pull/682)).

### Fixed
- Do not require MainNet attributes in "Without MainNet" mode ([#663](https://github.com/nspcc-dev/neofs-node/issues/663)).
- Stable alphabet list merge in Inner Ring governance ([#670](https://github.com/nspcc-dev/neofs-node/issues/670)).
- User can specify only wallet section without node key ([#690](https://github.com/nspcc-dev/neofs-node/pull/690)).
- Log keys in hex format in reputation errors ([#693](https://github.com/nspcc-dev/neofs-node/pull/693)).
- Connections leak and reduced amount of connection overall ([#692](https://github.com/nspcc-dev/neofs-node/issues/692)).

### Removed
- Debug output of public key in Inner Ring log ([#689](https://github.com/nspcc-dev/neofs-node/pull/689)).

## [0.22.2] - 2021-07-07

Updated broken version of NeoFS API Go.

### Updated
- NeoFS API Go: [v1.28.3](https://github.com/nspcc-dev/neofs-api-go/releases/tag/v1.28.3).

## [0.22.1] - 2021-07-07

### Added
- `GetCandidates` method to morph client wrapper ([#647](https://github.com/nspcc-dev/neofs-node/pull/647)).
- All-in-One Docker image that contains all NeoFS related binaries ([#662](https://github.com/nspcc-dev/neofs-node/pull/662)).
- `--version` flag to Storage Node binary ([#664](https://github.com/nspcc-dev/neofs-node/issues/664)).

### Changed
- Do not check NeoFS version in `LocalNodeInfo` requests and `Put` container operations; `v2.7.0` is genesis version of NeoFS ([#660](https://github.com/nspcc-dev/neofs-node/pull/660)).
- All error calls of CLI return `1` exit code ([#657](https://github.com/nspcc-dev/neofs-node/issues/657)).

### Fixed
- Do not use multisignature for audit operations ([#658](https://github.com/nspcc-dev/neofs-node/pull/658)).
- Skip audit for containers without Storage Groups ([#659](https://github.com/nspcc-dev/neofs-node/issues/659)).

### Updated
- NeoFS API Go: [v1.28.2](https://github.com/nspcc-dev/neofs-api-go/releases/tag/v1.28.2).

## [0.22.0] - 2021-06-29 - Muuido (무의도, 舞衣島)

Storage nodes with a group of network endpoints.

### Added
- Support of Neo wallet credentials in CLI ([#610](https://github.com/nspcc-dev/neofs-node/issues/610)).
- More reliable approval of trust value by IR ([#500](https://github.com/nspcc-dev/neofs-node/issues/500)).
- Storage node's ability to announce and serve on multiple network addresses ([#607](https://github.com/nspcc-dev/neofs-node/issues/607)).
- Validation of network addresses of netmap candidates in IR ([#557](https://github.com/nspcc-dev/neofs-node/issues/557)).
- Control service with healthcheck RPC in IR and CLI support ([#414](https://github.com/nspcc-dev/neofs-node/issues/414)).

### Fixed
- Approval of objects with with duplicate attribute keys or empty values ([#633](https://github.com/nspcc-dev/neofs-node/issues/633)). 
- Approval of containers with with duplicate attribute keys or empty values ([#634](https://github.com/nspcc-dev/neofs-node/issues/634)).
- Default path for CLI config ([#626](https://github.com/nspcc-dev/neofs-node/issues/626)). 

### Changed
- `version` command replaced with `--version` flag in CLI ([#571](https://github.com/nspcc-dev/neofs-node/issues/571)).
- Command usage text is not printed on errors in CLI ([#623](https://github.com/nspcc-dev/neofs-node/issues/623)).
- `netmap snapshot` command replaced with `control netmap-snapshot` one in CLI ([#651](https://github.com/nspcc-dev/neofs-node/issues/651)).
- IR does not include nodes with LOCODE derived attributes to the network map ([#412](https://github.com/nspcc-dev/neofs-node/issues/412)).
- IR uses morph/client packages for contract invocations ([#496](https://github.com/nspcc-dev/neofs-node/issues/496)).
- Writecache decreases local size when objects are flushed ([#568](https://github.com/nspcc-dev/neofs-node/issues/568)).
- IR can override global configuration values only in debug build ([#363](https://github.com/nspcc-dev/neofs-node/issues/363)).

### Updated
- Neo Go: [v0.95.3](https://github.com/nspcc-dev/neo-go/releases/tag/v0.95.3).
- NeoFS API Go: [v1.28.0](https://github.com/nspcc-dev/neofs-api-go/releases/tag/v1.28.0).
- protobuf: [v1.26.0](https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.26.0).
- uuid: [v1.2.0](https://github.com/google/uuid/releases/tag/v1.2.0).
- compress: [v1.13.1](https://github.com/klauspost/compress/releases/tag/v1.13.1).
- base58: [v1.2.0](https://github.com/mr-tron/base58/releases/tag/v1.2.0).
- multiaddr: [v0.3.2](https://github.com/multiformats/go-multiaddr/releases/tag/v0.3.2).
- ants: [v2.4.0](https://github.com/panjf2000/ants/releases/tag/v2.4.0).
- orb: [v0.2.2](https://github.com/paulmach/orb/releases/tag/v0.2.2).
- prometheus: [v1.11.0](https://github.com/prometheus/client_golang/releases/tag/v1.11.0).
- testify: [v1.7.0](https://github.com/stretchr/testify/releases/tag/v1.7.0).
- atomic: [v1.8.0](https://github.com/uber-go/atomic/releases/tag/v1.8.0).
- zap: [v1.17.0](https://github.com/uber-go/zap/releases/tag/v1.17.0).
- grpc: [v1.38.0](https://github.com/grpc/grpc-go/releases/tag/v1.38.0).
- cast: [v1.3.1](https://github.com/spf13/cast/releases/tag/v1.3.1).
- cobra: [1.1.3](https://github.com/spf13/cobra/releases/tag/v1.1.3).
- viper: [v1.8.1](https://github.com/spf13/viper/releases/tag/v1.8.1). 

## [0.21.1] - 2021-06-10

### Fixed
- Session token lifetime check (#589).
- Payload size check on the relayed objects (#580).

### Added
- VMagent to collect metrics from testnet storage image

### Changed
- Updated neofs-api-go to v1.27.1 release.

## [0.21.0] - 2021-06-03 - Seongmodo (석모도, 席毛島)

Session token support in container service, refactored config in storage node,
TLS support on gRPC servers.

### Fixed
- ACL service traverses over all RequestMetaHeader chain to find 
  bearer and session tokens (#548).
- Object service correctly resends complete objects without attached
  session token (#501).
- Inner ring processes `neofs.Bind` and `neofs.Unbind` notifications (#556).
- Client cache now gracefully closes all available connections (#567).

### Added
- Session token support in container service for `container.Put`, 
  `container.Delete` and `container.SetEACL` operations.
- Session token support in container and sign command of NeoFS CLI.
- TLS encryption support of gRPC service in storage node.

### Changed
- Inner ring listens RoleManagement contract notifications to start governance
  update earlier.
- Inner ring processes extended ACL changes.
- Inner ring makes signature checks of containers and extended ACLs.
- Refactored config of storage node. 
- Static clients from `morph/client` do not process notary invocations 
  explicitly anymore. Now notary support specified at static client creation.
- Updated neo-go to v0.95.1 release.
- Updated neofs-api-go to v1.27.0 release.

### Removed
- Container policy parser moved to neofs-sdk-go repository.
- Invoke package from inner ring.

## [0.20.0] - 2021-05-21 - Dolsando (돌산도, 突山島)

NeoFS is N3 RC2 compatible. 

### Fixed
- Calculations in EigenTrust algorithm (#527).
- NPE at object service request forwarding (#532, #543, #544).
- FSTree iterations in blobstor (#541).
- Inhume operation in storage engine (#546).

### Added
- Optional endpoint to main chain in storage app.
- Client for NeoFSID contract.

### Changed
- Reorganized and removed plenty of application configuration records 
  (#510, #511, #512, #514).
- Nodes do not resolve remote addresses manually.
- Presets for basic ACL in CLI are `private` ,`public-read` and
  `public-read-write` now.
- Updated neo-go to v0.95.0 release.
- Updated neofs-api-go to v1.26.1 release.
- Updated go-multiaddr to v0.3.1 release.

### Removed
- Unused external GC workers (GC is part of the shard in storage engine now).
- Unused worker pools for object service in storage app.
- `pkg/errors` dependency (stdlib errors used instead).

## [0.19.0] - 2021-05-07 - Daecheongdo (대청도, 大靑島)

Storage nodes exchange, calculate, aggregate and store reputation information
in reputation contract. Inner ring nodes support workflows with and without
notary subsystem in chains. 

### Fixed
- Build with go1.16.
- Notary deposits last more blocks. 
- TX hashes now prints in little endian in logs.
- Metabase deletes graves regardless of the presence of objects.
- SplitInfo error created from all shards instead of first matched shard.
- Possible deadlock at cache eviction in blobovnicza.
- Storage node does not send rebootstrap messages after it went offline.

### Added
- Reputation subsystem that includes reputation collection, exchange, 
calculation and storage components.
- Notary and non notary workflows in inner ring.
- Audit fee transfer for inner ring nodes that performed audit.
- Unified encoding for all side chain payment details.
- New write cache implementation for storage engine.
- NEP-2 and NEP-6 key formats in CLI.

### Changed
- Metabase puts data in batches.
- Network related new epoch handlers in storage node executed asynchronously. 
- Storage node gets epoch duration from global config.
- Storage node resign and resend Search, Range, Head, Get requests of object
service without modification.
- Inner ring does not sync side chain validators in single chain deployment.
- neo-go updated to v0.94.1
- neofs-api-go updated to v1.26.0

## [0.18.0] - 2021-03-26 - Yeongheungdo (영흥도, 靈興島)

NeoFS operates with updated governance model. Alphabet keys and inner ring keys
are accessed from side chain committee and `RoleManagement` contract. Each epoch
alphabet keys are synchronized with main chain.

### Fixed
- Metabase does not store object payloads anymore.
- TTLNetCache now always evict data after a timeout.
- NeoFS CLI keyer could misinterpret hex value as base58. 

### Added
- Local trust controller in storage node.
- Governance processor in inner ring that synchronizes list of alphabet keys.

### Changed
- Inner ring keys and alphabet keys are managed separately by inner ring and
  gathered from committee and `RoleManagement` contract.

## [0.17.0] - 2021-03-22 - Jebudo (제부도, 濟扶島)

Notary contract support, updated neofs-api-go with raw client, some performance 
tweaks with extra caches and enhanced metrics.

### Added
- Notary contract support.
- Cache for morph client.
- Metrics for object service and storage engine.
- Makefile target for fast and dirty docker images.
- GAS threshold value in inner ring GAS transfers.

### Changed
- RPC client cache now re-used per address instead of (address+key) tuple.
- Updated neofs-api-go version to v1.25.0 with raw client support.
- Updated neo-go to testnet compatible v0.94.0 version.

## [0.16.0] - 2021-02-26 - Ganghwado (강화도, 江華島)

Garbage collector is now running inside storage engine. It is accessed
via Control API, from `policer` component and through object expiration
scrubbers. 

Inner ring configuration now supports single chain mode with any number of
alphabet contracts.

Storage node now supports NetworkInfo method in netmap service.

### Fixed
- Storage engine now inhumes object only in single shard.
- Metabase correctly removes parent data at batched children delete.
- Metabase does not accept tombstone on tombstone records in graveyard anymore.
- Object service now rejects expired objects.
- CLI now correctly attaches bearer token in storage group operations.
- Container policy parser now works with strings in filter key.
- Policer component now removes redundant objects locally.

### Added
- GC job that monitors expired objects.
- GC job that removes marked objects from physical storage.
- Batch inhume operations in metabase, shard and engine.
- `control.DropObjects` RPC method.
- Support of `netmap.NetworkInfo` RPC method.
- Single chain inner ring configuration.

### Changed
- `UN-LOCODE` node attribute now optional.
- `engine.Delete` method now marks object to be removed by GC.
- Inner ring node supports any number of alphabet contracts from 1 up to 40.

## [0.15.0] - 2021-02-12 - Seonyudo (선유도, 仙遊島)

NeoFS nodes are now preview5-compatible. 

IR nodes are now engaged in the distribution of funds to the storage nodes:
for the passed audit and for the amount of stored information. All timers 
of the IR nodes related to the generation and processing of global system 
events are decoupled from astronomical time, and are measured in the number 
of blockchain blocks.

For the geographic positioning of storage nodes, a global NeoFS location
database is now used, the key in which is a UN/LOCODE, and the base itself 
is generated on the basis of the UN/LOCODE and OpenFlights databases.

### Added
- Timers with time in blocks of the chain.
- Subscriptions to new blocks in blockchain event `Listener`.
- Tracking the volume of stored information by containers in the 
  storage engine and an external interface for obtaining this data.
- `TransferX` operation in sidechain client.
- Calculators of audit and basic settlements.
- Distribution of funds to storage nodes for audit and for the amount 
  of stored information (settlement processors of IR).
- NeoFS API `Container.AnnounceUsedSpace` RPC service.
- Exchange of information about container volumes between storage nodes 
  controlled by IR through sidechain notifications.
- Support of new search matchers (`STRING_NOT_EQUAL`, `NOT_PRESENT`).
- Functional for the formation of NeoFS location database.
- CLI commands for generating and reading the location database.
- Checking the locode attribute and generating geographic attributes 
  for candidates for a network map on IR side.
- Verification of the eACL signature when checking Object ACL rules.

### Fixed
- Overwriting the local configuration of node attributes when updating 
  the network map.
- Ignoring the X-headers CLI `storagegroup` commands.
- Inability to attach bearer token in CLI `storagegroup` commands.

### Changed
- Units of epoch and emit IR intervals.
- Query language in CLI `object search` command.

### Updated
- neo-go v0.93.0.
- neofs-api-go v1.23.0.

## [0.14.3] - 2021-01-27

### Fixed
- Upload of objects bigger than single gRPC message.
- Inconsistent placement issues (#347, #349).
- Bug when ACL request classifier failed to classify `RoleOthers` in 
  first epoch.

### Added
- Debug section in readme file for testnet configuration.

### Changed
- Docker images now based on alpine and contain shell.
- Node bootstraps with active state in node info structure.

## [0.14.2] - 2021-01-20

Testnet4 related bugfixes.

### Fixed 
- Default values for blobovnicza object size limit and blobstor small object 
  size are not zero.
- Various storage engine log messages.
- Bug when inner ring node ignored bootstrap messages from restarted storage
  nodes. 
  
### Added
- Timeout for reading boltDB files at storage node initialization.

### Changed
- Increased default extra GAS fee for contract invocations at inner ring.

## [0.14.1] - 2021-01-15

### Fixed

- Inner ring node could not confirm `netmap.updateState` notification.
- `object.RangeHash` method ignored salt values.

### Added

- Control API service for storage node with health check, netmap and node state
  relate methods.
- Object service now looks to previous epoch containers.
- Possibility to configure up multiple NEO RPC endpoints in storage node.

### Changed

- Storage node shuts down if event producer RPC node is down.

## [0.14.0] - 2020-12-30 - Yeouido (여의도, 汝矣島)

Preview4 compatible NeoFS nodes with data audit.

### Added
- Data audit routines in inner ring nodes.
- Storage group operations in CLI (`neofs-cli storagegroup --help`).

### Fixed
- Loss of request X-headers during the forwarding in Object service.

### Changed
- Updated neo-go version for preview4 compatibility.

### Updated
- neo-go v0.92.0.
- neofs-api-go v1.22.0.

## [0.13.2] - 2020-12-24

Support changes from neofs-api-go v1.21.2 release.

### Added

- Support of request X-Headers in CLI commands.

### Changed

- Use updated API of container library.

## [0.13.1] - 2020-12-18

Fixes based on Modo release testing results.

### Added

- Verification of chain element addresses during object assembling.

### Changed

- Processing of filters by non-address fields in Object Range/RangeHash/Delete.

### Fixed

- `Graveyard` and `ToMoveIt` bucket names in metabase.
- Double formation of the parent title when transforming an object.
- Loss of session token during Object Put.
- Potential generating Range requests inside Get request execution context.

## [0.13.0] - 2020-12-15 - Modo (모도, 茅島)

Implementation of a local object storage engine.
Adaptation of the object service work scheme for the engine.

### Changed

- Object format after transformations.
- Handling of object operations.

### Added

- Local storage components: `Engine`, `Shard`, `BlobStor`,
  `Metabase`, `Blobovnicza`.
- Support of voting for sidechain governance in IR node.
- `Raw` flag support in Object Get/Head/GetRange CLI commands.

### Fixed

- Ignoring object address from session token in eACL validation.

## [0.12.1] - 2020-11-25

Bugfixes and small performance improvements.

### Fixed

- Routine leak by adding SDK client cache. (#184)
- Variety of ACL bugs. (#180, #190, #209)
- Policer tried to replicate virtual objects. (#182)
- Search queries with object ID field. (#177)
- Bug with extended ACL signature check in neofs-cli (#206)

### Added

- More debug logs in object service.
- Dial timeouts in object service config (`NEOFS_OBJECT_PUT_DIAL_TIMEOUT=5s`)

### Changed

- Routine pools in object service are non-blocking now.
- Container service now returns error if extended ACL is not set.

## [0.12.0] - 2020-11-17

NeoFS-API v2.0 support and updated brand-new storage node application.

### Fixed

- SetConfig method invocation of netmap contract. (#147)
- Balance response overflow. (#122)

### Added

- Gas emission routine in inner ring nodes.
- GRPC reflection service. (`NEOFS_GRPC_ENABLE_REFLECT_SERVICE=true`)
- New netmap query language parser.

### Changed

- Storage node application rebuilt from scratch.
- CLI supports accounting, object and container related operations.
- Inner ring node shutdowns on neo RPC node connection drop.
- Updated to preview4 compatible neo-go version.

## [0.11.0] - 2020-07-23

### Added

- Inner ring application to repository.
- Inner ring epoch processor.
- Inner ring asset processor for GAS deposit and withdraw.

### Changed

- The structure of source code tree.

## [0.10.0] - 2020-07-10

First public review release.

[Unreleased]: https://github.com/nspcc-dev/neofs-node/compare/v0.27.5...master
[0.27.5]: https://github.com/nspcc-dev/neofs-node/compare/v0.27.4...v0.27.5
[0.27.4]: https://github.com/nspcc-dev/neofs-node/compare/v0.27.3...v0.27.4
[0.27.3]: https://github.com/nspcc-dev/neofs-node/compare/v0.27.2...v0.27.3
[0.27.2]: https://github.com/nspcc-dev/neofs-node/compare/v0.27.1...v0.27.2
[0.27.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.27.0...v0.27.1
[0.27.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.26.1...v0.27.0
[0.26.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.26.0...v0.26.1
[0.26.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.25.1...v0.26.0
[0.25.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.25.0...v0.25.1
[0.25.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.24.1...v0.25.0
[0.24.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.24.0...v0.24.1
[0.24.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.23.1...v0.24.0
[0.23.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.23.0...v0.23.1
[0.23.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.22.3...v0.23.0
[0.22.3]: https://github.com/nspcc-dev/neofs-node/compare/v0.22.2...v0.22.3
[0.22.2]: https://github.com/nspcc-dev/neofs-node/compare/v0.22.1...v0.22.2
[0.22.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.22.0...v0.22.1
[0.22.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.21.1...v0.22.0
[0.21.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.21.0...v0.21.1
[0.21.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.20.0...v0.21.0
[0.20.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.19.0...v0.20.0
[0.19.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.18.0...v0.19.0
[0.18.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.17.0...v0.18.0
[0.17.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.16.0...v0.17.0
[0.16.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.15.0...v0.16.0
[0.15.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.14.3...v0.15.0
[0.14.3]: https://github.com/nspcc-dev/neofs-node/compare/v0.14.2...v0.14.3
[0.14.2]: https://github.com/nspcc-dev/neofs-node/compare/v0.14.1...v0.14.2
[0.14.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.14.0...v0.14.1
[0.14.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.13.2...v0.14.0
[0.13.2]: https://github.com/nspcc-dev/neofs-node/compare/v0.13.1...v0.13.2
[0.13.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.13.0...v0.13.1
[0.13.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.12.1...v0.13.0
[0.12.1]: https://github.com/nspcc-dev/neofs-node/compare/v0.12.0...v0.12.1
[0.12.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/nspcc-dev/neofs-node/compare/v0.10.0...v0.11.0
