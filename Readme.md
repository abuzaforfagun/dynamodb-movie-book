# MovieBook

## Getting Started: Running the application

- **Option 1**: `docker-compose up`
- **Option 2**: Follow the `Makefile` and start the services using `make [service-name]`. (Please make sure dynamodb and rabbitmq are running)

## Purpose

The purpose of the project is about to develop highly scalable product using microservice architecture, event driven development and dynamodb.

## Project scope

Create a movie book, where admin can create movies with actors. And user can check out movie details and submit their review. And we need to keep in mind that, this is a read heavy applicaiton.

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

Movie creation, is related with few other features, for an example list movies by a single genre, list actor movies. To make those operation scalable, we are going to populate required data for genre and actor details. After storing the data to table, it trigger an event, actor and movie event listener process the event to generate appropriate data for their services.
![Movie Creation](https://github.com/user-attachments/assets/afe41414-02a6-4161-85cd-3202edbdf2b0)

#### Add review

Reviews are related with movie and user. When a review get added, it trigger an event, that event get processed by movie event listener, and user event listener. Movie event listener change the overall score of movie and re calcualte the top reviewed movies. And user event listener populate the data for user details page where he/she can lookup the reviews he/she made.
![Add Review](https://github.com/user-attachments/assets/3dd1d0cf-934f-4d96-a45f-f83950cbf62a)

#### User name update

When a user submit review, for better scalability we have dupliated the user data to movie table, and as a result when a user updated their name, we also need to update the related data in movie table, and to achieve that we have triggered another event after updating the user name, and review event listener update the user name accordingly.

![User name updated](https://github.com/user-attachments/assets/d65bee2a-a321-4042-9379-d0d5042c0a78)
