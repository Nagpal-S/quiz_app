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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/upload/file": {
            "post": {
                "description": "upload file to s3",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User AUth"
                ],
                "summary": "upload file to s3",
                "parameters": [
                    {
                        "type": "file",
                        "description": "user file to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.fileUploadResponse"
                        }
                    }
                }
            }
        },
        "/users/check-user/{phone}": {
            "get": {
                "description": "generateing OTP for the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User AUth"
                ],
                "summary": "generateing OTP for the user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user phone number",
                        "name": "phone",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.OTPResponse"
                        }
                    }
                }
            }
        },
        "/users/edit-user-profile": {
            "post": {
                "description": "This API edit user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User AUth"
                ],
                "summary": "This API edit user profile",
                "parameters": [
                    {
                        "description": "this is user info json",
                        "name": "id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.profileInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.editProfileResponse"
                        }
                    }
                }
            }
        },
        "/users/get-user/{id}": {
            "get": {
                "description": "This API will provide user info bu id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User AUth"
                ],
                "summary": "This API will provide user info bu id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "this is user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.userInfo"
                        }
                    }
                }
            }
        },
        "/users/verify-otp/{userId}/{otp}": {
            "get": {
                "description": "This API will verify user OTP with userId",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User AUth"
                ],
                "summary": "This API will verify user OTP with userId",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user app Id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "user otp",
                        "name": "otp",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.verifyOtpResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.OTPResponse": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "object",
                    "properties": {
                        "otp": {
                            "type": "string",
                            "example": "8162"
                        },
                        "userId": {
                            "type": "integer",
                            "example": 3
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "New user created"
                },
                "status": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "controllers.editProfileResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "User info updated successfully."
                },
                "status": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "controllers.fileUploadResponse": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "object",
                    "properties": {
                        "url": {
                            "type": "string",
                            "example": "https://quizbuck.s3.ap-south-1.amazonaws.com/uploads/1734090491_new.jpg"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "File uploaded successfully"
                },
                "status": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "controllers.profileInfo": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "sn@gmail.com"
                },
                "gender": {
                    "type": "string",
                    "example": "Male"
                },
                "image": {
                    "type": "string",
                    "example": "url-of-the-image"
                },
                "name": {
                    "type": "string",
                    "example": "Shivam Nagpal"
                },
                "phone": {
                    "type": "string",
                    "example": "0987656"
                },
                "userId": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "controllers.userInfo": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "object",
                    "properties": {
                        "ID": {
                            "type": "integer",
                            "example": 3
                        },
                        "created": {
                            "type": "string",
                            "example": "2024-12-10T07:04:37Z"
                        },
                        "email": {
                            "type": "string",
                            "example": "shivam@gmail.com"
                        },
                        "name": {
                            "type": "string",
                            "example": "Shivam"
                        },
                        "otp": {
                            "type": "string",
                            "example": "8162"
                        },
                        "phone": {
                            "type": "string",
                            "example": "9144"
                        },
                        "register": {
                            "type": "string",
                            "example": "1"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "User logged in successfully"
                },
                "status": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "controllers.verifyOtpResponse": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "object",
                    "properties": {
                        "ID": {
                            "type": "integer",
                            "example": 3
                        },
                        "created": {
                            "type": "string",
                            "example": "2024-12-10T07:04:37Z"
                        },
                        "email": {
                            "type": "string",
                            "example": "shivam@gmail.com"
                        },
                        "name": {
                            "type": "string",
                            "example": "Shivam"
                        },
                        "otp": {
                            "type": "string",
                            "example": "8162"
                        },
                        "phone": {
                            "type": "string",
                            "example": "9144"
                        },
                        "register": {
                            "type": "string",
                            "example": "1"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "User logged in successfully"
                },
                "status": {
                    "type": "string",
                    "example": "1"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Quiz App API's",
	Description:      "This is list API's to be used in Quiz App.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}