// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "login user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Logs in an existing user",
                "parameters": [
                    {
                        "description": "Member username and password",
                        "name": "member",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Member"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "description": "logout user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Logs out a currently logged in user",
                "parameters": [
                    {
                        "description": "Member username",
                        "name": "member",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Member"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/member": {
            "get": {
                "description": "get members",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Lists the first 10 members",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Member"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "update user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Updates a members information",
                "parameters": [
                    {
                        "description": "Updated member info",
                        "name": "member",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Member"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Deletes a member from the system",
                "parameters": [
                    {
                        "description": "Member username",
                        "name": "member",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Member"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/member/{username}": {
            "get": {
                "description": "get member by username",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Gets a member's info by their username",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Member"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "resgister user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Registers a new user",
                "parameters": [
                    {
                        "description": "New member",
                        "name": "member",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Member"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Member": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "csrf_token": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "session_token": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "GatorShare API",
	Description:      "Backend API for the GatorShare app",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
