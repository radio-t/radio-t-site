# publisher

Publisher is a docker container designed to execute scipts and programs inside to publish the podcast.

## Before you start

1. You need a [Google Account](https://www.google.com/accounts/NewAccount) to access the Google API Console, request an API key, and register your application. 
1. Create a project in the [Google Developers Console](https://console.developers.google.com/) and [obtain authorization credentials](https://developers.google.com/youtube/registering_an_application) so your application can submit API requests.
1. Put obtained authorization credentials to `/secrets` folder.
1. After creating your project, make sure the YouTube Data API is one of the services that your application is registered to use:
    1. Go to the API Console and select the project that you just registered.
    1. Visit the [Enabled APIs page](https://console.developers.google.com/apis/enabled). In the list of APIs, make sure the status is ON for the YouTube Data API v3.

## Note

`add-to-youtube` apps uses `/secrets` folder to store and keep fresh an access token. 

## Actions

### Authorize an user at Youtube

`docker-compose run --rm publisher add-to-youtube authorize`

### Add a podcast episode to Youtube

`docker-compose run --rm publisher add-to-youtube {episodeID}`