# macaque

a kubernetes [chaos monkey]("https://netflix.github.io/chaosmonkey/") implementation in go.

```
o(..)o   I'm a chaos macaque.
  (-) __/
```

macaque runs in-cluster, uses a crontab for its schedule, can send kill confirmation to slack/webex, and select pods based on labels.

## deploy

The image is available for both `amd64` and `arm64`. To install : 

```
# apply in the proper namespace (comes with role/rolebinding/sa and minimal permissions)
kubectl apply -f https://raw.githubusercontent.com/k0rventen/macaque/main/macaque.yml

# edit the deployment to suit your preferences
kubectl edit deploy macaque
```


## config

configuration is done either through args or env vars:

```
Usage of macaque:
  -crontab string
        env 'MACAQUE_CRONTAB'
        crontab spec for macaque, eg 0 * * * * for every hour.
    
  -namespace string
        env 'MACAQUE_NAMESPACE'
        optionnal namespace in which to look for pods
        (if undefined, uses the ns from the service account).
    
  -selector string
        env 'MACAQUE_SELECTOR'
        optionnal pod selector to use in app=foo format
        (no selector will match any pod in the given namespace).
    
  -slack-channel string
        env 'MACAQUE_SLACK_CHANNEL'
        optionnal slack channel id.
    
  -slack-token string
        env 'MACAQUE_SLACK_TOKEN'
        optionnal slack bot token.
    
  -timezone string
        env 'MACAQUE_TIMEZONE'
        optionnal timezone to use, eg Europe/Paris (defaults to UTC).
    
  -webex-room-id string
        env 'MACAQUE_WEBEX_ROOM_ID'
        optionnal webex room id.
    
  -webex-token string
        env 'MACAQUE_WEBEX_TOKEN'
        optionnal webex bot token.
```
