package repository

import (
	"testing"

	_ "cdp-service/infra" //Implicit initialization infra
)

func TestClondToken(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *ClondRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	t.Log(repo.NewUptoken())
	t.Log(repo.NewUptoken("xxx.jpg"))
}

func TestClondUpload(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *ClondRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	err := repo.UpLoad("test.csv", []byte("1,2,3,4"))
	if err != nil {
		panic(err)
	}
	data, err := repo.PrivateDownload("test.csv")
	if err != nil {
		panic(err)
	}
	t.Log(string(data))
}
