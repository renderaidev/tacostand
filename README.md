# 🌮 tacostand

Simple, lightweight, and self-hostable Slack bot for collecting standup reports.
The bot is primarily intended for people who are familiar with working with
Slack applications.

## Getting Started

### Prerequisites

- Slack application token
- Slack bot token
- Go
- PostgreSQL

### Creating a Slack application

In order to use this software, you will first have to create a new application
on Slack.

The application needs to have the following **features** enabled:

- Slash Commands
- Event Subscriptions
- Bots
- Permissions

Additionally, enabling **Socket Mode** is necessary for this application.

The application needs to have the following **slash commands**:

- `/register`
- `/adduser`
- `/deluser`
- `/unregister`
- `/addquestion`

The following **Bot Token scopes** are mandatory:

- `channels:history` - for finding standup threads in the report channel
- `channels:join` - for viewing the report channel
- `channels:read` - for reading messages from the report channel
- `chat:write` - for writing messages to the report channel
- `chat:write.customize` - for writing messages as the user who submitted the
  report
- `chat:write.public` - for writing in public channels
- `commands` - to execute slash commands
- `im:history` - to view direct message history with users
- `team:read` - to read information about the team
- `users:read` - to list users and get information about them

The following **bot event subscriptions** are mandatory:

- `message.im` - for receiving direct messages from users

Once you have the application configured correctly, you should add it to your
workspace. You will have two tokens that you need to add to your configuration:
the app token (starting with `xapp-`) and the bot token (starting with `xoxb-`).

### Clone the repository

```sh
git clone git@github.com:SomusHQ/tacostand.git
cd tacostand
```

### Configure the app

```sh
cp .env.example .env
nano .env
```

The following configuration options are available:

- `SLACK_APP_TOKEN` - the app token (`xapp-` prefix)
- `SLACK_BOT_TOKEN` - the bot token (`xoxb-` prefix)
- `SLACK_REPORT_CHANNEL` - the name of the channel to which the reports will be
  posted
- `PGHOST` - the hostname of the PostgreSQL server
- `PGPORT` - the port of the PostgreSQL server
- `PGUSER` - the username to use when connecting to the PostgreSQL server
- `PGPASSWORD` - the password to use when connecting to the PostgreSQL server
- `PGDATABASE` - the name of the PostgreSQL database
- `CRON_EXPRESSION` - the cron expression to use for scheduling the bot's
  standup collection. You may use a tool such as [Crontab Guru][crontab-guru] to
  generate a cron expression.
- `WRAP_UP_TIME` - the amount of minutes before the standup reports are wrapped
  up.

You can also set debug options if you want to work on the app:

- `DEBUG_MODE` - will enable additional debug logging when `true`
- `SKIP_MIGRATIONS` - will skip the auto-migration process, resulting in faster
  startups.

Never set any of these to `true` in production.

### Build the app

Use `go` to build the application:

```sh
go build -o tacostand
```

### Run the bot

```sh
./tacostand
```

If you can't execute `./tacostand` you probably have to specify the permissions
for it:

```
chmod +x tacostand
```

## License

MIT. See the LICENSE file for more details.
