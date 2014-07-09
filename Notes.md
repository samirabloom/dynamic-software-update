Work done so far:

1. Submitted literature review on CATE

1. Created Docker container running latest version of Go programming language 


1. Load balancing proxy running in Docker and load balancing multiple docker containers
   
1. Performance testing and analysing the developed proxy
        - Analysis demonstrated that messaging systems especially ZeroMQ are not suitable for this project
        - To get the most efficient proxy with lowest over head the proxy was implemented 3 times (currently working with the most efficient proxy)

1. Researched about fleetctl, but decided to not use it because:
        - fleetctl is not fully released so it has several bugs e.g. the "stop" command does not work
        - Add extra complexity  
        
1. As an alternative to fleetctl I have implemented a basic JSON/REST web service for configuring the proxy (This is on going)

1. I have written integration and unit tests for the majority of the code I have written so far


  
