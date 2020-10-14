package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	// ハンドラの登録
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/callback", lineHandler)

	fmt.Println("http://localhost:8080 で起動中...")
	// HTTPサーバを起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!!!!"
	fmt.Fprintf(w, msg)
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	// BOTを初期化
	bot, err := linebot.New(
		"(自分のシークレットを入力)",
		"(自分のアクセストークンを入力)",
	)
	if err != nil {
		log.Fatal(err)
	}

	// リクエストからBOTのイベントを取得
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				replyMessage := message.Text
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Print(err)
				}
			// メッセージが位置情報の場合
			case *linebot.LocationMessage:
				sendRestoInfo(bot, event)
			}
			// 他にもスタンプや画像、位置情報など色々受信可能
		}
	}
}

func sendRestoInfo(bot *linebot.Client, e *linebot.Event) {
	msg := e.Message.(*linebot.LocationMessage)

	lat := strconv.FormatFloat(msg.Latitude, 'f', 2, 64)
	lng := strconv.FormatFloat(msg.Longitude, 'f', 2, 64)

	replyMsg := getRestoInfo(lat, lng)

	_, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
	if err != nil {
		log.Print(err)
	}
}

// Response APIレスポンス
type Response struct {
	Results Results `json:"results"`
}

// Results APIレスポンスの内容
type Results struct {
	Shop []Shop `json:"shop"`
}

// Shop レストラン一覧
type Shop struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func getRestoInfo(lat string, lng string) string {
	apikey := "(自分のAPIKEYを入力)"
	url := "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/?format=json" +
		"&key=" + apikey +
		"&lat=" + lat +
		"&lng=" + lng

	// リクエストしてボディを取得
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	info := ""
	for _, shop := range data.Results.Shop {
		info += shop.Name + "\n" + shop.Address + "\n\n"
	}
	return info
}
