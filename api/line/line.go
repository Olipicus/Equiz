package line

import (
	"log"
	"net/http"
	"regexp"

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
				profile, err := app.bot.GetProfile(event.Source.UserID).Do()

				if err != nil {
					log.Fatalf("Got Error when get Line profile: %v", err)
				}

				re := regexp.MustCompile("^#.*")

				if re.MatchString(message.Text) {
					user := equiz.User{UserName: profile.DisplayName, LineID: profile.UserID, Pic: profile.PictureURL}
					event := equiz.Event{EventTag: message.Text}
					app.equizService.RegisterEvent(&user, &event)
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
