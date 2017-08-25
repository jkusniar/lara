# github.com/jkusniar/lara

LARA REST API.

## Routes

<details>
<summary>`/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/***
	- _GET_
		- [fileServer.func1](/http/server.go#L188)

</details>
<details>
<summary>`/api/v1/*/breed/by-species/{id}`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/breed/by-species/{id}**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllBreedsBySpeciesHandler)-fm](/http/server.go#L131)

</details>
<details>
<summary>`/api/v1/*/city`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/city**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchCityHandler)-fm](/http/server.go#L132)

</details>
<details>
<summary>`/api/v1/*/gender`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/gender**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllGendersHandler)-fm](/http/server.go#L128)

</details>
<details>
<summary>`/api/v1/*/owner/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/owner/***
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createOwnerHandler)-fm](/http/server.go#L91)

</details>
<details>
<summary>`/api/v1/*/owner/*/{id}/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/owner/***
		- **/{id}/***
			- **/**
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getOwnerHandler)-fm](/http/server.go#L93)
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updateOwnerHandler)-fm](/http/server.go#L94)

</details>
<details>
<summary>`/api/v1/*/patient/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/patient/***
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createPatientHandler)-fm](/http/server.go#L100)

</details>
<details>
<summary>`/api/v1/*/patient/*/{id}/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/patient/***
		- **/{id}/***
			- **/**
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getPatientHandler)-fm](/http/server.go#L102)
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updatePatientHandler)-fm](/http/server.go#L103)

</details>
<details>
<summary>`/api/v1/*/productsearch`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/productsearch**
		- _POST_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchProductHandler)-fm](/http/server.go#L146)

</details>
<details>
<summary>`/api/v1/*/record/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/record/***
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createRecordHandler)-fm](/http/server.go#L109)

</details>
<details>
<summary>`/api/v1/*/record/*/{id}/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/record/***
		- **/{id}/***
			- **/**
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updateRecordHandler)-fm](/http/server.go#L112)
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getRecordHandler)-fm](/http/server.go#L111)

</details>
<details>
<summary>`/api/v1/*/report/income`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/report/income**
		- _POST_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getIncomeStatisticsHandler)-fm](/http/server.go#L143)

</details>
<details>
<summary>`/api/v1/*/search`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/search**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchHandler)-fm](/http/server.go#L138)

</details>
<details>
<summary>`/api/v1/*/search/patient-by-tag/{tag}`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/search/patient-by-tag/{tag}**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchPatientByTagHandler)-fm](/http/server.go#L140)

</details>
<details>
<summary>`/api/v1/*/species`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/species**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllSpeciesHandler)-fm](/http/server.go#L129)

</details>
<details>
<summary>`/api/v1/*/street/by-city/{id}`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/street/by-city/{id}**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchStreetByCityHandler)-fm](/http/server.go#L134)

</details>
<details>
<summary>`/api/v1/*/tag/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/tag/***
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createTagHandler)-fm](/http/server.go#L118)

</details>
<details>
<summary>`/api/v1/*/tag/*/{id}/*`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/tag/***
		- **/{id}/***
			- **/**
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getTagHandler)-fm](/http/server.go#L120)
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updateTagHandler)-fm](/http/server.go#L121)

</details>
<details>
<summary>`/api/v1/*/title`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/title**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllTitlesHandler)-fm](/http/server.go#L126)

</details>
<details>
<summary>`/api/v1/*/unit`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/api/v1/***
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L87)
	- **/unit**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllUnitsHandler)-fm](/http/server.go#L127)

</details>
<details>
<summary>`/login`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/login**
	- _POST_
		- [(*Server).(github.com/jkusniar/lara/http.authenticationHandler)-fm](/http/server.go#L85)

</details>
<details>
<summary>`/ping`</summary>

- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L30)
- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- **/ping**
	- _GET_
		- [(*Server).Router.func1](/http/server.go#L80)

</details>

Total # of routes: 22
