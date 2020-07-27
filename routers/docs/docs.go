// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/reports/{id}": {
            "get": {
                "tags": [
                    "Report"
                ],
                "summary": "get reports for the report id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "report id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "get latest report in main branch",
                        "name": "latest",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "coverage report",
                        "schema": {
                            "$ref": "#/definitions/core.Report"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "Report"
                ],
                "summary": "Upload coverage report",
                "parameters": [
                    {
                        "type": "string",
                        "description": "report id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "report",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Git commit SHA",
                        "name": "commit",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "report type",
                        "name": "type",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "branch ref",
                        "name": "branch",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "tag ref",
                        "name": "tag",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "error message",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/reports/{id}/comment/{number}": {
            "post": {
                "tags": [
                    "Report"
                ],
                "summary": "Leave a report summary comment on pull request",
                "parameters": [
                    {
                        "type": "string",
                        "description": "report id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "pull request number",
                        "name": "number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/reports/{id}/repo": {
            "get": {
                "tags": [
                    "Report"
                ],
                "summary": "get repository of the report id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "report id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "repository",
                        "schema": {
                            "$ref": "#/definitions/core.Repo"
                        }
                    },
                    "400": {
                        "description": "error message",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/reports/{id}/treemap/{commit}": {
            "get": {
                "produces": [
                    "image/svg+xml"
                ],
                "tags": [
                    "Report"
                ],
                "summary": "Get coverage difference treemap with main branch",
                "parameters": [
                    {
                        "type": "string",
                        "description": "report id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "commit sha",
                        "name": "commit",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "treemap svg",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/repos": {
            "get": {
                "tags": [
                    "Repository"
                ],
                "summary": "List repositories for all available SCM providers",
                "responses": {
                    "200": {
                        "description": "repositories",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.Repo"
                            }
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "Repository"
                ],
                "summary": "Create new repository for code coverage",
                "parameters": [
                    {
                        "description": "repository to create",
                        "name": "repo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.Repo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Created repository",
                        "schema": {
                            "$ref": "#/definitions/core.Repo"
                        }
                    }
                }
            }
        },
        "/repos/{scm}": {
            "get": {
                "tags": [
                    "Repository"
                ],
                "summary": "List repositories",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM source (github, gitea)",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "repositories",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.Repo"
                            }
                        }
                    }
                }
            }
        },
        "/repos/{scm}/{namespace}/{name}": {
            "get": {
                "tags": [
                    "Repository"
                ],
                "summary": "get repository",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.Repo"
                        }
                    }
                }
            },
            "patch": {
                "tags": [
                    "Repository"
                ],
                "summary": "sync repository information with SCM",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.Repo"
                        }
                    }
                }
            }
        },
        "/repos/{scm}/{namespace}/{name}/content/{path}": {
            "get": {
                "tags": [
                    "Repository"
                ],
                "summary": "Get a file content",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "file path",
                        "name": "path",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "specify git ref, default main branch",
                        "name": "ref",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "file content",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/repos/{scm}/{namespace}/{name}/files": {
            "get": {
                "tags": [
                    "Repository"
                ],
                "summary": "List all files in repository",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "specify git ref, default main branch",
                        "name": "ref",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "files",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/repos/{scm}/{namespace}/{name}/report": {
            "patch": {
                "tags": [
                    "Repository"
                ],
                "summary": "renew repository report id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "updated repository",
                        "schema": {
                            "$ref": "#/definitions/core.Repo"
                        }
                    }
                }
            }
        },
        "/repos/{scm}/{namespace}/{name}/setting": {
            "get": {
                "tags": [
                    "Repository"
                ],
                "summary": "get repository setting",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.RepoSetting"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "Repository"
                ],
                "summary": "update repository setting",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Namespace",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "repository setting",
                        "name": "setting",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.RepoSetting"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.RepoSetting"
                        }
                    }
                }
            }
        },
        "/scm/{scm}/repos": {
            "get": {
                "tags": [
                    "SCM"
                ],
                "summary": "Get repositories from SCM",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SCM source (github, gitea)",
                        "name": "scm",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "repositories",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.Repo"
                            }
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "tags": [
                    "User"
                ],
                "summary": "Get login user",
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/scm": {
            "get": {
                "tags": [
                    "User"
                ],
                "summary": "Get user's SCM binding state",
                "responses": {
                    "200": {
                        "description": "providers",
                        "schema": {
                            "$ref": "#/definitions/user.Providers"
                        }
                    },
                    "404": {
                        "description": "providers",
                        "schema": {
                            "$ref": "#/definitions/user.Providers"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "core.Repo": {
            "type": "object",
            "properties": {
                "branch": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "nameSpace": {
                    "type": "string"
                },
                "private": {
                    "type": "boolean"
                },
                "reportID": {
                    "type": "string"
                },
                "scm": {
                    "type": "SCMProvider"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "core.RepoSetting": {
            "type": "object",
            "properties": {
                "filters": {
                    "type": "FileNameFilters"
                }
            }
        },
        "core.Report": {
            "type": "object",
            "properties": {
                "branch": {
                    "type": "string"
                },
                "commit": {
                    "type": "string"
                },
                "coverage": {
                    "type": "CoverageReport"
                },
                "createdAt": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "reportID": {
                    "type": "string"
                },
                "tag": {
                    "type": "string"
                },
                "type": {
                    "type": "ReportType"
                }
            }
        },
        "user.Providers": {
            "type": "object",
            "additionalProperties": {
                "type": "boolean"
            }
        },
        "user.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "CodeCover API",
	Description: "REST API for CodeCover",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
