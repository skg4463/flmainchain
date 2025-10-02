-------------------cross-chain relayer setting----------------------------------

* flmainchaind keys

fls17nglqjy7d3uvg4q78ck34jaxy9gxav85tz8j2h

fls1lf06w0zzhgjxelagatlwvg6kuqae42kprwaluj



* flstoraged keys

fls1alksc8mrx4dwwgqc5fan8uv06jjsy0ah25h8dg

fls1w2evy6hkh9h807pdlsxmlhcgdrlt0wr0xw6hw6



* 잔액 쿼리

flmainchaind query bank balances fls1ppkrhzyg6pf2gg7xxhyq8z8jcclgrn5j3kr765 

flstoraged query bank balances fls1ppkrhzyg6pf2gg7xxhyq8z8jcclgrn5j3kr765  --node tcp://localhost:26658



1\. 니모닉으로 각 체인에 계정 추가 (/.hermes/keys/flmainchain/keyring-test/relayer-flmainchain.json)

HERMES\_APP=~/.ignite/apps/hermes/bin/v1.13.1

CONFIG\_PATH=~/.ignite/relayer/hermes/flmainchain\_flstorage



abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art



$HERMES\_APP --config $CONFIG\_PATH keys add --chain flmainchain --mnemonic-file /dev/stdin

$HERMES\_APP --config $CONFIG\_PATH keys add --chain flstorage --mnemonic-file /dev/stdin





2\. 생성된 계정 확인 

$HERMES\_APP --config $CONFIG\_PATH keys list --chain flmainchain

$HERMES\_APP --config $CONFIG\_PATH keys list --chain flstorage



3\. 각 계정에 자금 채우기 

curl -X POST -d '{"address": "fls1r5v5srda7xfth3hn2s26txvrcrntldjuqs34tt", "coins": \["10000000stake"]}' http://localhost:4501/credit

curl -X POST -d '{"address": "fls1r5v5srda7xfth3hn2s26txvrcrntldjuqs34tt", "coins": \["10000000stake"]}' http://localhost:4500/credit



4\. 릴레이어 및 헤르메스 시작

$HERMES\_APP --config $CONFIG\_PATH create channel --a-chain flmainchain --b-chain flstorage --a-port fedlearning --b-port fedstoraging --new-client-connection --channel-version fedlearning-1



$HERMES\_APP --config $CONFIG\_PATH start



5.1  1라운드 초기화

MEMBERS=$(flmainchaind keys show alice -a),$(flmainchaind keys show bob -a)

echo $MEMBERS



5.2 init-round 트잭 전송

\# \[MEMBERS 변수]를 사용하여 alice와 bob을 1라운드의 초기 위원으로 설정

flmainchaind tx fedlearning init-round "$MEMBERS" --from alice --node tcp://localhost:26658 -y

-> 확인

flmainchaind query fedlearning show-current-round --node tcp://localhost:26658

flmainchaind query fedlearning show-round 1 --node tcp://localhost:26658



-----------------------------------------------------------------------



-------------------L-node simul------------------------------------



b9be8a5e6612896cf8b84f9e3a2b591858b1e9e47f4bffbab6f9b81be253fe64



오리지널해시 메인체인 보고

flmainchaind tx fedlearning submit-weight 1 b9be8a5e6612896cf8b84f9e3a2b591858b1e9e47f4bffbab6f9b81be253fe64 "1-bob-flstorage" --from fls1p9jqgjchwvue5fdm2ltvcasq5svfpw6k7wufpg --node tcp://localhost:26658 -y





flstoraged tx fedstoraging request-data-access b9be8a5e6612896cf8b84f9e3a2b591858b1e9e47f4bffbab6f9b81be253fe64 --from fls1w2evy6hkh9h807pdlsxmlhcgdrlt0wr0xw6hw6 -y



go run main.go b9be8a5e6612896cf8b84f9e3a2b591858b1e9e47f4bffbab6f9b81be253fe64 ./downloaded\_weight\_bob.bin





\# \[alice\_mainchain\_주소]와 \[bob\_mainchain\_주소]를 채워 넣기

flmainchaind tx fedlearning submit-score 1 fls1vva47v57ke002xn7n3qm43j70rr60mpntangj3 85 --from fls1vva47v57ke002xn7n3qm43j70rr60mpntangj3 --node tcp://localhost:26658 -y







fls1p9jqgjchwvue5fdm2ltvcasq5svfpw6k7wufpg

fls1vva47v57ke002xn7n3qm43j70rr60mpntangj3



flmainchaind tx fedlearning submit-global-model 1 b9be8a5e6612896cf8b84f9e3a2b591858b1e9e47f4bffbab6f9b81be253fa64 --from fls1p9jqgjchwvue5fdm2ltvcasq5svfpw6k7wufpg --node tcp://localhost:26658 -y



flmainchaind query fedlearning show-global-model 1 --node tcp://localhost:26658



