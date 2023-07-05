package main

import (
  "encoding/json"
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
  "time"
)

type VideoData struct {
  LivestreamStatus string    `json:"livestreamStatus"`
  VideoID          string    `json:"videoId"`
  Updated          string    `json:"updated"`
  FetchedAt        time.Time `json:"fetched_at"`
}

var videoData VideoData

func fetchData(ChannelID string, ApiKey string) {
  url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&channelType=any&order=date&type=video&videoCaption=any&videoDefinition=any&videoDimension=any&videoDuration=any&videoEmbeddable=any&videoLicense=any&videoSyndicated=any&videoType=any&key=%s", ChannelID, ApiKey)
  res, err := http.Get(url)
  if err != nil {
    fmt.Println("Failed to fetch the data")
  }
  defer res.Body.Close()

  var result map[string]interface{}
  if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
    fmt.Println("Failed to decode the result")
  }

  items, ok := result["items"].([]interface{})
  if !ok {
    itemsInterf := result["items"]
    fmt.Println("itemsInterf:", itemsInterf)
    if itemsInterf == nil {

      videoData = VideoData{LivestreamStatus: "none", VideoID: "none", Updated: "none", FetchedAt: time.Now()}
      return
    }
    items, ok = itemsInterf.([]interface{})
    if !ok || len(items) == 0 {
      videoData = VideoData{LivestreamStatus: "none", VideoID: "none", Updated: "none", FetchedAt: time.Now()}
      return
    }
  }

  var liveItem map[string]interface{}
  for _, item := range items {
    itemMap := item.(map[string]interface{})
    if itemMap["snippet"].(map[string]interface{})["liveBroadcastContent"].(string) == "live" {
      liveItem = itemMap
      break
    }
  }

  if liveItem != nil {
    livestreamStatus := liveItem["snippet"].(map[string]interface{})["liveBroadcastContent"].(string)
    videoID := liveItem["id"].(map[string]interface{})["videoId"].(string)
    videoData = VideoData{LivestreamStatus: livestreamStatus, VideoID: videoID, Updated: "Stream is Live", FetchedAt: time.Now()}
  } else {
    var mostRecentVideo map[string]interface{}
    for i, item := range items {
      itemMap := item.(map[string]interface{})
      publishedAt := itemMap["snippet"].(map[string]interface{})["publishedAt"]
      if publishedAt == nil || publishedAt == "" {
        continue
      }
      if i == 0 {
        mostRecentVideo = itemMap
        continue
      }
      if mostRecentVideo["snippet"].(map[string]interface{})["publishedAt"].(string) < publishedAt.(string) {
        mostRecentVideo = itemMap
      }
    }
    if mostRecentVideo != nil {
      livestreamStatus := mostRecentVideo["snippet"].(map[string]interface{})["liveBroadcastContent"].(string)
      videoID := mostRecentVideo["id"].(map[string]interface{})["videoId"].(string)
      publishedAt := mostRecentVideo["snippet"].(map[string]interface{})["publishedAt"].(string)

      updated := fetchEndTime(videoID, ApiKey)
      if updated == "" {
        updated = publishedAt // Assign the value of `publishedAt` to `updated`
      }

      videoData = VideoData{LivestreamStatus: livestreamStatus, VideoID: videoID, Updated: updated, FetchedAt: time.Now()}
    } else {
      videoData = VideoData{LivestreamStatus: "none", VideoID: "none", Updated: "none", FetchedAt: time.Now()}
    }

  }

}

func fetchEndTime(videoId, apiKey string) string {
  url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/videos?part=liveStreamingDetails&id=%s&key=%s", videoId, apiKey)
  resp, err := http.Get(url)
  if err != nil {
    return ""
  }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return ""
  }
  var data struct {
    Items []struct {
      LiveStreamingDetails struct {
        ActualEndTime string `json:"actualEndTime"`
      } `json:"liveStreamingDetails"`
    } `json:"items"`
  }
  err = json.Unmarshal(body, &data)
  if err != nil {
    return ""
  }
  if len(data.Items) == 0 {
    return ""
  }
  return data.Items[0].LiveStreamingDetails.ActualEndTime
}

func main() {
  var ChannelID = os.Getenv("CHANNEL_ID")
  var ApiKey = os.Getenv("API_KEY")

  fetchData(ChannelID, ApiKey)

  ticker := time.NewTicker(15 * time.Minute)
  go func() {
    for range ticker.C {
      fetchData(ChannelID, ApiKey)
      fmt.Println("Fetching Data at:", time.Now().Format(time.RFC1123))
    }
  }()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(videoData)
    fmt.Println("New Request at:", time.Now().Format(time.RFC1123)+" "+videoData.Updated)
  })

  fmt.Println("Listening on port", 5000)
  log.Fatal(http.ListenAndServe(":5000", nil))
}
