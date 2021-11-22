package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	captcha "github.com/widaT/go-captcha/puzzle_captcha"
)

type RespJson struct {
	RepCode string      `json:"repCode"`
	RepData interface{} `json:"repData"`
	RepMsg  string      `json:"repMsg"`
}

type Data struct {
	OrgImg    string `json:"originalImageBase64"`
	BlkImg    string `json:"jigsawImageBase64"`
	Token     string `json:"token"`
	SecretKey string `json:"secretKey"`
}

type Item struct {
	Point   *captcha.Point
	Expired time.Time
}

type CheckParams struct {
	Point *captcha.Point `json:"point"`
	Token string         `json:"token"`
}

type VerificationParams struct {
	Verification string `json:"verification"`
}

var (
	lock  sync.Mutex
	cache = make(map[string]*Item) //实际项目组应该放到redis
)

func cos(handle func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	a := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")                                // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		rw.Header().Add("Access-Control-Allow-Credentials", "true")                        //设置为true，允许ajax异步请求带cookie信息
		rw.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers,x-requested-with")
		rw.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		handle(rw, req)
	}
	return a
}

func main() {

	captcha.LoadBackgroudImages("./images/puzzle_captcha/backgroud")
	captcha.LoadBlockImages("./images/puzzle_captcha/block")

	http.HandleFunc("/captcha/get", cos(func(rw http.ResponseWriter, req *http.Request) {
		ret, err := captcha.Run()
		if err != nil {
			return
		}
		data := &Data{
			BlkImg: ret.BlockImg,
			OrgImg: ret.BackgroudImg,
			Token:  uuid.NewV4().String(),
		}

		lock.Lock()
		cache[data.Token] = &Item{
			Point:   ret.Point,
			Expired: time.Now().Add(3 * time.Second),
		}
		lock.Unlock()

		resp := &RespJson{
			RepCode: "0000",
			RepData: data,
		}
		buf, err := json.Marshal(resp)
		if err != nil {
			return
		}
		rw.Write(buf)
	}))

	http.HandleFunc("/captcha/check", cos(func(rw http.ResponseWriter, req *http.Request) {
		reqdata := CheckParams{}
		resp := &RespJson{
			RepCode: "0000",
			RepData: nil,
		}
		err := json.NewDecoder(req.Body).Decode(&reqdata)
		if err != nil {
			fmt.Println(err)
			resp = &RespJson{
				RepCode: "0007",
				RepMsg:  "参数错误",
			}
			goto end
		}

		lock.Lock()
		defer lock.Unlock()
		if i, found := cache[reqdata.Token]; found {
			if i.Expired.Before(time.Now()) {
				resp.RepCode = "0005"
				resp.RepMsg = "过期了"
			} else {
				delete(cache, reqdata.Token) //删除缓存
				if err := captcha.Check(reqdata.Point, i.Point); err != nil {
					resp.RepCode = "0004"
					resp.RepMsg = "位置不正确"
				} else {
					//二次认证数据 入缓存
					verificationKey := fmt.Sprintf("second:%s",
						md5Str(fmt.Sprintf("%s--- %d---%d", reqdata.Token, reqdata.Point.X, reqdata.Point.Y)))
					//第二次判断只要判断key在不在缓存里头
					cache[verificationKey] = &Item{
						Expired: time.Now().Add(3 * time.Second),
					}
					resp.RepData = VerificationParams{Verification: verificationKey}
				}
			}
		} else {
			resp.RepCode = "0001"
			resp.RepMsg = "token不存在"
		}

	end:
		buf, _ := json.Marshal(&resp)
		rw.Write(buf)
	}))

	//二次验证
	http.HandleFunc("/captcha/verification", cos(func(rw http.ResponseWriter, req *http.Request) {
		v := VerificationParams{}
		json.NewDecoder(req.Body).Decode(&v)
		resp := &RespJson{
			RepCode: "0000",
			RepData: nil,
		}
		lock.Lock()
		defer lock.Unlock()
		if i, found := cache[v.Verification]; found {
			if i.Expired.Before(time.Now()) {
				resp.RepCode = "0005"
				resp.RepMsg = "过期了"
			}
		} else {
			resp.RepCode = "0001"
			resp.RepMsg = "token不存在"
		}
		buf, _ := json.Marshal(&resp)
		rw.Write(buf)
	}))

	http.Handle("/", http.FileServer(http.Dir("./dist")))
	panic(http.ListenAndServe(":8081", nil))
}
func md5Str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
