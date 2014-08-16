## Completed:

1. ZeroMQ load balancing proxy implemented as a proof of concept

1. Literature review (submitted on CATE)

1. Created Docker container running latest version of Go programming
language and added it to the Docker registry

1. Load balancing proxy running in Docker and load balancing multiple
docker containers

1. Performance testing and analyses of proxy demonstrated:
 - messaging systems especially ZeroMQ are not suitable for this project
    - careful use of sockets and threads is critical to get acceptable performance and avoid read errors
    - careful use of pointer and memory allocation (slices) critical for performance and reliability
 - To get an acceptable level of efficiency the proxy has been re-implemented 4 times

1. Researched about fleetctl, but decided to not use it because:
 - fleetctl is not fully released so it has several bugs e.g. the "stop"
command does not work
 - adds unnecessary complexity

1. As a simpler alternative to fleetctl I have implemented a basic
JSON/REST web service for configuring the proxy (on going)

1. To ensure reliable code and to allow for easy refactoring:
 - updated proxy architect to make code testable and more modular
 - created test utilities to enable mocking of net and http go API
components / methods
 - created test utilities to allow simple object comparisons and test
assertion
 - written integration and unit tests for 50% of the code, currently
working of covering most of the code

1. Finish the remaining unit and integration tests

1. Extend the proxy to record metrics for server responses (to enable detection of invalid response)

1. Extend the proxy to cover all scenarios described in the report submitted on CATE:
 - Multi-Version Update
 - Rapid Update - DONE
    - All the requests after the update goes to the new cluster. If the user doesn't define an upgrade mode then the default is instant 
 - New Session Update - DONE
    - Any request with no uuid or with new uuid goes to the new cluster, however, any request with the old uuid goes to the old cluster till it's session is timed out
 - Long Term Update
 
## Remaining: 

1. Taking out modes from route function and making it more encapsulated
 
1. Writing tests for route function

1. Write a makefile

1. Evaluation

1. Extend the metric analysis for multi-version response

1. Dockerisation

1. Write documentation on github and add link to REST service error messages

1. Command line implementation




