# a simple battery monitor cron tool

I wondered how quickly my laptop battery deteriorates over time. This little
tool creates a sqlite database and records battery values from
/sys/class/power_supply.

For convencience a simple dump method is included.

## install

```
go get github.com/krysopath/watchmon
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


Lastly:

```
watchmon -dbcreate
```

> Creates the sqlite file


## dump

```
watchmon -dump
```
