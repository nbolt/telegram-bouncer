package main

import "os"
import "fmt"
import "net/http"
import "strconv"
import "github.com/go-martini/martini"
import "gopkg.in/redis.v3"
import "gopkg.in/telegram-bot-api.v1"

func main() {
  db, _  := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
  client := redis.NewClient(&redis.Options{
    Addr:     os.Getenv("REDIS_URL"),
    Password: os.Getenv("REDIS_PASS"),
    DB:       db,
  })

  _, err := client.Ping().Result()
  if err != nil {
    fmt.Println("error connecting to redis")
  } else if os.Getenv("TELEGRAM_TOKEN") == "" {
    fmt.Println("no TELEGRAM_TOKEN env var found")
  } else {
    m := martini.Classic()
    
    m.Post("/telegram/:user/**", func(params martini.Params, req *http.Request) string {
      bot, _ := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
      _, err := client.Get(params["user"]).Result()
      if err != nil {
        offset, _ := client.Get("offset").Result()
        num, _  := strconv.Atoi(offset)
        u := tgbotapi.NewUpdate(num+1)
        updates, _ := bot.GetUpdates(u)
        for _, update := range updates {
          id := strconv.Itoa(update.Message.From.ID)
          client.Set("offset", update.UpdateID, 0)
          client.Set(id, update.Message.Chat.ID, 0)
        }
        _, err = client.Get(params["user"]).Result()
      }
      if err == redis.Nil {
        return "user has not messaged bot. please add betcoin_bot to your contact list."
      } else if err != nil {
        return "error with redis connection"
      } else {
        var msg tgbotapi.Chattable
        user, _ := client.Get(params["user"]).Result()
        userid, _ := strconv.Atoi(user)
        switch params["_1"] {
          case "message":
            msg = tgbotapi.NewMessage(userid, params["_2"])
          case "photo":
            file, head, err := req.FormFile("file")
            if err != nil {
              return "bad file"
            }
            fileReader := tgbotapi.FileReader{head.Filename, file, -1}
            msg = tgbotapi.NewPhotoUpload(userid, fileReader)
          case "document":
            file, head, err := req.FormFile("file")
            if err != nil {
              return "bad file"
            }
            fileReader := tgbotapi.FileReader{head.Filename, file, -1}
            msg = tgbotapi.NewDocumentUpload(userid, fileReader)
        }
        _, err = bot.Send(msg)
        if err != nil {
          return "there was an issue communicating with telegram api"
        }

        return "success"
      }
    })

    m.Run()
  }
}
