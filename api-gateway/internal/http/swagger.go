package http

import (
	"encoding/json"
	"net/http"
)

// SwaggerRoutes serves a minimal Swagger UI pointing to the gateway endpoints.
// This can later be extended to merge specs from additional services.
func SwaggerRoutes() http.Handler {
	mux := http.NewServeMux()

	// Raw spec
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(swaggerSpec)
	})

	// Simple HTML that loads Swagger UI from CDN.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerHTML))
	})

	return mux
}

var swaggerSpec = map[string]interface{}{
	"openapi": "3.0.1",
	"info": map[string]interface{}{
		"title":   "API Gateway",
		"version": "1.0.0",
	},
	"paths": map[string]interface{}{
		"/register": map[string]interface{}{
			"post": map[string]interface{}{
				"summary":     "Register user",
				"description": "Create a new user via auth-service.",
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"email": map[string]interface{}{
										"type":   "string",
										"format": "email",
									},
									"password": map[string]interface{}{
										"type":   "string",
										"format": "password",
									},
								},
								"required": []string{"email", "password"},
							},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "User created",
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"id": map[string]interface{}{
											"type":   "integer",
											"format": "int64",
										},
										"email": map[string]interface{}{
											"type": "string",
										},
										"created_at": map[string]interface{}{
											"type":   "string",
											"format": "date-time",
										},
									},
								},
							},
						},
					},
					"400": map[string]interface{}{
						"description": "Bad request",
					},
				},
			},
		},
		"/login": map[string]interface{}{
			"post": map[string]interface{}{
				"summary":     "Login",
				"description": "Login via auth-service and receive JWT.",
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"email": map[string]interface{}{
										"type":   "string",
										"format": "email",
									},
									"password": map[string]interface{}{
										"type":   "string",
										"format": "password",
									},
								},
								"required": []string{"email", "password"},
							},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "JWT issued",
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"token": map[string]interface{}{
											"type": "string",
										},
									},
								},
							},
						},
					},
					"401": map[string]interface{}{
						"description": "Invalid credentials",
					},
				},
			},
		},
	},
}

const swaggerHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>API Gateway Swagger</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.9.0/swagger-ui.min.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.9.0/swagger-ui-bundle.min.js"></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/swagger/swagger.json',
      dom_id: '#swagger-ui'
    });
  };
</script>
</body>
</html>`
