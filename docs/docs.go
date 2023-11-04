// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Sam Komarov",
            "url": "github.com/k0marov",
            "email": "sam@skomarov.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/profiles/me": {
            "get": {
                "description": "Get profile of the caller if it has been created.\nIf profile was not yet created, create it.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profiles"
                ],
                "summary": "Get caller's profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/internal.Profile"
                            }
                        }
                    }
                }
            },
            "patch": {
                "description": "Update profile of the caller, only updating the specified fields.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profiles"
                ],
                "summary": "Update profile of the caller",
                "parameters": [
                    {
                        "description": "fields to update",
                        "name": "account",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal.ProfileUpdatable"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal.Profile"
                        }
                    }
                }
            }
        },
        "/api/v1/profiles/{id}": {
            "get": {
                "description": "Get profile by user id. Returns 404 if profile does not exist.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profiles"
                ],
                "summary": "Get profile by user id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the user for which you want to get its profile.",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal.Profile"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/internal.ClientError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "internal.ClientError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "internal.Profile": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "internal.ProfileUpdatable": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "name": {
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
	BasePath:         "",
	Schemes:          []string{"https", "http"},
	Title:            "profiles-svc",
	Description:      "A microservice for handling user profiles",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
