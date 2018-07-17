Use this issue template to integrate a new language.


- [ ] Announcer
  - [ ] Start() method (should be idempotent)
  - [ ] Stop() method (should be idempotent)
  - [ ] Announcement at 5 second interval
  - [ ] Command line client (args: <service name> <port>)

- [ ] Discoverer
  - [ ] Start() method (should be idempotent)
  - [ ] Get services method (wait parameter, a bool to block for the interval period)
  - [ ] Stop() method (should be idempotent)
  - [ ] Timeout at 10 second interval (service appears offline)
  - [ ] Number of known services GC (limit to 1000)
  - [ ] Command line client (args: <service name>)

- [ ] Validation
  - [ ] Port range
  - [ ] Service name length
  - [ ] Service name character set (ASCII)

- [ ] Tests
  - [ ] Announce/Discover/Timeout
  - [ ] Invalid port
  - [ ] Invalid character set on service name
  - [ ] Service name too long
  - [ ] Foreign protocol rejection (banana)
  - [ ] Foreign service rejection (bar)
  - [ ] Start/stop idempotency test for Announcer
  - [ ] Start/stop idempotency test for Discoverer
