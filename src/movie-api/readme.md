## Components

- **Movie API**: Mostly serve the request from end users.
- **Movie gRPC**: Mostly serve data to other services.

### How does list top reviewed movies work?

- When a new review get submited, a new event triggered, and movie consumer process the event to update the top reviewed movies.
