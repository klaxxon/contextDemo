# contextDemo
## Golang context demonstration
<br/>
For this demo, I have a main context which can show how a cancel can cascade through
to all other derived contexts.  The child contexts also create their own context off of the parent
so they can cancel just their children.<br/>
<br/>
Main creates the initial context with a cancel function:

```
mainCtx, cancel := context.WithCancel(context.Background())
```
<br/>
Each child creates their own context off of the parent.  This allows a call to the parent cancel() to propogate through to the children.<br/>

```
mctx, cancel := context.WithCancel(ctx)
	defer cancel()
```
<br/>
The defer is a good practice to get into since you can create contexts with timers or deadlines and if they are not cancelled, the timers will hang around (whatever the timer/deadline is) and keep the context alive in memory.<br/>
<br/>
<br/>
In the first run, the parent is going to cancel before any of the children.  This will result in all goroutines being ended through the one main cancel with each child making sure their own goroutines end properly.
<br/>
<br/>
	PARENTWAIT = 1<br/>
	CHILDWAIT  = 10
<br/>
<br/>

```
// All children are created
2020/08/20 14:02:00 Child 0 started
2020/08/20 14:02:00 Grandchild 2 of 0 started
2020/08/20 14:02:00 Grandchild 0 of 0 started
2020/08/20 14:02:00 Grandchild 1 of 0 started
2020/08/20 14:02:00 Child 1 started
2020/08/20 14:02:00 Grandchild 2 of 1 started
2020/08/20 14:02:00 Grandchild 0 of 1 started
2020/08/20 14:02:00 Grandchild 1 of 1 started
2020/08/20 14:02:00 Child 2 started
2020/08/20 14:02:00 Grandchild 0 of 2 started
2020/08/20 14:02:00 Grandchild 2 of 2 started
2020/08/20 14:02:00 Grandchild 1 of 2 started
// One second later, parent issues cancel and everyone exits
2020/08/20 14:02:01 Main parent cancelling!
2020/08/20 14:02:01 Grandchild 0 of 2 done
2020/08/20 14:02:01 Grandchild 1 of 0 done
2020/08/20 14:02:01 Grandchild 0 of 1 done
2020/08/20 14:02:01 Grandchild 0 of 0 done
2020/08/20 14:02:01 Grandchild 2 of 1 done
2020/08/20 14:02:01 Grandchild 2 of 0 done
2020/08/20 14:02:01 Child 0 done
2020/08/20 14:02:01 Grandchild 1 of 1 done
2020/08/20 14:02:01 Child 1 done
2020/08/20 14:02:01 Grandchild 2 of 2 done
2020/08/20 14:02:01 Grandchild 1 of 2 done
2020/08/20 14:02:01 Child 2 done
2020/08/20 14:02:01 Finished!

```
<br/>
As you can see all of the grandchildren ended without a need to call each of the childrens cancels, it was passed through to the childs context given to the grandchild.  However, each child still waited for it's child to exit before notifying the main it was done.
<br/>
<br/>
In this example, the parent waits longer than it will take the children to cancel their own children.  The children wait CHILDWAIT seconds plus their ID so they are staggered.
<br/>
<br/>
	PARENTWAIT = 5<br/>
	CHILDWAIT  = 1
<br/>

```
// All children are created
2020/08/20 13:56:45 Child 0 started
2020/08/20 13:56:45 Grandchild 2 of 0 started
2020/08/20 13:56:45 Grandchild 0 of 0 started
2020/08/20 13:56:45 Grandchild 1 of 0 started
2020/08/20 13:56:45 Child 1 started
2020/08/20 13:56:45 Grandchild 2 of 1 started
2020/08/20 13:56:45 Grandchild 0 of 1 started
2020/08/20 13:56:45 Grandchild 1 of 1 started
2020/08/20 13:56:45 Child 2 started
2020/08/20 13:56:45 Grandchild 2 of 2 started
2020/08/20 13:56:45 Grandchild 1 of 2 started
2020/08/20 13:56:45 Grandchild 0 of 2 started
// Children start cancelling thier children spaced by one second
2020/08/20 13:56:46 Child 0 cancelling kids
2020/08/20 13:56:46 Grandchild 1 of 0 done
2020/08/20 13:56:46 Grandchild 2 of 0 done
2020/08/20 13:56:46 Grandchild 0 of 0 done
2020/08/20 13:56:46 Child 0 done
2020/08/20 13:56:47 Child 1 cancelling kids
2020/08/20 13:56:47 Grandchild 2 of 1 done
2020/08/20 13:56:47 Grandchild 1 of 1 done
2020/08/20 13:56:47 Grandchild 0 of 1 done
2020/08/20 13:56:47 Child 1 done
2020/08/20 13:56:48 Child 2 cancelling kids
2020/08/20 13:56:48 Grandchild 0 of 2 done
2020/08/20 13:56:48 Grandchild 1 of 2 done
2020/08/20 13:56:48 Grandchild 2 of 2 done
2020/08/20 13:56:48 Child 2 done
// Main still waits 5 seconds, then finishes since all of it's imeediate children are done
2020/08/20 13:56:50 Main parent cancelling!
2020/08/20 13:56:50 Finished!

```
