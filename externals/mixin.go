package externals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MixinNetwork/safe/governance/config"
)

const (
	mixinRPC = "https://rpc.mixin.dev"
)

type Node struct {
	Id          string `json:"id"`
	Signer      string `json:"signer"`
	Payee       string `json:"payee"`
	State       string `json:"state"`
	Timestamp   int64  `json:"timestamp"`
	Transaction string `json:"transaction"`
}

type Transaction struct {
	Asset string `json:"asset"`
	Extra string `json:"extra"`
	Hash  string `json:"hash"`
}

func ListAllNodes() ([]*Node, error) {
	if config.AppConfig.Environment != "prod" {
		return nodes, nil
	}
	data, err := callMixinRPC("listallnodes", []any{0, false})
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func ReadTransaction(hash string) (*Transaction, error) {
	data, err := callMixinRPC("gettransaction", []any{hash})
	if err != nil {
		return nil, err
	}
	var tx Transaction
	err = json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func callMixinRPC(method string, params []any) ([]byte, error) {
	client := &http.Client{Timeout: 20 * time.Second}

	body, err := json.Marshal(map[string]any{
		"method": method,
		"params": params,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", mixinRPC, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data  any `json:"data"`
		Error any `json:"error"`
	}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("ERROR %s", result.Error)
	}

	return json.Marshal(result.Data)
}

var nodes = []*Node{
	{
		Id:     "394e7b2131b7d0a996bb094e30d05ac7d51f5a09156e5f7349cac55d2179a144",
		Signer: "XIN4qtYcAuAsJFnHp61waUheVsiK1byouLqbhrA8VpSQwxHs4z8LPjpFRrx3zdmiXZuFSwJ8CAMCwLkxap1LbRWHk2iVsLyx",
		Payee:  "XINYvDWLAqoa1PxNxAaJcecrrehHVaaqqT4owg7ST1Yt2Gs5VUX62ArnVW7rx3vBMxfRdA5Y6kEg1Y5jSdQDFF3msunpmED4",
		State:  "ACCEPTED",
		// custodian XINJYiri2BU4dLGdsj33C5pvDuhzxK7DmWB9PvABa7u53tCoabApajFRsNTbsLjm2tjPfRQJEN2Awpe8SP3V35CMGRm2A5N1
	},
	{
		Id:     "2b0636403194b897a2d92d54060dd84acab78139626db2d919ce9ca84d64a433",
		Signer: "XINWJr86Q6h8SHfpYSFNKviyqQpcBC4c2NyHNRf1g6JrZP6LJFvd9eKYjz78yGcPasMpW5qGJZHKzAJugCfct95ywuoXbRt",  // 828bbdfc591db85aa202fe1d69ee832a86f1d789d70d5abc900da4aa45f2a30d
		Payee:  "XINNHgSRGWP9fF7dm7SWdaKcKfGwVZZPn4Xuypxx1n9n93XNqEHbD7QLwAE8xVFg6UNHACqPo4sdAsvvXUbnpAQe9EWVPUCk", // 83fd6d9b0969a6e450ea0acec51b88d7a2eea137e0e19d3cfa96225e8a729209
		State:  "ACCEPTED",
		// custodian XINR9cLBybXuQ7g2S3QFmFFe3pUXKcrvrZtEAChqFXxDo5a8vQAaRmA7SqaxoyRBoVohyf7kmkMT7UGLzXEXVAeXHB4wAnUk 1bf1c616d321bc5ef9c615eaea12dd8996f953ee2d9313cdce20f543dac2bb0f
	},
	{
		Id:     "cb5cd1a02f94ca98c060769e7c98f62cd14559b62d10762465ae46a51b69a432",
		Signer: "XINZjfRoqpppuoQY87fCP15BvTAh4khKr9PQ3UCKeSs4bQ3g1ZJvCrWSziWxCrg2KxNEfY7KGpak1DfQUUJirFx3iohXF5KL", // fbd36e62926af4e83b8c7c89fc37440125cbc1e0e4657b29ae010f579bf42907
		Payee:  "XIN3u5Q8qf1C96WRJuCiGCSx3qy3awy7hmT7m5rUbE7CcVHJpSqxYjoyDeGKmVhNAZkVitcye534uXXXuMCZTsfbNevqLsz9", // 26c0176f0aa98d2f6861d702572200ca8483e1f6a557a45e0379f2aa73bf0609
		State:  "ACCEPTED",
	},
	{
		Id:     "ff19e930753846c676c340042d5547021b5d59be3181f471b7fd7c332d613672",
		Signer: "XINW3Akp1yC1LKtQAv7YdAiRufdLQN6uwfVUkcpt976k4fdzJFtazMH8vgtfFux7qiHPbnL6LL3ShFjzK3HPgt3yewydnSm6", // d8ec369499f613fd30914678a746d02a1ff58b686e33db732637aeaf6eec660b
		Payee:  "XINHCU4KJj3XJT3shyYSoRp3RPQag3MaQc36xaDwqraVs6HZDu4r5t7vSHk6zm6rFmXENGMQcphq5ZhikwA5bfeZexXKqsof", // 1e6b011f81bb243228ca7056c1c8daa9db66868e60a3bcb9eb417f50e189be05
		State:  "ACCEPTED",
	},
	{
		Id:     "c681e456ef357bf721c6321c9585d49769e98af721464169a92c595be13708e9",
		Signer: "XINA8SjXhu3jHrie7VtAiCFtTJQ9iBqLtTEuJA3vDge97UMvEMrX8HCmwFvyv4VjLhvFdTSVjS4Z4yMq84oNBrN31xji2UYK", // bfd4a043d4e8a34e19b452987273ed954a0e403230321e361c5454d2ae46cf09
		Payee:  "XINU3AnRQQL61XsxRNKLRpQQ6Jbzgf8Brbrmrdhruj6xf7ggfhJ2c2aQ9ekZtNHXuFA5FkVkzusK4XuH82ycdYuRD3C3Zomv", // 7c77a98f39ca91c9a22a7a704400157a9e0464e41d65f96a29ce210850125a0b
		State:  "ACCEPTED",
	},
	{
		Id:     "3c15dc056af8141ec271790e278f5698a9cfc1c73711ff07099218f09967e125",
		Signer: "XINYTkPnoV73s6yYWzj56TwHLqnczPsueAoQM5RtBoY6uSAucBgutgjZ5vjwuFtSRN9nbMZDZYJ16wbGqTAZGVsDeLywrhwP", // 3f525bde115db71945772eeee79d2826047b61141a88332258354cef540c0b08
		Payee:  "XINALJ7YWynYaQYPqF6pc5nz8wvo7CPYkTtF4dSmV9b5qLQoz2YWy27AS4WK3TGoG7o3C2FqWHo4uZqTnseAYJwq3suHYDM3", // eaf458dd97ca87149fe7582ae4b029c58c8cf7e684d67f65a8d2eab713dca307
		State:  "ACCEPTED",
	},
	{
		Id:     "7593faea3acfa45c4f784be9fe3239a97b4d48da7ed8ea63725c3cd0c10ce4ca",
		Signer: "XINYv6MuELTyED1KAz8ea5cx4SdhZPF2XqXuarQKAAYrMbiJpx4sj2jRHAKnNu5AvJv6BWF3EAtAaXCAPm4GNZySiYcjg3Cz", // bec96ce8e1c7183d2b21e8a7ead42e2029b0bacf075801a7b60b7d9c32bddf06
		Payee:  "XINCpBi1JXRQwzXQLYRXQBLD6P2G3Ri3VnKiFSmUbwaus9CBZBV1uCzKCYd8HZRch6QMg6JR3FKfYM7vsMuVXnpr7gFhw5cp", // d7100ee458f36d4f080877d50cc492116ccbe5051a34eb8f99bac83ef09e3408
		State:  "ACCEPTED",
	},
	{
		Id:     "d5b022e636fc22f659acdefc9ca28160290a08e6c71287fce114b89eced5bf52",
		Signer: "XINEJRz5htVMnTU9q6VRBnbX8obGVdxaeKR3xHbE6vSY99pbm59RCB4Uc9ijr7dmGczMPH3DrDWGERbxfhNFfkckPFzcps35", // e2076ffffedfa64253582764e66115b3be047794e1a409b03e909cde60225e01
		Payee:  "XINWYwofRvRoK5e3CjBStuYN2rW4vyBy8VD1ct9MJSmB44kXw5JbSVSLfbBL7sWGq3FYe6Xxv8kpHhi3Fpzuio9yePFw84ME", // 8bfb5cf1fa932515ef6e08de2d02bc2b68fee650b0c4a0058a77f85a41cbbc07
		State:  "ACCEPTED",
	},
	{
		Id:     "3034606592d091cc956f16ced742bfa437460b1679bb303cd9010ab4e042d5d7",
		Signer: "XINZRTkU6HCnBESi75E9cZtLnwBxSEKoAudmVwEuNh2V6t1Jq7cProBBFFwR81NRt5WaUc3qyranXdCkkQC2aquPKUTHLoVm", // c0b3e03ddab4c00e2147fed20878599094f3d8f270a46c1c9aec9aac0d932107
		Payee:  "XINJ59qzeH4offuyQfznrhmMQ9faGkQbbmQpGrRKutPXv3qKk3GztA7ZnEiuXHK6SCFU7a4mH7Sd82z4dp5aCTCym8xxdsQ7", // 11580753d0f136d867e544432775fe75893621657e28dd915dcb96125d196504
		State:  "ACCEPTED",
	},
	{
		Id:     "9b604fc84e1bfaa789d3ecf0f520ca795dd97bde57a13a39a9c71c69f04f8820",
		Signer: "XINa5ZfAdcXb5DDV5ZQqtqMg3RQ6pYdegbsXji7FuagLxDEsKYWzDEqktChJuVNUin5MEveMh8meV7z4bhD24Zww7MydNpws", // 1a277b1cc7425f374969003fe7306575a8301bd262ae75c1e80406a0f8834f0a
		Payee:  "XINRPo2yD2CQPHiC44kW4gJaBsAmWbLNsMUJPX3dPFTY921WALVF9nJfLiq5oN1jCBPVCnu5LScyAfcZ1qfe4pW5M8fozHJ4", // bee7b26fefd23a22bf5eb355ec0564719283abae43980567a4cbfb3e0ee86807
		State:  "ACCEPTED",
	},
	{
		Id:     "e926c344d57017c3a0b3437b3180a7c3a1c7d81d2bdecee090ca06b4d1563905",
		Signer: "XIN85VNtKeGsEEcPg5tFt58xdfzi3mRaLWFAZFhHaucoVyKBbsGAV6FPMrozyUsDvah4EDPRfBdTKKJ2vKfphTQ9HNgiBBKW", // 5b4850807626fddc02002221666964727af475b63e4eaad744b34a5d10868a05
		Payee:  "XINRCXhM7XXgy9duDxYrbLYC9XTFmwqfGyF2Um2zdvwDd2i8x95aes885VJ5w3uJNWCAdgtyFc8o4My1uX3FNeTRS887KJ7U", // 47f313a47dbf85ef0f32171fff3b0b4deb5b0c4f06bf0ddbc9d595ce20a5310b
		State:  "ACCEPTED",
	},
}
