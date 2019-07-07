Distributed Hash Craker
=======================

How To Build
------------

```bash
$ sh ./scripts/build.sh
```

Generate MD5
------------

```bash
$ sh ./scripts/md5sum.sh abc 
900150983cd24fb0d6963f7d28e17f72
```

How To Run
----------

Run server:

```console
$ sh ./master/master
```

Run worker:

```console
$ sh ./worker/worker
```

Decrypt MD5:

```console
$ time sh ./client/client --code 900150983cd24fb0d6963f7d28e17f72
abc

real    0m0.077s
user    0m0.000s
sys     0m0.000s
```
