# a simple battery monitor cron tool

I wondered how quickly my laptop battery deteriorates over time. This little
tool creates a slite database and records battery values from
/sys/class/power_supply.

For convencience a simple dump method id included.

## install

```
go install github.com/krysopath/watchmon

```
> Get the binary and build it.


```
crontab -e

```
> Open the crontab editor


```
 */N * * * * $HOME/go/bin/watchmon | logger
```

> Adapt the `N` to set how many minutes between executions should pass. Also check the PATH.
