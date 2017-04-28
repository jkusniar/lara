# github.com/jkusniar/lara

LARA REST API.

## Routes

<details>
<summary>`/*`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/***
	- _GET_
		- [(*Mux).FileServer.func1](/vendor/github.com/pressly/chi/mux.go#L317)

</details>
<details>
<summary>`/api/v1/breed/by-species/:id`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/breed/by-species/:id**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllBreedsBySpeciesHandler)-fm](/http/server.go#L130)

</details>
<details>
<summary>`/api/v1/city`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/city**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchCityHandler)-fm](/http/server.go#L131)

</details>
<details>
<summary>`/api/v1/gender`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/gender**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllGendersHandler)-fm](/http/server.go#L127)

</details>
<details>
<summary>`/api/v1/owner`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/owner**
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createOwnerHandler)-fm](/http/server.go#L90)

</details>
<details>
<summary>`/api/v1/owner/:id`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/owner**
		- **/:id**
			- **/**
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updateOwnerHandler)-fm](/http/server.go#L93)
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getOwnerHandler)-fm](/http/server.go#L92)

</details>
<details>
<summary>`/api/v1/patient`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/patient**
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createPatientHandler)-fm](/http/server.go#L99)

</details>
<details>
<summary>`/api/v1/patient/:id`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/patient**
		- **/:id**
			- **/**
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getPatientHandler)-fm](/http/server.go#L101)
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updatePatientHandler)-fm](/http/server.go#L102)

</details>
<details>
<summary>`/api/v1/productsearch`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/productsearch**
		- _POST_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchProductHandler)-fm](/http/server.go#L145)

</details>
<details>
<summary>`/api/v1/record`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/record**
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createRecordHandler)-fm](/http/server.go#L108)

</details>
<details>
<summary>`/api/v1/record/:id`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/record**
		- **/:id**
			- **/**
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getRecordHandler)-fm](/http/server.go#L110)
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updateRecordHandler)-fm](/http/server.go#L111)

</details>
<details>
<summary>`/api/v1/report/income`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/report/income**
		- _POST_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getIncomeStatisticsHandler)-fm](/http/server.go#L142)

</details>
<details>
<summary>`/api/v1/search`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/search**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchHandler)-fm](/http/server.go#L137)

</details>
<details>
<summary>`/api/v1/search/patient-by-tag/:tag`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/search/patient-by-tag/:tag**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchPatientByTagHandler)-fm](/http/server.go#L139)

</details>
<details>
<summary>`/api/v1/species`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/species**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllSpeciesHandler)-fm](/http/server.go#L128)

</details>
<details>
<summary>`/api/v1/street/by-city/:id`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/street/by-city/:id**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.searchStreetByCityHandler)-fm](/http/server.go#L133)

</details>
<details>
<summary>`/api/v1/tag`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/tag**
		- **/**
			- _POST_
				- [requirePermission.func1](/http/auth.go#L102)
				- [(*Server).(github.com/jkusniar/lara/http.createTagHandler)-fm](/http/server.go#L117)

</details>
<details>
<summary>`/api/v1/tag/:id`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/tag**
		- **/:id**
			- **/**
				- _PUT_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.updateTagHandler)-fm](/http/server.go#L120)
				- _GET_
					- [requirePermission.func1](/http/auth.go#L102)
					- [(*Server).(github.com/jkusniar/lara/http.getTagHandler)-fm](/http/server.go#L119)

</details>
<details>
<summary>`/api/v1/title`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/title**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllTitlesHandler)-fm](/http/server.go#L125)

</details>
<details>
<summary>`/api/v1/unit`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/api/v1**
	- [(*Server).(github.com/jkusniar/lara/http.requireAuthorizedUser)-fm](/http/server.go#L86)
	- **/unit**
		- _GET_
			- [requirePermission.func1](/http/auth.go#L102)
			- [(*Server).(github.com/jkusniar/lara/http.getAllUnitsHandler)-fm](/http/server.go#L126)

</details>
<details>
<summary>`/login`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/login**
	- _POST_
		- [(*Server).(github.com/jkusniar/lara/http.authenticationHandler)-fm](/http/server.go#L84)

</details>
<details>
<summary>`/ping`</summary>

- [RequestID](/vendor/github.com/pressly/chi/middleware/request_id.go#L63)
- [Logger](/vendor/github.com/pressly/chi/middleware/logger.go#L26)
- [Recoverer](/vendor/github.com/pressly/chi/middleware/recoverer.go#L16)
- **/ping**
	- _GET_
		- [(*Server).Router.func1](/http/server.go#L79)

</details>

Total # of routes: 22
