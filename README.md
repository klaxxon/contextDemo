# contextDemo
## Golang context demonstration
<br/>
For this demo, I have a main context which can show how a cancel can cascade through
to all other derived contexts.  The child contexts also create their own context off of the parent
so they can cancel just their children.
<br/>
<br/>
In the first run, the parent is going to cancel before any of the children.  This will result
is all goroutines being ended throught the one main cancel with each child making sure their own goroutines end properly.
<br/>
<br/>
	PARENTWAIT = 1
	CHILDWAIT  = 10
<br/>
<br/>

```
// All children are created
2020/08/20 13:38:09 Child 2 started
2020/08/20 13:38:09 Grandchild 2 of 2 started
2020/08/20 13:38:09 Child 0 started
2020/08/20 13:38:09 Grandchild 2 of 0 started
2020/08/20 13:38:09 Grandchild 0 of 2 started
2020/08/20 13:38:09 Grandchild 1 of 2 started
2020/08/20 13:38:09 Child 1 started
2020/08/20 13:38:09 Grandchild 2 of 1 started
2020/08/20 13:38:09 Grandchild 0 of 0 started
2020/08/20 13:38:09 Grandchild 1 of 0 started
2020/08/20 13:38:09 Grandchild 0 of 1 started
2020/08/20 13:38:09 Grandchild 1 of 1 started
// As second later, main parent context cancel() is called
2020/08/20 13:38:10 Main parent cancelling!
2020/08/20 13:38:10 Grandchild 2 of 0 done
2020/08/20 13:38:10 Grandchild 2 of 1 done
2020/08/20 13:38:10 Grandchild 1 of 2 done
2020/08/20 13:38:10 Grandchild 0 of 1 done
2020/08/20 13:38:10 Grandchild 2 of 2 done
2020/08/20 13:38:10 Grandchild 0 of 0 done
2020/08/20 13:38:10 Grandchild 0 of 2 done
2020/08/20 13:38:10 Child 2 done
2020/08/20 13:38:10 Grandchild 1 of 1 done
2020/08/20 13:38:10 Child 1 done
2020/08/20 13:38:10 Grandchild 1 of 0 done
2020/08/20 13:38:10 Child 0 done
2020/08/20 13:38:10 Finished!
```

<br/>
As you can see all of the grandchildren ended without a need to call each of the childrens cancels, it was passed through to the childs context given to the grandchild.  However, each child still waited for it's child to exit before notifying the main it was done.
<br/>
<br/>
In this example, the parent waits longer than it will take the children to cancel their own children.  The children wait CHILDWAIT seconds plus their ID so they are staggered.
<br/>
	PARENTWAIT = 5
	CHILDWAIT  = 1
<br/>

```
// All children are created
2020/08/20 13:39:16 Child 2 started
2020/08/20 13:39:16 Grandchild 2 of 2 started
2020/08/20 13:39:16 Child 0 started
2020/08/20 13:39:16 Grandchild 2 of 0 started
2020/08/20 13:39:16 Grandchild 0 of 2 started
2020/08/20 13:39:16 Child 1 started
2020/08/20 13:39:16 Grandchild 1 of 0 started
2020/08/20 13:39:16 Grandchild 1 of 2 started
2020/08/20 13:39:16 Grandchild 1 of 1 started
2020/08/20 13:39:16 Grandchild 2 of 1 started
2020/08/20 13:39:16 Grandchild 0 of 0 started
2020/08/20 13:39:16 Grandchild 0 of 1 started
// First child cancels its children
2020/08/20 13:39:17 Child 0 cancelling kids
2020/08/20 13:39:17 Grandchild 2 of 0 done
2020/08/20 13:39:17 Grandchild 0 of 0 done
2020/08/20 13:39:17 Grandchild 1 of 0 done
2020/08/20 13:39:17 Child 0 done
// Second child cancels
2020/08/20 13:39:18 Child 1 cancelling kids
2020/08/20 13:39:18 Grandchild 0 of 1 done
2020/08/20 13:39:18 Grandchild 2 of 1 done
2020/08/20 13:39:18 Grandchild 1 of 1 done
2020/08/20 13:39:18 Child 1 done
// Third child cancels
2020/08/20 13:39:19 Child 2 cancelling kids
2020/08/20 13:39:19 Grandchild 2 of 2 done
2020/08/20 13:39:19 Grandchild 1 of 2 done
2020/08/20 13:39:19 Grandchild 0 of 2 done
2020/08/20 13:39:19 Child 2 done
// Main still waits 5 seconds, but everyone is done so it just finishes
2020/08/20 13:39:21 Main parent cancelling!
2020/08/20 13:39:21 Finished!
```
