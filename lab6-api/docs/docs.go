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
        "/api/matches": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Partidos"
                ],
                "summary": "Listar todos los partidos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Match"
                            }
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Partidos"
                ],
                "summary": "Crear nuevo partido",
                "parameters": [
                    {
                        "description": "Datos del partido",
                        "name": "match",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Partidos"
                ],
                "summary": "Obtener un partido por ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Partidos"
                ],
                "summary": "Actualizar partido",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Datos actualizados",
                        "name": "match",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Partidos"
                ],
                "summary": "Eliminar partido",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/extratime": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Estadísticas"
                ],
                "summary": "Registrar tiempo extra",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 2,
                        "description": "Minutos de tiempo extra",
                        "name": "minutes",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/goals": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Estadísticas"
                ],
                "summary": "Actualizar goles del partido",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Goles del equipo local",
                        "name": "homeGoals",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Goles del equipo visitante",
                        "name": "awayGoals",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/redcards": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Estadísticas"
                ],
                "summary": "Registrar tarjeta roja",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Cantidad de tarjetas a añadir",
                        "name": "count",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/yellowcards": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Estadísticas"
                ],
                "summary": "Registrar tarjeta amarilla",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Cantidad de tarjetas a añadir",
                        "name": "count",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Match": {
            "description": "API para gestión de partidos de fútbol",
            "type": "object",
            "required": [
                "awayTeam",
                "homeTeam",
                "matchDate"
            ],
            "properties": {
                "awayGoals": {
                    "type": "integer"
                },
                "awayTeam": {
                    "type": "string"
                },
                "extraTime": {
                    "type": "integer"
                },
                "homeGoals": {
                    "type": "integer"
                },
                "homeTeam": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "matchDate": {
                    "type": "string"
                },
                "redCards": {
                    "type": "integer"
                },
                "yellowCards": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "lab6-api",
	Description:      "API para gestión de partidos de fútbol",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
