dispatch
========
[![GoDoc](https://godoc.org/github.com/huangml/dispatch?status.svg)](https://godoc.org/github.com/huangml/dispatch)

Package dispatch implements a message processing mechanism.

The basic idea is goroutine-per-request, which means each request runs in a
separate goroutine.

                  +  +  +
                  |  |  |
                  |  |  | goroutine per request
                  |  |  |
                  v  v  v
         +---------------------+
         |      Dispatcher     |
         +----------+----------+
           +--------+--------+
           |        |  +-----|----------------------------------+
           v        v  |     v       In-memory entity           |
         +----+  +----+|+---------+                             |
         |Dest|  |Dest|||  Dest   |                             |
         +----+  +----+|+---------+                             |
                       ||  Mutex  |             +------------+  |
                       ||+-------+|   access    |            |  |
                       |||Handler+------------->|            |  |
                       ||+-------+|sequentially | Associated |  |
                       ||+Handler+------------->|            |  |
                       ||+-------+|             | Resources  |  |
                       ||+Handler+------------->|            |  |
                       ||+-------+|             |            |  |
                       |+---------+             +------------+  |
                       |                                        |
                       +----------------------------------------+
