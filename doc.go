/*
Package workerpool queues work to a limited number of goroutines.

The purpose of the worker pool is to limit the concurrency of the task
performed by the workers.  This is useful when performing a task that requires
sufficient resources (CPU, memory, etc.), that running too many tasks at the
same time would exhaust resources.

Non-blocking task submission

A task is a function submitted to the worker pool for execution.  Submitting
tasks to this worker pool will not block, regardless of the number of tasks.
Tasks read from the input task queue are immediately dispatched to an available
worker.  If no worker is immediately available, then the task is passed to a go
routine which is started to wait for an available worker.  This clears the task
from the task queue immediately, whether or not a worker is currently
available, and will not block the submission of tasks.

The intent of the worker pool is to limit the concurrency of task execution,
not limit the number of tasks queued to be executed. Therefore, this unbounded
input of tasks is acceptable as the tasks cannot be discarded.  If the number
of inbound tasks is too many to even queue for pending processing, then the
solution is outside the scope of workerpool, and should be solved by
distributing load over multiple systems, and/or storing input for pending
processing in intermediate storage such as a database, file system, distributed
message queue, etc.

Dispatcher

This worker pool uses a single dispatcher goroutine to read tasks from the
input task queue and dispatch them to a worker goroutine.  This allows for a
small input channel, and lets the dispatcher queue as many tasks as are
submitted when there are no available workers (using goroutines).
Additionally, the dispatcher can adjust the number of workers as appropriate
for the work load, without having to utilize locked counters and checks
incurred on task submission.

When no tasks have been submitted for a period of time, a worker is removed by
the dispatcher.  This is done until there are no more workers to remove.  The
minimum number of workers is always zero, because the time to start new workers
is insignificant.

Usage note

It is advisable to use different worker pools for tasks that are bound by
different resources, or that have different resource use patterns.  For
example, tasks that use X Mb of memory may need different concurrency limits
than tasks that use Y Mb of memory.

Credits

This implementation builds on ideas from the following:

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

*/
package workerpool
