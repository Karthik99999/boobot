# Running the bot

This bot uses GO version 1.19

To run the bot, run the following commands

```bash
git clone https://github.com/Karthik99999/boobot.git
cd boobot
go build
./boobot -token <token>
```

For the google sheet related commands to work, you need to provide a `credentials.json` file in the root of the folder, containing data for a user-managed google service account that has read perms for the google sheets api (go [here](https://cloud.google.com/iam/docs/service-accounts#user-managed) to see how to create the service account and [here](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) to see how to create the json file).
