package main

import (
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

	//s := `{"rank":["1,603897,N长城,25.43,7.77,44.00%,0.00,353,898264,17.66,25.43,25.43,25.43,-,-,-,-,-,-,-,-,0.00%,-,0.08,25.57,2018-04-10","1,600929,湖南盐业,12.60,1.15,10.04%,0.00,49877,62845360,11.45,12.60,12.60,12.60,-,-,-,-,-,-,-,-,0.00%,6.48,3.33,72.52,2018-03-26","1,603106,恒银金融,26.68,2.43,10.02%,2.76,136226,362691648,24.25,26.68,26.68,26.01,-,-,-,-,-,-,-,-,0.00%,1.61,19.46,45.40,2017-09-20","2,002930,宏川智慧,23.95,2.18,10.01%,0.00,1644,3938026,21.77,23.95,23.95,23.95,-,-,-,-,-,-,-,-,0.00%,1.21,0.27,62.90,2018-03-28","1,603963,大理药业,38.46,3.50,10.01%,12.19,84788,305547632,34.96,34.20,38.46,34.20,-,-,-,-,-,-,-,-,0.00%,1.29,33.92,73.92,2017-09-22","1,603214,爱婴室,46.27,4.21,10.01%,0.00,1401,6483445,42.06,46.27,46.27,46.27,-,-,-,-,-,-,-,-,0.00%,2.21,0.56,49.45,2018-03-30","1,603289,泰瑞机器,18.92,1.72,10.00%,0.00,58518,110716643,17.20,18.92,18.92,18.92,-,-,-,-,-,-,-,-,0.00%,0.52,11.47,46.22,2017-10-31","2,300670,大烨智能,36.20,3.29,10.00%,5.17,185957,663184112,32.91,34.50,36.20,34.50,-,-,-,-,-,-,-,-,0.00%,4.22,68.87,53.17,2017-07-03","2,002873,新天药业,53.24,4.84,10.00%,11.86,54465,280705920,48.40,48.30,53.24,47.50,-,-,-,-,-,-,-,-,0.00%,1.35,31.63,41.36,2017-05-19","2,300705,九典制药,52.56,4.78,10.00%,13.50,127578,651973520,47.78,46.50,52.56,46.11,-,-,-,-,-,-,-,-,0.00%,0.94,43.48,89.94,2017-10-10","1,603383,顶点软件,59.09,5.37,10.00%,7.35,83432,482910704,53.72,55.50,59.09,55.14,-,-,-,-,-,-,-,-,0.00%,2.23,39.64,49.79,2017-05-22","2,300726,宏达电子,42.25,3.84,10.00%,11.06,213075,858802464,38.41,38.08,42.25,38.00,-,-,-,-,-,-,-,-,0.00%,1.02,53.14,63.44,2017-11-21","2,300660,江苏雷利,63.25,5.75,10.00%,4.61,18129,114202447,57.50,61.00,63.25,60.60,-,-,-,-,-,-,-,-,0.00%,4.16,7.17,22.15,2017-06-02","2,002900,哈三联,37.84,3.44,10.00%,14.65,203289,714939216,34.40,33.15,37.84,32.80,-,-,-,-,-,-,-,-,0.00%,1.70,38.53,44.11,2017-09-22","2,300688,创业黑马,72.61,6.60,10.00%,7.21,71501,513316608,66.01,67.85,72.61,67.85,-,-,-,-,-,-,-,-,0.00%,1.37,42.06,77.98,2017-08-10","1,603229,奥翔药业,24.53,2.23,10.00%,4.17,230922,564106320,22.30,24.53,24.53,23.60,-,-,-,-,-,-,-,-,0.00%,4.53,57.73,74.04,2017-05-09","1,603712,七一二,48.19,4.38,10.00%,10.02,371131,1740539216,43.81,48.19,48.19,43.80,-,-,-,-,-,-,-,-,0.00%,0.78,37.11,363.05,2018-02-26","2,300554,三超新材,99.03,9.00,10.00%,9.72,20523,196384450,90.03,92.00,99.03,90.28,-,-,-,-,-,-,-,-,0.00%,2.03,15.79,44.56,2017-04-21","2,002927,泰永长征,61.63,5.60,9.99%,9.82,116417,689369344,56.03,56.89,61.63,56.13,-,-,-,-,-,-,-,-,0.00%,1.18,49.64,41.52,2018-02-23","2,002919,名臣健康,40.63,3.69,9.99%,16.97,107918,407061504,36.94,36.07,40.63,34.36,-,-,-,-,-,-,-,-,3.15%,2.10,53.01,33.74,2017-12-18"],"pages":17,"total":338}`

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

	Run()

	elapsed := time.Since(t)
	fmt.Println("耗时:", elapsed)

}

func Run() {
	share := new(Share)
	share.Page = "1"
	share.PageSize = "3526"
	share.Downloader("http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx").ApiSpider()
	fmt.Println(share.MarketData, share.Total)
}

type Market struct {
	/*
		?,代码，名称，最新价，涨跌额，涨跌幅，振幅， 成交量(手)，成交额， 昨收，今开，最高，最低，
		-,-,-,-,-,-,-,-,五分钟涨跌，量比，换手率，市盈率，上市时间

		2,300303,聚飞光电,3.36,0.31,10.16%,11.80,451664,149043339,3.05,3.01,3.36,3.00,
		-,-,-,-,-,-,-,-,0.00%,5.65,4.42,52.64,2012-03-19
	*/

	Id               string
	Code             string
	Name             string
	NewPrice         string
	UpDownPrice      string
	UpDownRange      string
	ShakeRange       string
	TradeNum         string
	TradePrice       string
	PrevClosePrice   string
	TodayStartPrice  string
	HighPrice        string
	LowPrice         string
	UpDownRange5     string
	TradeNumRelative string
	TurnoverRate     string
	PeRatio          string
	ListedTime       string
}

type Share struct {
	Total      int64
	Pages      int64
	PageSize   string
	Page       string
	Body       io.ReadCloser
	MarketData []Market
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

	dataJson := []byte(data)

	//清洗数据到MarketData
	s.Pages, _ = jsonparser.GetInt(dataJson, "pages")
	s.Total, _ = jsonparser.GetInt(dataJson, "total")

	jsonparser.ArrayEach(dataJson, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Fatalln(err)
		}

		data := strings.Split(string(value), ",")

		m := Market{
			Id:               data[0],
			Code:             data[1],
			Name:             data[2],
			NewPrice:         data[3],
			UpDownPrice:      data[4],
			UpDownRange:      data[5],
			ShakeRange:       data[6],
			TradeNum:         data[7],
			TradePrice:       data[8],
			PrevClosePrice:   data[9],
			TodayStartPrice:  data[10],
			HighPrice:        data[11],
			LowPrice:         data[12],
			UpDownRange5:     data[21],
			TradeNumRelative: data[22],
			TurnoverRate:     data[23],
			PeRatio:          data[24],
			ListedTime:       data[25],
		}
		s.MarketData = append(s.MarketData, m)
	}, "rank")

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
