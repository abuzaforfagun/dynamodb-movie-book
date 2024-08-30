## Purpose

- Consume the events where movie service need to do some actions.
- To make the application scalable, we are using separate table, to avoid table scan and smaller amount of data read during query. As a result we need to prepare the data for each service and update the data when any data changed occured. This service is going to serve that purpose.
