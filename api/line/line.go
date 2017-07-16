package line

import (
	"log"
	"net/http"
	"regexp"

	"gopkg.in/mgo.v2/bson"

	"code.olipicus.com/equiz/api/equiz"
	"github.com/line/line-bot-sdk-go/linebot"
)

//LineApp :
type LineApp struct {
	bot          *linebot.Client
	equizService *equiz.EquizService
}

//NewLineApp : New LineApp
func NewLineApp(channelSecret string, channelToken string, service *equiz.EquizService) (*LineApp, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &LineApp{
		bot:          bot,
		equizService: service,
	}, nil
}

//CallbackHandler : handler
func (app *LineApp) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := app.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			log.Println("Invalid Signature")
			log.Println("X-Line-Signature: " + r.Header.Get("X-Line-Signature"))
		} else {
			w.WriteHeader(500)
			log.Println("Unknow error")
		}
		return
	}

	log.Printf("Got events %v", events)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				log.Printf("Got message %s", message.Text)
				re := regexp.MustCompile("^#.*")

				if re.MatchString(message.Text) {
					log.Printf("Message matched event template")
					profile, err := app.bot.GetProfile(event.Source.UserID).Do()

					if err != nil {
						log.Fatalf("Got Error when get Line profile: %v", err)
					}

					u := equiz.User{ID: bson.NewObjectId(), UserName: profile.DisplayName, LineID: profile.UserID, Pic: profile.PictureURL}
					e := equiz.Event{EventTag: message.Text}

					err = app.equizService.RegisterEvent(&u, &e)

					if err != nil {
						switch err {
						case equiz.ErrorNotFoundEvent:
							app.replyText(event.ReplyToken, "ไม่พบ event ที่คุณอยากร่วม")
						case equiz.ErrorUserExist:
							app.replyText(event.ReplyToken, "คุณได้ทำการลงทะเบียนไปแล้ว")
						default:
							log.Printf("Got Error when Register Event: %v", err)
						}
					}

					app.replyText(event.ReplyToken, "ยินดีต้อนรับ เตรียมรอคำถามได้เลย")
				}
			}
		}
	}
}

func (app *LineApp) replyText(replyToken, text string) error {
	if _, err := app.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}
