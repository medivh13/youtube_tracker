# YouTube Tracker

YouTube Tracker is an application using Golang and YouTube API v3 to collect and analyze video and channel data from YouTube. This application can be used to search for videos, get channel statistics, and calculate popularity scores based on video data.

The purpose of this project is to examine how viral an artist is based on videos using keywords related to that artist. By using these search and analysis features, you can see which artists are currently trending and how their popularity is developing on YouTube.

## Features

1. **Search Videos**: Search for videos based on keywords.
2. **Get Channel Stats**: Retrieve channel statistics such as the number of subscribers, views, and videos.
3. **Calculate Popularity Score**: Calculate the popularity score based on Z-score and median.
4. **Get Trending Videos**: Get the trending videos on YouTube.
5. **Search Channels**: Search for channels based on keywords.

## Prerequisites

- [Golang](https://golang.org/) version 1.16 or newer
- API Key from Google Developer Console for YouTube API v3

## Installation

1. Clone this repository:
    ```sh
    git clone https://github.com/yourusername/youtube-tracker.git
    ```
2. Navigate to the project directory:
    ```sh
    cd youtube-tracker
    ```
3. Create a `.env` file and add your PostgreSQL DSN and YouTube API key:
    ```env
    POSTGRES_DSN="user=username password=password dbname=yourdb sslmode=disable"
    YOUTUBE_API_KEY="your_api_key"
    ```

4. Install the dependencies:
    ```sh
    go mod download
    ```

## Usage

1. Start the application:
    ```sh
    go run main.go
    ```

2. Use the provided endpoints to interact with the application. Here are some examples:

### Search Videos

   - **Endpoint**: `/search?keyword=foo,bar`
   - **Method**: `Get

### analitik with z-score
   - **Endpoint**: `/analitic?keyword=foo,bar`