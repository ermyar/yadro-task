# Yadro-task

System prototype for biathlon competitions

## Configuration (json)

- **Laps** - Amount of laps for main distance
- **LapLen** - Length of each main lap
- **PenaltyLen** - Length of each penalty lap
- **FiringLines** - Number of firing lines per lap
- **Start** - Planned start time for the first competitor
- **StartDelta** - Planned interval between starts

## Events

All events are characterized by time and event identifier. Outgoing events are events created during program operation. Events related to the "incoming" category cannot be generated and are output in the same form as they were submitted in the input file.

- All events occur sequentially in time. (**_Time of event N+1_**) >= (**_Time of event N_**)
- Time format **_[HH:MM:SS.sss]_**. Trailing zeros are required in input and output

#### Common format for events:

[***time***] **eventID** **competitorID** extraParams

```
Incoming events
EventID | extraParams | Comments
1       |             | The competitor registered
2       | startTime   | The start time was set by a draw
3       |             | The competitor is on the start line
4       |             | The competitor has started
5       | firingRange | The competitor is on the firing range
6       | target      | The target has been hit
7       |             | The competitor left the firing range
8       |             | The competitor entered the penalty laps
9       |             | The competitor left the penalty laps
10      |             | The competitor ended the main lap
11      | comment     | The competitor can`t continue
```

An competitor is disqualified if he/she does not start during his/her start interval. This marked as **NotStarted** in final report.
If the competitor can`t continue it should be marked in final report as **NotFinished**

```
Outgoing events
EventID | extraParams | Comments
32      |             | The competitor is disqualified
33      |             | The competitor has finished
```

## Final report

The final report should contain the list of all registered competitors
sorted by ascending time.

- Total time includes the difference between scheduled and actual start time or **NotStarted**/**NotFinished** marks
- CompetitorID
- Time taken to complete each lap
- Average speed for each lap [m/s]
- Time taken to complete penalty laps
- Average speed over penalty laps [m/s]
- Number of hits/number of shots

## Build/Run/Test

```
make build		   # for build
make test 		   # for testing
go test -v ./... # also testing.
make clean 			 # clean extra files
```
