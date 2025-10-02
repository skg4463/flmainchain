# flmainchain
**flmainchain** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Environment
Ignite CLI version:             v29.3.1-dev

Ignite CLI source hash:         845a1a8886b8a098ed56372bab45ddee5caea526

Ignite CLI config version:      v1

Cosmos SDK version:             v0.53.3

Your go version:                go version go1.24.7 linux/amd64

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

Port info: 
    tendermint node(RPC): 26658,
    blockchain API: 1318,
    token faucet 4501,
    gPRC: 9091

```
query options: 


flmainchaind query bank balances ADDRESS --node tcp://localhost:26658
```


Relayer setting & hermes setting
```
⦁	잔액 쿼리
flmainchaind query bank balances fls1ppkrhzyg6pf2gg7xxhyq8z8jcclgrn5j3kr765 
flstoraged query bank balances fls1ppkrhzyg6pf2gg7xxhyq8z8jcclgrn5j3kr765  --node tcp://localhost:26658

1. 니모닉으로 각 체인에 계정 추가 (/.hermes/keys/flmainchain/keyring-test/relayer-flmainchain.json)
HERMES_APP=~/.ignite/apps/hermes/bin/v1.13.1
CONFIG_PATH=~/.ignite/relayer/hermes/flmainchain_flstorage

abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art

$HERMES_APP --config $CONFIG_PATH keys add --chain flmainchain --mnemonic-file /dev/stdin
$HERMES_APP --config $CONFIG_PATH keys add --chain flstorage --mnemonic-file /dev/stdin


2. 생성된 계정 확인 
$HERMES_APP --config $CONFIG_PATH keys list --chain flmainchain
$HERMES_APP --config $CONFIG_PATH keys list --chain flstorage

3. 각 계정에 자금 채우기 
curl -X POST -d '{"address": "fls1r5v5srda7xfth3hn2s26txvrcrntldjuqs34tt", "coins": ["10000000stake"]}' http://localhost:4501/credit
curl -X POST -d '{"address": "fls1r5v5srda7xfth3hn2s26txvrcrntldjuqs34tt", "coins": ["10000000stake"]}' http://localhost:4500/credit

4. 릴레이어 및 헤르메스 시작
$HERMES_APP --config $CONFIG_PATH create channel --a-chain flmainchain --b-chain flstorage --a-port fedlearning --b-port fedstoraging --new-client-connection --channel-version fedlearning-1

$HERMES_APP --config $CONFIG_PATH start

5.1  1라운드 초기화
MEMBERS=$(flmainchaind keys show alice -a),$(flmainchaind keys show bob -a)
echo $MEMBERS

5.2 init-round 트잭 전송
# [MEMBERS 변수]를 사용하여 alice와 bob을 1라운드의 초기 위원으로 설정
flmainchaind tx fedlearning init-round "$MEMBERS" --from alice --node tcp://localhost:26658 -y
-> 확인
flmainchaind query fedlearning show-current-round --node tcp://localhost:26658
flmainchaind query fedlearning show-round 1 --node tcp://localhost:26658

```

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers a frontend scaffolding feature (based on Vue) to help you quickly build a web frontend for your blockchain:

Use: `ignite scaffold vue`
This command can be run within your scaffolded blockchain project.


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/flmainchain@latest! | sudo bash
```
`username/flmainchain` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/ignite/installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.com/invite/ignitecli)
