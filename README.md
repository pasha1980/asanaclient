# Asana client

Demo client for ASANA. Client fetching the data from ASANA by Access Token and saving as json

## Todo

- Test with retry response from ASANA
- Add error mapping from asana (not just throwing)

## Install

1. Install
```shell
go install github.com/pasha1980/asanaclient@latest
```

2. Export you ASANA access token in env
```shell
export ASANA_ACCESS_TOKEN="<YOUR_ASANA_ACCESS_TOKEN>"
```

3. Run
```shell
asanaclient
```


