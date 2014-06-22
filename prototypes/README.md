performance testing prototypes
==============================

Used following command

```bash
./wrk -t15 -c15 -d120s --latency http://...
```

## wrk testing
 1. simple server with 15 accept loop threads
 1. pipe connecting to two servers (exactly as in step above) also with 15 accept loop threads

