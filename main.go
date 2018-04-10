package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	// index : http://quote.eastmoney.com/centerv2/hsgg
	// list  : http://quote.eastmoney.com/stocklist.html
	/*
		var urls = map[string]string{
			"沪深A股":     "http://quote.eastmoney.com/centerv2/hsgg/hsag",
			"上证A股":     "http://quote.eastmoney.com/centerv2/hsgg/shag",
			"深证A股":     "http://quote.eastmoney.com/centerv2/hsgg/szag",
			"新股":       "http://quote.eastmoney.com/centerv2/hsgg/xg",
			"中小板":      "http://quote.eastmoney.com/centerv2/hsgg/zxb",
			"创业板":      "http://quote.eastmoney.com/centerv2/hsgg/cyb",
			"沪股通":      "http://quote.eastmoney.com/centerv2/hsgg/hgt",
			"深股通":      "http://quote.eastmoney.com/centerv2/hsgg/sgt",
			"B股":       "http://quote.eastmoney.com/centerv2/hsgg/bg",
			"AB股比价(沪)": "http://quote.eastmoney.com/centerv2/hsgg/abgbj/shabgbj",
			"AB股比价(深)": "http://quote.eastmoney.com/centerv2/hsgg/abgbj/szabgbj",
			"风险警示板":    "http://quote.eastmoney.com/centerv2/hsgg/fxjsb",
			"两网及退市":    "http://quote.eastmoney.com/centerv2/hsgg/lwjts",
		}
	*/

	//s := `
	//{rank:[
	//"2,000908,景峰医药,6.11,0.56,10.09%,0.00,38150,23309643,5.55,6.11,6.11,6.11,-,-,-,-,-,-,-,-,0.00%,1.09,0.69,17.37,1999-02-03",
	//"1,600354,敦煌种业,8.42,07%,12.94,789760,631131936,7.65,7.75,8.42,7.43,-,-,-,-,-,-,-,-,0.00%,1.79,17.64,26.84,2004-01-15"
	//],pages:177,total:3526}`

	//maxWorker := 1000
	//maxQueue := 10000
	//
	////初始化工作池
	//dispatch := NewDispatcher(maxWorker)
	////指定任务的队列长度
	//JobQueue = make(chan Job, maxQueue)
	////一直运行调度器
	//dispatch.Run()
	//
	//for i := 0; i < 10*10000; i++ {
	//	p := PayLoad{
	//		Name: fmt.Sprintf("玩家-[%s] \r\n", strconv.Itoa(i)),
	//	}
	//	job := Job{PayLoad: p}
	//	JobQueue <- job
	//	//time.Sleep(time.Millisecond)
	//}
	//close(JobQueue)

	t := time.Now()

	//ch := make(chan int)
	//for i := 0; i < 5; i++ {
	//	go func(n int) {
	//		ch <- n
	//	}(i)
	//}
	//<-ch

	s := `2,300303,聚飞光电,3.36,0.31,10.16%,11.80,451664,149043339,3.05,3.01,3.36,3.00,-,-,-,-,-,-,-,-,0.00%,5.65,4.42,52.64,2012-03-19`
	datas := strings.Split(s, ",")
	fmt.Println(datas)

	//Run()

	elapsed := time.Since(t)
	fmt.Println("耗时:", elapsed)

}

func Run() {
	share := new(Share)
	share.Page = "1"
	share.PageSize = "100"
	share.Downloader("http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx").ApiSpider()

	fmt.Println(share.Data)
}

type Market struct {
	/*
		?,代码，名称，最新价，涨跌额，涨跌幅，振幅， 成交量(手)，成交额， 昨收，今开，最高，最低，
		-,-,-,-,-,-,-,-,五分钟涨跌，量比，换手率，市盈率，上市时间

		2,300303,聚飞光电,3.36,0.31,10.16%,11.80,451664,149043339,3.05,3.01,3.36,3.00,
		-,-,-,-,-,-,-,-,0.00%,5.65,4.42,52.64,2012-03-19
	*/
	id   string
	code string
	name string
}

type DataSource struct {
	Source *[]ResultData
}

type ResultData struct {
	Rank  []string `json:"rank"`
	Pages int64    `json:"pages"`
	Total int64    `json:"total"`
}

type Share struct {
	Total    int64
	Pages    int64
	PageSize string
	Page     string
	Body     io.ReadCloser
	Data     ResultData
	Market   Market
}

// Using HTTP.NewRequest access to specify the url of the body
func (s *Share) Downloader(url string) *Share {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//req.Header.Add("User-Agent", "浏览器")

	param := req.URL.Query()
	param.Add("type", "CT")
	param.Add("cmd", "C._A")
	param.Add("sty", "FCOIATA")
	param.Add("sortType", "(ChangePercent)")
	param.Add("sortRule", "-1")
	param.Add("page", s.Page)
	param.Add("pageSize", s.PageSize)
	param.Add("js", "var ztCJfqUC={rank:[(x)],pages:(pc),total:(tot)}")
	param.Add("token", "7bc05d0d4c3c22ef9fca8c2a912d779c")
	param.Add("jsName", "quote_123")
	param.Add("_g", "0.628606915911589")
	param.Add("_", "1522921744988")
	req.URL.RawQuery = param.Encode()

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	s.Body = res.Body
	return s
}

//Deal with illegal string, and illegal json format,
//finally deserialization for struct
func (s *Share) ApiSpider() *Share {
	body, err := ioutil.ReadAll(s.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Body.Close()

	//字符串处理，截取"="后半段
	str := string(body)
	start := strings.IndexAny(str, "=") + 1
	end := len(str)

	//处理非法json
	r_rank := strings.Replace(str[start:end], "rank", "\"rank\"", 1)
	r_page := strings.Replace(r_rank, "pages", "\"pages\"", 1)
	data := strings.Replace(r_page, "total", "\"total\"", 1)

	bjson := []byte(data)

	jsonparser.Get(bjson, "rank")
	jsonparser.GetInt(bjson, "pages")
	jsonparser.GetInt(bjson, "total")

	// 2,300303,聚飞光电,3.36,0.31,10.16%,11.80,451664,149043339,3.05,3.01,3.36,3.00,-,-,-,-,-,-,-,-,0.00%,5.65,4.42,52.64,2012-03-19
	// ?,代码，名称，最新价，涨跌额，涨跌幅，振幅， 成交量(手)，成交额， 昨收，今开，最高，最低，-,-,-,-,-,-,-,-,五分钟涨跌，量比，换手率，市盈率，上市时间
	//deserialization to struct

	result := ResultData{}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	s.Data = result
	return s
}

func (s *Share) PageSpider() {

	doc, err := goquery.NewDocumentFromReader(s.Body)
	defer s.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	node := doc.Find(".job-profile").Children().Eq(1).Find("dt")
	name, is_name := node.Find("a").Attr("title")
	fmt.Println(is_name)
	price, is_price := node.Find("span").Attr("title")

	fmt.Println(is_price)
	fmt.Println(name, price)
}

func (s *Share) PipLine() {

}
