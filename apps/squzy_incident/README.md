# Squzy Incident server

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

Provide possability for handaling users incident from storage

## API

[**GRPC API**](https://github.com/squzy/squzy_proto/blob/master/proto/v1/squzy_incident_server.proto#L11) 

## Environment variables

Bold is required

- PORT(9097) - on with port run squzy_incident
- **MONGO_UR*I* - mongo URI for connect
- MONGO_DB(incident_manager) - mongo DB for connect
- MONGO_COLLECTION(rules) - collection name
- **STORAGE_HOST** - squzy storage host

# Rules definition

For creating the Rule we use [**expr**](https://github.com/antonmedv/expr)  with a number of additional rules.

## Agent rules

First of all, let us provide the agent structure:

```
Agent
    Time    timestamp
    CpuInfo
        Cpus: slice{
            Load  float64
        }
    MemoryInfo
        Mem
            Total       uint64 
            Used        uint64
            Free        uint64
            Shared      uint64
            UsedPercent float64
        Swap
            Total       uint64
            Used        uint64
            Free        uint64
            Shared      uint64
            UsedPercent float64
    DiskInfo
        Disks: map[string]{
            Total       uint64
            Used        uint64
            Free        uint64
            UsedPercent float64
        }
    NetInfo
        Interfaces: map[string]{
            BytesSent     uint64 
            BytesRecv     uint64
            PacketsSent   uint64
            PacketsRecv   uint64
            ErrIn         uint64
            ErrOut        uint64
            DropIn        uint64
            DropOut       uint64
        }
```

You can use requests like `MemoryInfo.Mem.Total` to get the field value.

The new function which can be used with agent:

- `Last(count, filters...)`: receive the last count agents with provided filters. The result of execution is an array of agent struct.

Possible filters:

- `UseTimeFrom("05/02/2020")`: set the time from which entities should be taken

- `UseTimeTo("05/02/2020")`: set the time till which entities should be taken

- `UseType(type)`: set the type of information to recieve about agent. Possible arguemtns:

  - `All`: take all statistics
  - `CPU`: take statistics about CPUs load
  - `Disk`: take statistics about Disk load
  - `Memory`: take statistics about memory load
  - `Net`: take statistics about net load

Example:

```
    any(
        Last(10, UseType(CPU), UseTimeFrom("05/02/2020")),
        {
            all(.CpuInfo.Cpus, {.Load > 80})
        }
    )
```

This rule means next: 
if at least one of the last 10 cpu measurements taken after 05.02.2020 has all the cpus load more then 80 percent,
then there is an incident.

## Application

Application operates with transactions. So, let us provide the transaction structure:

```
TransactionInfo
    Id              string
    ApplicationId   string 
    ParentId        string
    Meta
        Host        string
        Path        string  
        Method      string 
    Name        string
    StartTime   timestamp
    EndTime     timestamp
    Status      TransactionStatus
    Type        TransactionType
    Error       
        Message string
```

The additional to exec functions which can be used with transactions:

- `Last(count, filters...)` - receive the last count transactions info with provided filters.

- `First(count, filters...)` - receive the first count transactions info with provided filters.

- `Index(count, filters...)` - receive the transaction info on given index with provided filters.

- `Duration(transaction)` - calculate the duration of given transaction.

Possible filters:

- `UseTimeFrom("05/02/2020")`: set the time from which entities should be taken

- `UseTimeTo("05/02/2020")`: set the time till which entities should be taken

- `UseType(type)` - set the transaction type. Possible types:

  - `Xhr`: take Xhr transactions
  - `Fetch`: take Fetch transactions
  - `Websocket`: take Websocket transactions
  - `HTTP`: take HTTP transactions
  - `GRPC`: take GRPC transactions
  - `DB`: take DB transactions
  - `Internal`: take Internal transactions
  - `Router`: take Router transactions

- `UseStatus(status)` - set the transaction type. Possible statuses:

  - `Success`: return successful transactions
  - `Failed`: return failed transactions
  
- `UseHost("host")` - set the provided host.

- `UseName("name")` - set the provided transaction name.

- `UsePath("path")` - set the provided path.

- `UseMethod("method")` - set the provided method.


The example of a rule:

```
    len(
        First(10, UseType(HTTP), UseStatus(Success), UsePath("http://localhost"), UseMethod("GET"))
    ) <= 1
```

This rule will create incident, when the first 10 transaction have one ore less successful GET HTTP call to the "http://localhost".

## Scheduler

Application operates with snapshots. So, let us provide the snapshot structure:

```
Snapshot
    Code    SchedulerCode
    Type    SchedulerType  //constant for given scheduler
    Error
        Message string
    Meta
        StartTime   timestamp
        EndTime     timestamp
        Value       string
```

The additional to exec functions which can be used with transactions:

- `Last(count, filters...)` - receive the last count snapshots with provided filters.

- `First(count, filters...)` - receive the first count snapshots with provided filters.

- `Index(count, filters...)` - receive the snapshot info on given index with provided filters.

- `Duration(snapshot)` - calculate the duration of given snapshot.

- `UnixNanoNow()` - return time.Now() in unixNano

- `timeDiff(t1, t2)` - return t1 - t2

- `durationLess(d1, d2)` - check is d1 < d2

- `durationMore(d1, d2)` - check is d1 > d2

- `durationEqual(d1, d2)` - check is d1 == d2

- `durationToSecond(d1)` - convert duration in seconds

- `NowTime()` - return time.Now()

- `float64ToInt64(f1)` - convert float64 to int64

- `getValue(snapshot)` - get value from snapshot

- `unixToTime(u)` - convert unix to time.Time 

- `unixNanoToTime(u)` - convert unixnano to time.Time

- `null` - it is nil

- `mulDuration(int, duration)` - multiply duration on integer

- `Week/Day/Hour/Minute/Second` - constant for duration

Possible filters:

- `UseTimeFrom("05/02/2020")`: set the time from which entities should be taken

- `UseTimeTo("05/02/2020")`: set the time till which entities should be taken

- `UseCode(code)` - set the snapshot code. Possible statuses:

  - `Ok`: return successful snapshots
  - `Error`: return failed snapshots

Example of a rule:

```
    all(
        Last(5), {.Code === Error}
    )
```

This rule will create an incident if last 5 snapshot have an error code.
