package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)


var kwd_stock = []string{"주가", "주주", "증시", "주식", "주린이", "종목", "업종", "투자", "마감", "브리핑", "애널리스트", "기관", "개미", "코인", "광풍", "저평가", "우량주", "KOSDAQ", "KOSPI", "다우", "지수", "리서치",
	"수혜주", "관련주", "산업주", "요약주", "통신주", "테마주", "특징주", "레포트", "코스피", "코스닥", "재테크", "제테크", "032640", "017670", "030200", "계열사", "시총", "기준가", "시작가", "현재가", "시작가",
	"양회", "유화", "타이어", "제강", "바이오", "화학", "정유", "손정의", "매수", "매도", "채굴", "지갑", "클립", "글로벌", "건축", "경영", "국면", "추세", "연봉", "자산", "금융", "소득", "MOU", "무역", "유통", "쇼핑",
	"유통", "물류", "공모"}
var kwd_estate = []string{"임대", "월세", "전세", "보증금", "부동산", "원룸", "고시텔", "오피스텔", "관리비", "풀옵션", "건물", "신축", "역세권", "대출", "펜션", "모텔", "호텔", "청약", "분양", "모델하우스", "집값", "중개"}
var kwd_sport = []string{"골프", "관중", "직관", "야구", "농구", "축구", "좌익수", "우익수", "중견수", "외야", "홈런", "주전", "공격", "수비", "감독", "스포츠", "운동", "선발", "KBO", "FC", "원정", "경기", "체육"}
var kwd_etc = []string{"렌터카", "신용카드", "판매완료", "판매 완료", "대란", "충전", "딜러", "중고차", "현대차", "수소차", "매월", "커피", "디자이너", "나태", "은행", "핀테크", "창업", "스타트업", "퀴즈", "백신", "년차",
	"증명", "접종", "신앙", "종교", "기독교", "공장", "자재", "주얼리", "쥬얼리", "보석", "정수기", "알바", "아르바이트", "이삿짐", "냉장고", "세탁기", "에어컨", "전자레인지", "청소기", "다이슨", "요리", "트럭",
	"환자", "보호자", "도서", "원서", "봉사", "패키지", "병원", "치과", "한의원", "블루오션", "가구", "마이크", "장난감", "랭킹", "추천인", "버스", "지하철", "리조트", "오션뷰", "정품", "한정판", "인구", "도시",
	"블라인드", "학원", "변기", "CD", "맛집", "비타민", "밸리", "자퇴", "뷰티", "화장품", "다이소", "이어폰", "불펜", "회복", "팬데믹", "네비", "내비", "카플레이", "안드로이드오토", "안드로이드 오토", "교수",
	"다이어트", "원단", "침대", "베개", "주행속도", "포토뉴스", "중증", "질환"}

var number = 0

type dataframe struct {
	keyword string
	title   string
	time    string
	text    string
	url     string
}

// 네이버 블로그   키워드별
func naver_blog_func(term ...string) {
	var url_1 = "https://search.naver.com/search.naver?where=view&sm=tab_viw.blog&query=" // URL 설정
	var url_2 = "&nso=so%3Add%2Cp%3A1h%2Ca%3Aall"

	var datas []dataframe
	//var df[]dataframe
	cha := make(chan []dataframe)
	lenge := len(term)

	for i := 0; i < lenge; i++ {
		go naver_blog_getPage(term[i], cha, url_1, url_2)
	}
	for i := 0; i < lenge; i++ {
		dataframe := <-cha
		datas = append(datas, dataframe...)
	}

	//writeExcel(datas)

	number = number + len(datas)
	//fmt.Println("done, extracted", len(datas))

}

// 채널 단위 당 페이지별 진행
func naver_blog_getPage(terms string, mainC chan<- []dataframe, url_1 string, url_2 string) {

	c := make(chan dataframe)
	var datas []dataframe
	pageURL := url_1 + terms + url_2
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	searchCards := doc.Find(".total_wrap")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go naver_blog_crawling(card, c, terms, pageURL)

	})

	for i := 0; i < searchCards.Length(); i++ {
		dataed := <-c
		datas = append(datas, dataed)
	}

	mainC <- datas

	return
}

