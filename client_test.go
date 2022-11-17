package go_aptos_sdk

var err error
var client *AleoClient

func init() {
	endpoint := ""
	client, err = NewAleoClient(endpoint)
	if err != nil {
		panic(err)
	}
}
