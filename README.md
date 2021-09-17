# macaque

a kubernetes [chaos monkey]("https://netflix.github.io/chaosmonkey/") implementation in go.

```

o(..)o   I'm a chaos macaque.
  (-) __/

```

## config

configuration is done through env vars:

Mandatory :
- `MACAQUE_CRONTAB`: the crontab spec that macaque will use (eg 0 * * * * for running each hour)

Optionnal : 
- `MACAQUE_SELECTOR`: the kubernetes pod selector to use if you want to select pods in `app=foo` format **(by default no selector will match any pod in the given namespace)**

- `MACAQUE_SLACK_TOKEN`: the slack bot token for sending messages on slack,
- `MACAQUE_SLACK_CHANNEL`: the slack channel ID,
- `MACAQUE_TIMEZONE`: the timezone to use in `Region/City` format. (defaults to UTC


## deploy

The image is available for both `amd64` and `arm64`. To install : 

```
# apply in the proper namespace
kubectl apply -f https://raw.githubusercontent.com/k0rventen/macaque/main/macaque.yml

# edit the deployment env vars if necessary
kubectl edit deploy macaque
```
