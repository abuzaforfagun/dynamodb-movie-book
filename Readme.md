# MovieBook

## Getting Started: Running the application

- **Option 1**: `docker-compose up`
- **Option 2**: Follow the `Makefile` and start the services using `make [service-name]`. (Please make sure dynamodb and rabbitmq are running)

## Purpose

The purpose of this project is to develop a highly scalable product using microservice architecture, event-driven development, and DynamoDB.

## Project scope

Create a movie book, where an admin can create movies with associated actors, and users can view movie details and submit their reviews. Since this is a read-heavy application, we need to design with scalability in mind.

## How to contribute

- Please create an issue/follow the exiting issues.
- Create a PR containing, change description, unit and integration tests.

## Features

- Add actors
- Add/Update/Delete users
- Add movies
- Add review
- Get all movies
- Search movies by name
- Search by genre
- List top rated movies
- List actor details with movies [Feel free to create a PR]
- List user details with reviews [Feel free to create a PR]
- Submit review
- Delete review [Feel free to create a PR]

## Tech stack

- Golang
- Web API
- gRPC
- RabbitMQ
- Dynamodb

## Testing stretegy

- Unit tests (Feel free to contribute)
- Integration tests (Feel free to contribute)
- E2E tests (Feel free to contribute)
- Manual testing

## Overall architecture design

#### Core components

![Core components](https://github.com/user-attachments/assets/a205cc5d-af14-4a96-96d7-c843b9d7af15)

#### Component Communication

![Service communication](https://github.com/user-attachments/assets/90c4d251-64ee-4374-9a6a-ea727cf2e28d)

Looks scary? Let's break down the user actions:

#### User and actor creation

![User and Actor Creation](https://github.com/user-attachments/assets/cfcfa29b-48f2-40ea-a65e-4ceed83c9d82)

User and actor creation are mostly straight forward, it store the data into their own table.

#### Movie creation

Movie creation is connected to several other features, such as listing movies by genre and listing movies by a particular actor. To make these operations scalable, we will pre-populate the necessary data for genres and actor details. After storing the data in the table, it triggers an event that is processed by actor and movie event listeners to generate the appropriate data for their respective services.

![Movie Creation](https://github.com/user-attachments/assets/afe41414-02a6-4161-85cd-3202edbdf2b0)

#### Add review

Reviews are linked to both movies and users. When a review is added, it triggers an event that is processed by both the movie event listener and the user event listener. The movie event listener updates the overall score of the movie and recalculates the top-reviewed movies. The user event listener updates the user details page, allowing the user to view the reviews they have submitted.

![Add Review](https://github.com/user-attachments/assets/3dd1d0cf-934f-4d96-a45f-f83950cbf62a)

#### User name update

To improve scalability, we have duplicated user data in the movie table when a review is submitted. As a result, if a user updates their name, we also need to update the corresponding data in the movie table. To achieve this, another event is triggered after the user's name is updated, and the review event listener updates the user name accordingly.

![User name updated](https://github.com/user-attachments/assets/d65bee2a-a321-4042-9379-d0d5042c0a78)
