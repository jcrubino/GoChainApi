// json deserialization not implemented, json structs for that lead the api code
// consider changing the gcfg pkg to built-in solution; csv or json
// clean up verbose statements for structured clean output
// add address, txn, hash client side validations ?
// check recheck json values
// remove log.Fatal
// repackage to "package chain" when ready
// structs naming convention is lowercase with underscore for spaces

package main

import (
	"bytes"
	"code.google.com/p/gcfg"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type config struct {
	Auth struct {
		Key    string
		Secret string
	}
	Mode struct {
		Network string
		Verbose bool
	}
}

type chain struct {
	APIKey    string
	APISecret string
	Network   string
	BaseUrl   string
	verbose   bool
}

type msg struct {
	message string `json:"message"`
}

// btc units in satoshis
type block struct {
	Confirmations     int      `json:"confirmations"`
	Difficulty        int      `json:"difficulty"`
	Hash              string   `json:"hash"`
	Height            int      `json:"height"`
	MerkleRoot        string   `json:"merkle_root"`
	Nonce             int      `json:"nonce"`
	PreviousBlockHash string   `json:"previous_block_hash"`
	Time              string   `json:"time"`
	TransactionHashes []string `json:"transaction_hashes"`
}

type address struct {
	Hash                string `json:"hash"`
	Balance             int64  `json:"balance"`
	Recieved            int64  `json:"recieved"`
	Sent                int64  `json:"sent"`
	UnconfirmedRecieved int64  `json:"unconfirmed_recieved"`
	UnconfirmedSent     int64  `json:"unconfirmed_sent"`
	UnconfirmedBalance  int64  `json:"unconfirmed_balance"`
}

type tx struct {
	Amount        int64  `json:"amount"`
	BlockHash     string `json:"block_hash"`
	BlockHeight   int64  `json:"block_height"`
	BlockTime     string `json:"block_time"`
	Confirmations int64  `json:"confirmations"`
	Fees          int64  `json:"fees"`
	Hash          string `json:"hash"`
	Inputs        []struct {
		Addresses       []string `json:"addresses"`
		OutputHash      string   `json:"output_hash"`
		OutputIndex     int64    `json:"output_index"`
		ScriptSignature string   `json:"script_signature"`
		TransactionHash string   `json:"transaction_hash"`
		Value           int64    `json:"value"`
	} `json:"inputs"`
	Outputs []struct {
		Addresses          []string `json:"addresses"`
		OutputIndex        int64    `json:"output_index"`
		RequiredSignatures int64    `json:"required_signatures"`
		Script             string   `json:"script"`
		ScriptHex          string   `json:"script_hex"`
		ScriptType         string   `json:"script_type"`
		Spent              bool     `json:"spent"`
		TransactionHash    string   `json:"transaction_hash"`
		Value              int64    `json:"value"`
	} `json:"outputs"`
}

type op_return struct {
	Hex               string   `json:"hex"`
	ReceiverAddresses []string `json:"receiver_addresses"`
	SenderAddresses   []string `json:"sender_addresses"`
	Text              string   `json:"text"`
	TransactionHash   string   `json:"transaction_hash"`
}

type block_op struct {
	Hex               string   `json:"hex"`
	ReceiverAddresses []string `json:"receiver_addresses"`
	SenderAddresses   []string `json:"sender_addresses"`
	Text              string   `json:"text"`
	TransactionHash   string   `json:"transaction_hash"`
}

func (api *chain) Initialize(apikey, apisecret, network string, verbose bool) error {
	if network != "bitcoin" && network != "testnet3" {
		var msg string
		msg = "Network must be 'bitcoin' or 'testnet3"
		log.Println(msg)
		return errors.New(msg)
	}
	api.verbose = verbose
	api.Network = network
	api.APIKey = apikey
	api.APISecret = apisecret
	api.BaseUrl = "https://api.chain.com/v1"
	return nil
}

func Client(api chain, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	req.SetBasicAuth(api.APIKey, api.APISecret)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if api.verbose == true {
			log.Printf("Error: %s", err)
		}

		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if api.verbose == true {
			log.Printf("Error: %s", err)
		}
	}
	return body, nil
}

func PutClient(api chain, url, data string) ([]byte, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(data)))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	req.SetBasicAuth(api.APIKey, api.APISecret)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if api.verbose == true {
			log.Printf("Error: %s", err)
		}

		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if api.verbose == true {
			log.Printf("Error: %s", err)
		}
	}
	return body, nil
}

