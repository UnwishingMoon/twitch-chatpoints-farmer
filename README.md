# twitch_chatpoints_farmer
Opens a browser tab on the specified twitch channel and clicks to get the reward points

## WARNING: Read before using
As of 28 October 2020 the use of this program could lead you to punishments on your Twitch Account. Use at your **OWN RISK**.

From Twitch Community Guidelines: https://www.twitch.tv/p/legal/community-guidelines/
> Any content or activity that disrupts, interrupts, harms, or otherwise violates the integrity of Twitch services or another user's experience or devices is prohibited. Such activity includes:
>- ...
>- Cheating a Twitch rewards system (such as the Drops or **channel points systems**)

## How to run
Download the repository, build the program, configure the json and execute it
```
git clone https://github.com/UnwishingMoon/twitch_chatpoints_farmer
cd twitch_chatpoints_farmer/cmd
go build -o farmer -ldflags="-s -w" main.go
./farmer
```
- **username** is the username of the channel you want to farm on
- **clientID** is the Client ID of the App registered on the Twitch Dev site
- **clientSecret** is the Client Secret of the App registered on the Twitch Dev site
- **authCookie** is the authentication cookie (called "twilight-user") from your browser instance (copy the entire value inside)
