# Task

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

# Design

### Prerequisites

- The window for counting the number of requests is a rolling window. 
    - Request count should reflect the number of request in the last 60 seconds
- There can be a burst of requests all at the same time
- The requests may or may not be chronologically ordered

### Design

- Recording requests count on all endpoints and persisting  the request count:
    - A global "object" for recording the requests count (with its own lock) is created when the server is first started
    - This object is made available to every request to update the count
    - Also, an encoded file is created when the server is first started
    - The requests count is written only on clean shutdown of the service
    - The rationale behind it is to avoid to many file writes (every second)
- Data type decisions
    - Could not use atomic operations for add/delete/read of the count in the global object because 
        - Recording a new count and deleting counts older than 60 seconds need to maintain count integrity
        - This is not possible if recording count (increment) and delete are to 2 atomic operations (a lot can change between 2 atomic operations)
    - Timestamp is Unix timestamp (the number of seconds elapsed since January 1, 1970 UTC)
    - Suits better for time calculations where the quantum of time is in seconds
    - Used map of time stamps because
        - Multiple requests can come at the same time stamp
        - Count of all requests in the same timestamp can be aggregated in one entry 

# Testing

### Commands

Single requests

curl --location --request GET 'http://localhost:3000/requests-count'

curl --location --request GET 'http://localhost:3000/test-request'

Run the command in intervals of 5 seconds or a an interval of you choice. The request can be tallied after 60 seconds. 

Burst requests

Install hey - https://github.com/rakyll/hey#installation

hey -m GET -c 100 -n 500000 http://localhost:3000/requests-count

hey -m GET -c 100 -n 500000 http://localhost:3000/test-request