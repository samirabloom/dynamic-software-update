Work done so far:

1. Submitted literature review on CATE

1. Created Docker container running latest version of Go programming language and added it to the Docker registry 

1. have the Load balancing proxy running in Docker and load balancing multiple docker containers
   
1. Done Performance testing and analysing on proxy

        - Analysis demonstrated that messaging systems especially ZeroMQ are not suitable for this project
        - To get the most efficient proxy with lowest overhead the proxy was reimplemented 3 times (currently working with the most efficient proxy)

1. Researched about fleetctl, but decided to not use it because:

        - fleetctl is not fully released so it has several bugs e.g. the "stop" command does not work
        - Add extra complexity  
        
1. As an alternative to fleetctl I have implemented a basic JSON/REST web service for configuring the proxy (on going)

1. I have written integration and unit tests for the majority of the code I have written so far

To do:

1. Finish the remaining tests

1. extend the proxy for the following scenarios described in the report submitted on CATE: 

        - Rapid Update 
        - New Session Update
        - Long Term Update 
        - Multi-Version Update
        


  
