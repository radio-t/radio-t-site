# add-to-youtube

`add-to-youtube` is command-line tool, that make a video from an audio podcast and upload it to YouTube.

## Before you start

1. You need a [Google Account](https://www.google.com/accounts/NewAccount) to access the Google API Console, request an API key, and register your application. 
1. Create a project in the [Google Developers Console](https://console.developers.google.com/) and [obtain authorization credentials](https://developers.google.com/youtube/registering_an_application) so your application can submit API requests.
1. Put obtained authorization credentials to ENV variable `ADD_RADIOT_TO_YOUTUBE_CLIENT_SECRET_JSON`.
1. After creating your project, make sure the YouTube Data API is one of the services that your application is registered to use:
    1. Go to the API Console and select the project that you just registered.
    1. Visit the [Enabled APIs page](https://console.developers.google.com/apis/enabled). In the list of APIs, make sure the status is ON for the YouTube Data API v3.

## ENV Variables

- `ADD_RADIOT_TO_YOUTUBE_SECRET_TOKEN_PATH` — token fullpath (e.g.: `/secrets/add-to-youtube.json`), may not exists
- `ADD_RADIOT_TO_YOUTUBE_CLIENT_SECRET_JSON` — authorization credentials in json format as string

## Actions

### Authorize an user at Youtube

`add-to-youtube authorize`

### Add a podcast episode to Youtube

`add-to-youtube {episodeID}`