DefferedQ
===========
DefferedQ is a simple and fast work queue.

Protocol
========
Protocol runs over TCP using UTF-8 encoding.

Interface
---------
General command structure:

```
<ATTR_0> <ATTR_1> <ATTR_2> ... <ATTR_N>\n
```

where ATTR_0 is a command, ATTR_1, ..., ATTR_N are command parameters.

Commands
--------

Commands:

1. Add task to queue

```
ADD <DELAY_MS> <TASK_BODY>
```

2. Get and reserve next task from queue

```
RESERVE
``` 

3. Delete reserved task

```
DELETE <TASK_ID>
``` 

4. Return the reserved task back to the queue

```
RETURN <TASK_ID> <DELAY_MS>
``` 

5. Show statistic of service

```
STATS
``` 

Task lifecycle
==============

Typical success lifecycle
-------------------------

```mermaid
flowchart LR
    ADD-->B[time ticker]-->C[ready]-->RESERVE-->D[reserved]-->E[success execution]-->DELETE
    style ADD fill:#2ecc71,stroke:#2ecc71,color:#fff
    style RESERVE fill:#2ecc71,stroke:#2ecc71,color:#fff
    style DELETE fill:#2ecc71,stroke:#2ecc71,color:#fff
```

Lifecycle with retry
--------------------

```mermaid
flowchart LR
    ADD-->B[time ticker]-->C[ready]-->RESERVE-->D[reserved]-->E[success execution]-->DELETE
    D-.->F[failed execution]
    F-.->RETURN
    RETURN-.->B

    style ADD fill:#2ecc71,stroke:#2ecc71,color:#fff
    style RESERVE fill:#2ecc71,stroke:#2ecc71,color:#fff
    style DELETE fill:#2ecc71,stroke:#2ecc71,color:#fff
    style RETURN fill:#2ecc71,stroke:#2ecc71,color:#fff
```

Lifecycle of "stuck" task
-------------------------
A stuck task is a reserved task that has not been deleted or returned to the queue. A "Watcher" monitors such tasks.

```mermaid
flowchart LR
    ADD-->B[time ticker]-->C[ready]-->RESERVE-->D[reserved]-.->E[no actions]

    E-.->W[Watcher]
    W-.->A[time ticker]-.->SA[attempts counter]-.->AR[auto return]-.->B
    SA-.->AD[auto delete]

    D--->F[success execution]
    F--->DELETE


    style ADD fill:#2ecc71,stroke:#2ecc71,color:#fff
    style RESERVE fill:#2ecc71,stroke:#2ecc71,color:#fff
    style DELETE fill:#2ecc71,stroke:#2ecc71,color:#fff
    style W fill:#f1c40f,stroke:#f1c40f,color:#fff
```

Copyright
=========
Copyright (c) 2023 Vladimir Lila. See LICENSE for details.
