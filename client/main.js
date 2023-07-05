import { updateTimer } from "./modules/time.js"
import YouTubePlayer from "youtube-player"

const URL = "/api"
const videoPlayer = document.getElementById("video-player")
const timer = document.getElementById("timer")
videoPlayer.style.display = "none"

let playing = false
let lastBroadcast
let isLive = "none"

// Get the data from the server
async function getData() {
  try {
    const response = await fetch(URL)
    if (!response.ok) {
      throw new Error(response.statusText)
    }
    const data = await response.json()
    return data
  } catch (error) {
    console.error(error)
  }
}

async function getVideoStatus() {
  let player
  const data = await getData()
  let videoID = data.videoId
  isLive = data.livestreamStatus

  if (isLive !== "none") {
    lastBroadcast = "Stream is Live"
    videoPlayer.style.display = "flex"
    if (!playing) {
      timer.style.display = "none"
      player = YouTubePlayer("video-player")
      player.loadVideoById(videoID)
      playing = true
    }
  } else {
    lastBroadcast = new Date(data.updated)
    document.getElementById("video-player").style.display = "none"
    timer.style.display = "flex"
    playing = false
  }
}

document.addEventListener("DOMContentLoaded", function () {
  getVideoStatus()
  setInterval(() => {
    getVideoStatus()
  }, 10 * 1000)

  setInterval(() => {
    if (lastBroadcast === "Stream is Live") {
      return
    }
    const currentTime = Date.now()
    updateTimer(currentTime, lastBroadcast)
  }, 1000)
})
