// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/actors": {
            "post": {
                "description": "Add acotr with thumbnail image and multiple picture files",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Add new actor",
                "parameters": [
                    {
                        "type": "string",
                        "name": "date_of_birth",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "Upload thumbnail image",
                        "name": "thumbnail",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)",
                        "name": "pictures",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/actors/{id}": {
            "get": {
                "description": "Get actor details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Get actor details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Actor id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response_model.ActorDetails"
                            }
                        }
                    }
                }
            }
        },
        "/actors/{id}/photos": {
            "post": {
                "description": "Add pictures of the actor",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Add picture of actor",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "actor id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)",
                        "name": "pictures",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/movie/{id}/reviews/{review_id}": {
            "delete": {
                "description": "Add review",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Add movie review",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie Id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Review Id",
                        "name": "review_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/movies": {
            "get": {
                "description": "Get all movies",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "movies"
                ],
                "summary": "Get movies",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response_model.Movie"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Add new movie",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "movies"
                ],
                "summary": "Add movie",
                "parameters": [
                    {
                        "description": "movie payload",
                        "name": "AddMovieRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request_model.AddMovie"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/movies/genre/{id}": {
            "get": {
                "description": "Get movies by genre",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "movies"
                ],
                "summary": "Get movies by genre",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Genre Id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response_model.Movie"
                            }
                        }
                    }
                }
            }
        },
        "/movies/{id}": {
            "get": {
                "description": "Get all movies",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "movies"
                ],
                "summary": "Get movie",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response_model.MovieDetails"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete movie by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "movies"
                ],
                "summary": "Delete movie",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Movie Id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/movies/{id}/photos": {
            "post": {
                "description": "Add pictures to the movie",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "movies"
                ],
                "summary": "Add pictures to the movie",
                "parameters": [
                    {
                        "type": "string",
                        "description": "movie id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)",
                        "name": "pictures",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/movies/{id}/reviews": {
            "post": {
                "description": "Add review",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Add movie review",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Movie Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Review payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request_model.AddReview"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Add new user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Add user",
                "parameters": [
                    {
                        "description": "User payload",
                        "name": "AddUserRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request_model.AddUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Get user details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response_model.User"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Update existing user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Update user payload",
                        "name": "UpdateUserRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request_model.UpdateUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    }
                }
            }
        }
    },
    "definitions": {
        "core_models.ActorRole": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                4,
                5,
                6
            ],
            "x-enum-varnames": [
                "LeadHero",
                "LeadHeroin",
                "LeadBillen",
                "Hero",
                "Heroin",
                "Billen",
                "Other"
            ]
        },
        "request_model.ActorRole": {
            "type": "object",
            "properties": {
                "actor_id": {
                    "type": "string"
                },
                "role": {
                    "$ref": "#/definitions/core_models.ActorRole"
                }
            }
        },
        "request_model.AddMovie": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request_model.ActorRole"
                    }
                },
                "release_year": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "request_model.AddReview": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "userId": {
                    "description": "TODO: Need to get from the logged in user",
                    "type": "string"
                }
            }
        },
        "request_model.AddUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "request_model.UpdateUser": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "response_model.ActorDetails": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "movies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response_model.MovieOfActor"
                    }
                },
                "pictures": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "string": {
                    "type": "string"
                },
                "thumbnail_url": {
                    "type": "string"
                }
            }
        },
        "response_model.Creator": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "response_model.Movie": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "release_year": {
                    "type": "integer"
                },
                "score": {
                    "type": "number"
                },
                "thumbnail_url": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "response_model.MovieActor": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "response_model.MovieDetails": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response_model.MovieActor"
                    }
                },
                "genre": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "pictures": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "rating": {
                    "type": "number"
                },
                "release_year": {
                    "type": "string"
                },
                "reviews": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response_model.Review"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "response_model.MovieOfActor": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "score": {
                    "type": "number"
                },
                "thumbnail_url": {
                    "type": "string"
                }
            }
        },
        "response_model.Review": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string"
                },
                "created_by": {
                    "$ref": "#/definitions/response_model.Creator"
                },
                "id": {
                    "type": "integer"
                },
                "score": {
                    "type": "number"
                }
            }
        },
        "response_model.ReviewOfUser": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "movie_id": {
                    "type": "string"
                },
                "movie_thumbnail": {
                    "type": "string"
                },
                "movie_title": {
                    "type": "string"
                },
                "rating": {
                    "type": "integer"
                }
            }
        },
        "response_model.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "reviews": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response_model.ReviewOfUser"
                    }
                },
                "watch_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response_model.WatchListMovie"
                    }
                }
            }
        },
        "response_model.WatchListMovie": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "thumbnail_url": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:5001",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
