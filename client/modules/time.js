function getTimeDifference(currentTime, lastBroadcast) {
  if (lastBroadcast === undefined) {
    return {
      days: ".",
      hours: ".",
      minutes: ".",
      seconds: ".",
    }
  }

  const timeDifference = currentTime - lastBroadcast
  const days = Math.floor(timeDifference / (1000 * 60 * 60 * 24))
  const hours = Math.floor((timeDifference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((timeDifference % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((timeDifference % (1000 * 60)) / 1000)

  return {
    days: days.toString(),
    hours: hours.toString().padStart(2, "0"),
    minutes: minutes.toString().padStart(2, "0"),
    seconds: seconds.toString().padStart(2, "0"),
  }
}

export function updateTimer(currentTime, lastBroadcast) {
  let timeDifference = getTimeDifference(currentTime, lastBroadcast)
  let timeElements = {
    days: document.getElementById("days"),
    hours: document.getElementById("hours"),
    minutes: document.getElementById("minutes"),
    seconds: document.getElementById("seconds"),
  }
  let daysText = document.getElementById("days-text")

  // Handle case where seconds is 60
  if (timeDifference.seconds === "60") {
    timeDifference.seconds = "00"
  }

  // Handle case where time difference is NaN
  Object.keys(timeDifference).forEach((key) => {
    if (isNaN(timeDifference[key])) {
      timeDifference[key] = "."
      timeElements[key].className = "bounce"
    } else {
      timeElements[key].className = ""
    }
  })

  // Update days text
  if (timeDifference.days.endsWith("1") && !timeDifference.days.endsWith("11")) {
    daysText.textContent = "Day"
  } else {
    daysText.textContent = "Days"
  }

  // Update time elements
  Object.keys(timeElements).forEach((key) => {
    timeElements[key].textContent = timeDifference[key]
  })
}
