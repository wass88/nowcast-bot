# NowCast Bot

[![cron](https://github.com/wass88/nowcast-bot/actions/workflows/cron.yml/badge.svg)](https://github.com/wass88/nowcast-bot/actions/workflows/cron.yml)

```
SLACK_WEBHOOK=https://hooks.slack.com/services/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX \
GYAZO_TOKEN=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX \
MAP_ID=10 \
POS_X=490 \
POS_Y=323 \
TRIM=[100,60,600,600] \
./nowcast-bot
```

![image](https://user-images.githubusercontent.com/26019458/175773761-25b80f0b-beba-4617-9177-48b2c5a4f811.png)

* SLACK_WEBHOOK: SlackのWebhook URL
* GYAZO_TOKEN: GyazoのApplication Token
* MAP_ID: map\d
* POS_X, POS_Y: 降水判定箇所
* TRIM: (optional) 画像のトリム `[x, y, w, h]`