# MovieBook

## Development

_How to generate the swagger docs?_

- `swag init -g cmd/api/main.go` from the root directory

<!-- ## Access patterns

MOVIE#1 MOVIE#1 MOVIE MOVIE#1 // MOVIE
USER#1 USER#1 USER USER#1 // USER
MOVIE#1 USER#1 REVIEW MOVIE1#USER1 // REVIEW
ACTOR#1 ACTOR#1 ACTOR ACTOR#1 // ACTOR
MOVIE#1 ACTOR#1 ACTOR-MOVIE MOVIE#1ACTOR#1 // Movie acttors
GENRE#1 MOVIE#1 GENRE GENRE#1MOVIE#1 // GENRE

Get all movies: GSI_PK: MVOIE
Get single movie: PK:MOVIE#1 - movie details - movie reviews - movie actors
Get all users: GSI_PK: USER
Get single user: PK:USER#1 - user details - added review to movies
Get all actors: GSI_PK: ACTOR
Get single actor: ACTOR#1 - actor details - actor movies
Get genre movies: GENRE#1 -->
