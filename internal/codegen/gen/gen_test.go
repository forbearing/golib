package gen

import "os"

var dataServiceUserCreate string

func init() {
	var data []byte
	var err error
	if data, err = os.ReadFile("./testdata/service/user_create.go"); err != nil {
		panic(err)
	}
	dataServiceUserCreate = string(data)
}
