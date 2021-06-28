# Running the bot

To run the bot, run the following commands

```bash
git clone https://github.com/Karthik99999/boobot.git
cd boobot
go build
./boobot -token <token>
```

For the google sheet related commands to work, you need to provide a `credentials.json` file in the root of the folder, containing data for a google service account (go [here](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) to see how to create that file).
