package main

import (
  "testing"
)

func TestFetchData(t *testing.T) {
  // Call the fetchData function

  // assert that the values are not null

  if videoData.LivestreamStatus == "" {
    t.Errorf("videoData.LivestreamStatus is null")
  }

}
