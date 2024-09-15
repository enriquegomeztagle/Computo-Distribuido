[Source](https://golang.withcodeexample.com/blog/patterns-for-effective-concurrency/)

## Parallelism vs. Concurrency

Before we dive into the intricacies of concurrency patterns, it’s crucial to grasp the fundamental distinction between parallelism and concurrency.

### Parallelism

Parallelism involves the simultaneous execution of multiple tasks, usually with the primary objective of enhancing performance by harnessing the capabilities of multiple processor cores. In a true parallelism scenario, tasks execute concurrently without the need for synchronization or coordination between them. Parallelism is commonly employed for compute-intensive tasks like scientific simulations, rendering, and data processing.

### Concurrency

Concurrency, on the other hand, is a broader concept. It refers to a system’s capability to oversee and execute multiple tasks that overlap in time. These tasks may not necessarily run in parallel but rather in an interleaved manner. Concurrency aims at efficiently utilizing resources, improving responsiveness, and handling tasks concurrently, even in situations where genuine parallelism isn’t achievable.

With this foundational understanding of parallelism and concurrency, let’s delve into practical patterns for achieving effective concurrency in Go.

## Task Decomposition

Task decomposition is a fundamental pattern for designing concurrent systems. This pattern involves breaking down a complex task into smaller, more manageable sub-tasks that can be executed concurrently. This approach not only helps harness the full potential of your hardware but also enhances code modularity and maintainability.

### The Need for Task Decomposition

Imagine a scenario where you need to process a massive dataset. Without task decomposition, you might opt to process each item in a sequential manner. However, this approach can be painfully slow, especially in the context of modern multi-core processors, which remain underutilized.


In this age of advanced software programming, it is important to use parallelism for better results. With the increasing complexity of applications and demands for data processing, there’s a need for coders to write efficient and reliable concurrent programs. To tackle this problem, designers have come up with templates and methods that help in the design and management of concurrent systems. In this article we will take a look at 5 basic concurrency patterns in Go that should be known: distinguishing between parallelism and concurrency, breaking down tasks, using worker pools, cancellation and context, and testing concurrent code.
