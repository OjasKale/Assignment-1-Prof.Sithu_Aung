package main

import (
	"bytes"
	"math"
	"net/rpc"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Call struct {
	List flist `json:"list"`
}

type fmeta struct {
	Type  string `json:"-"`
	Start int32  `json:"-"`
	Count int32  `json:"-"`
}

type flist struct {
	Meta      fmeta        `json:"-"`
	Resources []fresources `json:"resources"`
}

type fresources struct {
	Resource fresource `json:"resource"`
}

type fresource struct {
	Classname string  `json:"classname"`
	Fields    ffields `json:"fields"`
}

type ffields struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}

type Args struct {
	StockSymbolAndPercentage string
	UserBudget               float64
}

type Quote struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
	TrdId         int     `json:"id"`
}

type Id struct {
	TrdId int `json:"id"`
}

type UpdQuote struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
}

type StkC int

var M map[int]Quote


func (t *StkC) StockPrice(args *Args, quote *Quote) error {
  quote.TrdId++
	a := string(args.StockSymbolAndPercentage[:])

	a = strings.Replace(a, ":", ",", -1)
	a = strings.Replace(a, "%", ",", -1)
	a = strings.Replace(a, ",,", ",", -1)
	a = strings.Trim(a, " ")
	a = strings.Replace(a, "\"", "", -1)
	a = strings.TrimSpace(a)
	a = strings.TrimSuffix(a, ",")
	Stkarr := strings.Split(a, ",")

	Total := 0.0
	var ReqUrl string

	for i := 0; i < len(Stkarr); i++ {
		i = i + 1

		temp, _ := strconv.ParseFloat(Stkarr[i], 64)
		Total = (temp * args.UserBudget * 0.01)
		fmt.Println(Stkarr[i-1], Total)
		ReqUrl = ReqUrl + (Stkarr[i-1] + ",")

	}

	ReqUrl = strings.TrimSuffix(ReqUrl, ",")

	UrlStr := "http://finance.yahoo.com/webservice/v1/symbols/" + ReqUrl + "/quote?format=json"

	client := &http.Client{}

	resp, _ := client.Get(UrlStr)
	req, _ := http.NewRequest("GET", UrlStr, nil)

	req.Header.Add("If-None-Match", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	resp, _ = client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var C Call
		body, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(body, &C)

		n := len(Stkarr)

		Quo := make([]float64, n, n)

		for i := 0; i < n; i++ {
			i = i + 1
			TempFloat, _ := strconv.ParseFloat(Stkarr[i], 64)
			Quo[i] = (TempFloat * args.UserBudget * 0.01)
			fmt.Println(Quo)
			fmt.Println(Stkarr[i-1], Quo[i])
		}

		var buffer bytes.Buffer
		q := 0
		for _, Sample := range C.List.Resources {

			temp1 := Sample.Resource.Fields.Symbol
			temp2, _ := strconv.ParseFloat(Sample.Resource.Fields.Price, 64)
			temp3 := (int)(Quo[q+1] / temp2)
			temp4 := math.Mod(Quo[q+1], temp2)
			q = q + 2

			quote.Stocks = fmt.Sprintf("%s:%g:%d", temp1, temp2, temp3)
			quote.UnvestedAmount = quote.UnvestedAmount + temp4
			buffer.WriteString(quote.Stocks)
			buffer.WriteString(",")
		}

		quote.Stocks = (buffer.String())
		quote.Stocks = strings.TrimSuffix(quote.Stocks, ",")
		fmt.Println(quote.Stocks)
		fmt.Println(quote.UnvestedAmount)

		M = map[int]Quote{
			quote.TrdId: {quote.Stocks, quote.UnvestedAmount, quote.TrdId},
		}

		fmt.Print(M)
		if err == nil {
			fmt.Println("Completed")
		}
	} else {
		fmt.Println(resp.Status)

	}


	return nil
}




func (t *StkC) Portfo(id *Id,qt *Quote) error{
if(id.TrdId==1){
*qt = M[id.TrdId]
}

return nil

}

func main() {

	StkC := new(StkC)
	rpc.Register(StkC)


	rpc.HandleHTTP()

	err := http.ListenAndServe(":1528", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
