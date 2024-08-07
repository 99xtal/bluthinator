package models

type Frame struct {
    ID        int    `json:"id"`
    Episode   string `json:"episode"`
    Timestamp int    `json:"timestamp"`
}

type Subtitle struct {
    ID              int    `json:"id"`
    Episode         string `json:"episode"`
    Text            string `json:"text"`
    StartTimestamp  int    `json:"start_timestamp"`
    EndTimestamp    int    `json:"end_timestamp"`
    FrameTimestamp  int    `json:"frame_timestamp"`
}

type Episode struct {
    EpisodeNumber    int        `json:"episode_number"`
    Season           int        `json:"season"`
    Title            string     `json:"title"`
    Director         string     `json:"director"`
}

type EpisodeResponse struct {
    EpisodeNumber    int        `json:"episode_number"`
    Season           int        `json:"season"`
    Title            string     `json:"title"`
    Director         string     `json:"director"`
    Subtitles        []Subtitle `json:"subtitles"`
}

type FrameResponse struct {
    Frame   Frame     `json:"frame"`
    Episode Episode   `json:"episode"`
    Subtitle Subtitle `json:"subtitle"`
}