// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://example.com/support",
            "email": "support@example.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/auth/sign-in": {
            "post": {
                "description": "Logs the user in using their credentials (username and password)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Auth"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Account is locked",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/v1/notifications": {
            "get": {
                "description": "Retrieve notifications for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Notifications"
                ],
                "summary": "Get user notifications",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003ctoken\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of notifications",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Notification"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/v1/profile": {
            "get": {
                "description": "Retrieve the profile of the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Get user profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization token (Bearer \u003ctoken\u003e)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User profile",
                        "schema": {
                            "$ref": "#/definitions/domain.Profile"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/v1/profile/change-email": {
            "post": {
                "description": "Change the email for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Change user email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization token (Bearer \u003ctoken\u003e)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "New email",
                        "name": "new_email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Email changed successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/v1/profile/change-password": {
            "post": {
                "description": "Change the password for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Change user password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization token (Bearer \u003ctoken\u003e)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "New password",
                        "name": "new_password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password changed successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Auth": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "passwd": {
                    "type": "string"
                }
            }
        },
        "domain.Notification": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "domain.Profile": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "string"
                },
                "balance": {
                    "type": "number"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "internet_status": {
                    "type": "boolean"
                },
                "last_name": {
                    "type": "string"
                },
                "middle_name": {
                    "type": "string"
                },
                "next_pay_date": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "tariff": {
                    "type": "string"
                },
                "to_pay": {
                    "type": "number"
                }
            }
        },
        "rest.ResponseError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Example API",
	Description:      "This is a sample server for managing users.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
