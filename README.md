# macaque

a kubernetes chaos monkey implementation in go.

```

o(..)o   I'm a chaos macaque.
  (-) __/

```

## config

configuration is done through env vars:

Mandatory :
- `MACAQUE_CRONTAB`: the crontab spec that macaque will use (eg 0 * * * * for running each hour)
- `MACAQUE_NAMESPACE`: the namespace to use

Optionnal : 
- `MACAQUE_SELECTOR`: optionnal, the kubernetes pod selector to use if you want to select pods in `app=foo` format,
- `MACAQUE_SLACK_TOKEN`: optionnal, the slack bot token for sending messages on slack,
- `MACAQUE_SLACK_CHANNEL`: optionnal, the slack channel ID,
- `MACAQUE_TIMEZONE`: optionnal, the timezone to use if you are not on UTC in `Region/City` format. (defaults to UTC)

**IMPORTANT: by default no selector will match any pod in the given namespace**

## deploy

The image is available for both `amd64` and `arm64`.

make sure your context has the proper namespace, then simply download the deployment spec : 

`curl -LO https://raw.githubusercontent.com/k0rventen/macaque/main/macaque.yml`

edit the env vars, then just `kubectl apply -f macaque.yml`

