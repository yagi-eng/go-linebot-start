package main

import (
	"fmt"
	"log"
	"net/http"

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
			}
			// 他にもスタンプや画像、位置情報など色々受信可能
		}
	}
}