func (api *chain) GetAddresses(addresses []string) (string, error) {
	var addr string
	if len(addresses) > 1 {
		addr = strings.Join(append(addresses), ",")
	} else if len(addresses) == 1 {
		addr = addresses[0]
	} else {
		if api.verbose == true {
			log.Println("Error: No Addresses Provided")
		}
		return "", errors.New("Error: No Addresses Provided")
	}

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + addr
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Get information on a single address
func (api *chain) GetAddress(address string) (string, error) {

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + address
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil
}

func (api *chain) GetAddressTransactions(address string, limit int) (string, error) {
	if limit > 500 {
		log.Println("Max Txn must not exceed 500\nAutomatically reduced to 500")
		limit = 500
	}

	str_limit := strconv.Itoa(limit)

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + address + "/transactions?limit=" + str_limit
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil

}

func (api *chain) GetAddressesTransactions(addresses []string, limit int) (string, error) {
	if limit > 500 {
		log.Println("Max Txn must not exceed 200\nAutomatically reduced to 200")
		limit = 500
	}

	var addr string
	if len(addresses) > 1 {
		addr = strings.Join(append(addresses), ",")
	} else if len(addresses) == 1 {
		addr = addresses[0]
	}

	str_limit := strconv.Itoa(limit)

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + addr + "/transactions?limit=" + str_limit
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil

}

func (api *chain) GetAddressUnspents(address string) (string, error) {

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + address + "/unspent"
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil

}

func (api *chain) GetAddressesUnspents(addresses []string) (string, error) {
	if len(addresses) > 200 {
		log.Println("Max # addresses must not exceed 200\nAutomatically reduced to first 200")
		addresses = addresses[0:200]
	}

	var addr string
	if len(addresses) > 1 {
		addr = strings.Join(append(addresses), ",")
	} else if len(addresses) == 1 {
		addr = addresses[0]
	}

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + addr + "/unspent"
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil

}

func (api *chain) GetAddressOpReturns(address string) (string, error) {

	url := api.BaseUrl + "/" + api.Network + "/addresses/" + address + "/op-returns"
	if api.verbose == true {
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil

}

func (api *chain) GetTransaction(hash string) (string, error) {
	url := api.BaseUrl + "/" + api.Network + "/transactions/" + hash

	if api.verbose == true {
		log.Println("\nGetTransaction")
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil
}

func (api *chain) GetTransactionOpReturn(hash string) (string, error) {
	url := api.BaseUrl + "/" + api.Network + "/transactions/" + hash + "/op-return"

	if api.verbose == true {
		log.Println("\nGetTransactionOpReturn")
		log.Println(url)
	}

	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil
}

// not yet tested
func (api *chain) SendTransaction(hexJson string) (string, error) {
	url := api.BaseUrl + "/transactions"

	if api.verbose == true {
		log.Println("\nSendTransaction")
		log.Println(url)
	}

	body, err := PutClient(*api, url, hexJson)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	return string(body), nil

}

func (api *chain) GetBlockByHash(hash string) (string, error) {
	url := api.BaseUrl + "/blocks/" + hash
	if api.verbose == true {
		log.Println("\nGetBlockByHash")
		log.Println(url)
	}
	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return string(body), nil
}

func (api *chain) GetBlockByHeight(height int) (string, error) {
	url := api.BaseUrl + "/blocks/" + strconv.Itoa(height)
	if api.verbose == true {
		log.Println("\nGetBlockByHeight")
		log.Println(url)
	}
	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return string(body), nil
}

func (api *chain) GetBlockLatest() (string, error) {
	url := api.BaseUrl + "/blocks/"
	if api.verbose == true {
		log.Println("\nGetBlockLatest")
		log.Println(url)
	}
	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return string(body), nil
}

func (api *chain) GetBlockOpReturnsByHash(hash string) (string, error) {
	url := api.BaseUrl + "/blocks/" + hash + "/op-returns"
	if api.verbose == true {
		log.Println("\nGetBlockByHash")
		log.Println(url)
	}
	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return string(body), nil
}

func (api *chain) GetBlockOpReturnsByHeight(height int) (string, error) {
	url := api.BaseUrl + "/blocks/" + strconv.Itoa(height) + "/op-returns"
	if api.verbose == true {
		log.Println("\nGetBlockByHeight")
		log.Println(url)
	}
	body, err := Client(*api, url)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return string(body), nil
}

func main() {

	api := chain{}

	var cfg config
	err := gcfg.ReadFileInto(&cfg, "config.cfg")
	if err != nil {
		log.Fatalln("Config not properly setup.")
	}

	// init api with APIKey, APISecret, Network (bitcoin or testnet3), verbose logging bool
	err = api.Initialize(cfg.Auth.Key, cfg.Auth.Secret, cfg.Mode.Network, cfg.Mode.Verbose)
	if err != nil {
		log.Fatalf("\n\nError: %s\n\n", err)
	} else {
		log.Println("Success: Api initialized\n\n")
	}

	address, err := api.GetAddress("17x23dNjXJLzGMev6R63uyRhMWP1VHawKc")
	if err != nil {
		log.Fatalf("\n\nError: %s\n\n", err)

	} else {
		log.Println("Success: GetAddress executed\n\n")
	}

	log.Println(address)

	addresses, err := api.GetAddresses([]string{"1VayNert3x1KzbpzMGt2qdqrAThiRovi8", "1Fi57hAqyYYwaQVdA7a9qSKfiukBbt31G3"})
	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Println(addresses)
	}

	txns, err := api.GetAddressTransactions("1K4nPxBMy6sv7jssTvDLJWk1ADHBZEoUVb", 10)

	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Println(txns)
	}

	txn_set, err := api.GetAddressesTransactions([]string{"1K4nPxBMy6sv7jssTvDLJWk1ADHBZEoUVb", "1VayNert3x1KzbpzMGt2qdqrAThiRovi8"}, 2)

	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Println(txn_set)
	}

	op_return, err := api.GetAddressOpReturns("1Bj5UVzWQ84iBCUiy5eQ1NEfWfJ4a3yKG1")
	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		log.Println(op_return)
	}

}
