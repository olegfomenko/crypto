# rho-Pollrad algorithm to find big number factorization via PARCS

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## What is PARCS?

PARCS — *Parallel Asynchronous Recursive Controlled System*

> A task‑tree model that lets you spawn work freely **and** keep tight control over how many tasks run, how long they
> live, and how their errors bubble up.

### Core Idea

| Term             | What it means in PARCS                                                                                |
|------------------|-------------------------------------------------------------------------------------------------------|
| **Parallel**     | Multiple tasks can run at the same time on different cores/nodes.                                     |
| **Asynchronous** | A parent task doesn’t block on its children; it may finish or keep going while they run.              |
| **Recursive**    | Any task can create subtasks, which can in turn spawn their own subtasks, forming a *task tree*.      |
| **Controlled**   | A supervisor restricts concurrency, propagates cancellation/time‑outs, and aggregates results/errors. |


### Shape of a PARCS Tree (Go‑style)

```text
root(ctx, spawn)        // ← your “main” task
 ├─ spawn(taskA)
 │    ├─ spawn(taskA1)
 │    └─ spawn(taskA2)
 └─ spawn(taskB)
      └─ spawn(taskB1)
```

Each task receives:

```go
func(ctx context.Context, spawn func (Task))
```

* **ctx** – carries deadlines, cancellation, per‑request values.
* **spawn** – closure for launching child tasks that inherit the same `ctx`.

## What is a rho-Pollrad algorithm?

## rho-Pollard Algorithm — *Memory‑Light Integer Factorisation*

> Finds a non‑trivial factor of *n* by following a pseudo‑random walk  
> x_{k+1} = f(x_k) mod n, detecting a cycle, and taking a single gcd.


```text
pollard_rho(n):
    if n is even: return 2
    choose random c, x0          # 1 ≤ c, x0 < n
    f(x) = (x*x + c) mod n
    x = y = x0                   # tortoise = hare
    d = 1
    while d == 1:
        x = f(x)                 # tortoise: 1 step
        y = f(f(y))              # hare:     2 steps
        d = gcd(|x - y|, n)
    if d == n:                   # unlucky cycle
        restart with new c and x0
    else:
        return d                 # found non‑trivial factor
```

* **Time (expected):** O(sqrt p), where *p*is the smallest prime factor of *n*
* **Memory:** O(1) — just a few integers

### Practical Tips

* **Randomise** both *c* and *x0*; restart if `gcd == n`.
* Avoid **c=0, +-2** — they fall into tiny degenerate cycles.
* Use random **c**