// 페이지 정보 얻어오기
func naver_blog_crawling(card *goquery.Selection, c chan<- dataframe, terms string, purl string) {
	keyword := terms
	title := CleanString(card.Find(".total_tit").Text())
	time := CleanString(card.Find(".sub_time").Text())
	time = getTime(strings.Trim(time, "분전 "))
	text := CleanString(card.Find(".total_dsc").Text())
	url := purl
	c <- dataframe{
		keyword: keyword,
		time:    time,
		title:   title,
		text:    text,
		url:     url}
}

func naver_news_func(term ...string) {

	var url_1 = "https://search.naver.com/search.naver?where=news&query=" // URL 설정
	var url_2 = "&sm=tab_srt&sort=1&photo=0&field=0&reporter_article=&pd=0&ds=&de=&docid=&nso=so%3Add%2Cp%3Aall%2Ca%3Aall&mynews=0&refresh_start=0&related=0"

	var datas []dataframe
	//var df[]dataframe
	cha := make(chan []dataframe)
	lenge := len(term)

	for i := 0; i < lenge; i++ {
		go naver_news_getPage(term[i], cha, url_1, url_2)
	}
	for i := 0; i < lenge; i++ {
		dataframe := <-cha
		datas = append(datas, dataframe...)
	}

	//writeExcel(datas)
	number = number + len(datas)

	//fmt.Println("done, extracted", len(datas))
}
func naver_news_getPage(terms string, mainC chan<- []dataframe, url_1 string, url_2 string) {

	c := make(chan dataframe)

	var datas []dataframe
	pageURL := url_1 + terms + url_2
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".news_area")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go naver_news_crawling(card, c, terms, pageURL)

	})

	for i := 0; i < searchCards.Length(); i++ {
		dataed := <-c
		datas = append(datas, dataed)
	}

	mainC <- datas

	return
}
func naver_news_crawling(card *goquery.Selection, c chan<- dataframe, terms string, purl string) {
	keyword := terms
	title := CleanString(card.Find(".news_tit").Text())
	time := CleanString(card.Find("span").Text())
	time = getTime(strings.Trim(time, "분전 "))
	text := CleanString(card.Find(".news_dsc").Text())
	url := purl
	c <- dataframe{
		keyword: keyword,
		time:    time,
		title:   title,
		text:    text,
		url:     url}
}

func writeExcel(datas []dataframe) {
	file, err := os.Create("data.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Keyword", "Ttile", "Time  ", "Text", "Url"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, data := range datas {
		dataSlice := []string{data.keyword, data.title, data.time, data.text, data.url}
		jwErr := w.Write(dataSlice)
		checkErr(jwErr)
	}

}

func getTime(diff string) string {
	nowTime := time.Now()
	dif, _ := strconv.Atoi(diff)
	var pastTime = nowTime
	pastTime = nowTime.Add(time.Duration(dif) * time.Minute * -1)
	return pastTime.String()
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func main() {
	
	startTime := time.Now()
	fmt.Print("star time: ")
	fmt.Println(startTime)
	//naver_news_func(kwd_sk...)
	//naver_blog_func(kwd_sk...)
	naver_news_func(kwd_estate...)
	naver_blog_func(kwd_estate...)
	naver_news_func(kwd_stock...)
	naver_blog_func(kwd_stock...)
	naver_news_func(kwd_sport...)
	naver_blog_func(kwd_sport...)
	//naver_news_func(kwd_etc...)
	//naver_blog_func(kwd_etc...)

	endTime := time.Now()
	fmt.Print("end time: ")
	fmt.Println(endTime)
	elapsedTime := time.Since(startTime)

	fmt.Println("수집데이터 수 : ", number)
	fmt.Println("Golang during time : ", elapsedTime)

}

/*



 */
